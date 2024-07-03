package logger

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"path/filepath"
	"sync"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	formatter "github.com/x-cray/logrus-prefixed-formatter"
)

const endWithZip = ".zip"

type output struct{}

func newOutput() output                            { return output{} }
func (l output) Write(p []byte) (n int, err error) { return }

var rw sync.RWMutex

type Log struct {
	log        *logrus.Logger
	path       string //日志路径
	name       string //日志名称
	watchStop  chan struct{}
	deleteStop chan struct{}
	hook       *lfshook.LfsHook
	index      int32  //日志重命名下标
	formatName string //格式化名字（去掉路径）
	caller     bool
	size       int64 //大小
	maxAge     int64 //有效期
}

var log *Log

func New(name string, opts ...option) *Entry {
	rw.Lock()
	defer rw.Unlock()
	if log == nil {
		log = NewLog("", opts...)
	}
	return log.NewEntry(name)
}

func NewLog(name string, opts ...option) *Log {
	log := &Log{
		log:        logrus.New(),
		watchStop:  make(chan struct{}),
		deleteStop: make(chan struct{}),
	}

	if 0 == len(opts) {
		opts = append(opts, WithLogLevel("debug"))
		opts = append(opts, WithLogName("./log/system.log"))
		opts = append(opts, WithWatchEnable(true))
	}

	log.initLocal(name, opts...)
	return log
}

func (l *Log) NewEntry(name string) *Entry {
	return &Entry{Log: l.log.WithField("model", name), caller: l.caller}
}

func (l *Log) Stop() {
	if nil != l.watchStop {
		l.NewEntry("stop").Warnf("log stop watchStop channel")
		l.watchStop <- struct{}{}
	}
	if nil != l.deleteStop {
		l.NewEntry("stop").Warnf("log stop deleteStop channel")
		l.deleteStop <- struct{}{}
	}
}

func (l *Log) logFileName() string {
	year, month, day := time.Now().Date()
	return fmt.Sprintf("%s_%d-%02d-%02d", l.formatName, year, month, day)
}

func (l *Log) getIndexFileName() {
	fileName := l.logFileName()
	if l.path == "" {
		return
	}
	name := filepath.Base(fileName)
	files, _ := ioutil.ReadDir(l.path)
	tempIndex := 0
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		fileName := f.Name()
		sub := fmt.Sprintf("%s-", name)
		if strings.HasSuffix(fileName, endWithZip) {
			fileName = fileName[:len(fileName)-len(endWithZip)]
		}
		if strings.HasPrefix(fileName, sub) {
			index := fileName[len(sub):]
			if len(index) > 1 && strings.HasPrefix(index, "0") {
				index = index[1:]
			}
			if k, err := strconv.Atoi(index); err == nil {
				if tempIndex < k {
					tempIndex = k
				}
			}

		}
	}
	if tempIndex != 0 && l.name == filepath.Base(l.logFileName()) {
		l.index = int32(tempIndex)
	}
}

/**
初始化设置
*/
func (l *Log) initLocal(name string, opts ...option) {
	level := findLevel(opts...)
	l.caller = findCaller(opts...)
	l.size = findWatchLogsBySize(opts...)
	l.maxAge = findMaxAge(opts...) //日志保存有效期

	l.defPath(name, level, opts...)
	go l.cutLog()
	if findWatcherEnable(opts...) {
		go l.deleteLog(l.path)
	}
}

func (l *Log) defPath(name string, level logrus.Level, opts ...option) {
	l.formatName = findLogName(opts...)
	name = l.logFileName()
	l.name = filepath.Base(name)

	l.path = filepath.Dir(name)

	l.hook = lfshook.NewHook(name, &formatter.TextFormatter{
		TimestampFormat:  "2006-01-02 15:04:05.000000",
		ForceColors:      true,
		QuoteEmptyFields: true,
		FullTimestamp:    true,
	})
	l.log.SetOutput(newOutput())
	l.log.AddHook(l.hook)
	l.log.SetLevel(level)
}

/**
切日志
*/
func (l *Log) cutLog() {
	var name = ""
	tick := time.Tick(time.Second * 1)
	renameLog := func() {
		files, err := ioutil.ReadDir(l.path)
		if err != nil {
			return
		}
		l.getIndexFileName()
		for _, f := range files {
			if f.IsDir() || l.name != f.Name() {
				continue
			}

			if f.Size() >= l.size {
				val := fmt.Sprintf("%s/%s", l.path, l.name)
				rename := fmt.Sprintf("%s-%02d", val, atomic.AddInt32(&l.index, 1))
				if err := os.Rename(val, rename); err != nil {
					l.log.Debugf("log cutLog Rename file:%s error:%s", val, err.Error())
				} else {
					go l.handleFile(rename)
				}
				break
			}
		}

		name = filepath.Base(l.logFileName())
		oldName := l.name
		if oldName != name && name != "" {
			l.name = name

			atomic.StoreInt32(&l.index, 0)
			l.hook.SetDefaultPath(fmt.Sprintf("%s/%s", l.path, l.name))
			delName := fmt.Sprintf("%s/%s", l.path, oldName)
			go l.handleFile(delName)
		}
	}
	for {
		select {
		case <-tick:
			renameLog()
		case <-l.watchStop:
			close(l.watchStop)
			l.watchStop = nil
			return
		}
	}
}

/**
1.压缩file
2.压缩成功后，删除原文件
*/
func (l *Log) handleFile(filename string) {
	//压缩文件，删除源文件
	if err := zipfile(filename); err != nil {
		l.log.Debugf("log cutLog zipfile filename:%s error:%s", filename, err.Error())
		return
	}
	if err := os.Remove(filename); err != nil {
		l.log.Debugf("log cutLog remove filename:%s error:%s", filename, err.Error())
	}
}

func zipfile(filename string) error {
	newFile, err := os.Create(filename + endWithZip)
	if err != nil {
		return err
	}
	defer newFile.Close()

	zipit := zip.NewWriter(newFile)

	defer zipit.Close()

	zipfile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	info, err := zipfile.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Method = zip.Deflate

	writer, err := zipit.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, zipfile)
	return err
}

/**
定时删除日志
*/
func (l *Log) deleteLog(path string) {
	tick := time.Tick(time.Second * 30)
	del := func() {
		var removeArr []string
		var arr []fileInfo
		arr = listFile(path, arr)

		for _, k := range arr {
			modeTime := k.file.ModTime().Unix()
			nowTime := time.Now().Unix()
			result := (modeTime+l.maxAge)-nowTime > 0 //当前文件最后修改时间+1个月(默认)>当前时间
			//fmt.Printf("modeTime:%v maxAge:%v,nowTime:%v,modeTime+l.maxAge:%v\n", modeTime, l.maxAge, nowTime, modeTime+l.maxAge)
			if !result {
				removeArr = append(removeArr, k.path)
			}
		}
		for _, path := range removeArr {
			err := os.Remove(path)
			if err != nil {
				l.log.Debugf("file too long time,remove file err:%s,path:%s", err.Error(), path)
			} else {
				l.log.Debugf("file too long time,remove file ok path:%s", path)
			}
		}
	}

	for {
		select {
		case <-tick:
			del()
		case <-l.deleteStop:
			close(l.deleteStop)
			l.deleteStop = nil
			return
		}
	}

}

type fileInfo struct {
	file os.FileInfo
	path string
}

func listFile(path string, arr []fileInfo) []fileInfo {
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		if f.IsDir() {
			filePath := path + f.Name()
			arr = listFile(filePath, arr)
		} else {
			temp := fileInfo{
				path: path + f.Name(),
				file: f,
			}
			arr = append(arr, temp)
		}
	}
	return arr
}

func (l *Log) SetLevel(level logrus.Level) {
	l.log.Level = level
}

func (l *Log) SetMaxAge(age int64) {
	if age == OneWeek {
		l.maxAge = age
		return
	}
	if age < OneMonth {
		age = OneMonth
		return
	}
	l.maxAge = age
}
