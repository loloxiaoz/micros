package server

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/micros/config"
	"github.com/micros/logger"
	//	"github.com/micros/monitor"
	"github.com/micros/orm"
	//	"github.com/micros/registry"
	"github.com/micros/toolkit"
)

const (
	REQID     = "reqID"
	beginTime = "beginTime"
	endTime   = "endTime"
)

type Server struct {
	Route *gin.Engine
}

func NewServer(name, configPath string) *Server {
	server := new(Server)
	config.App, _ = config.NewConfig("yaml", configPath)
	//handler
	server.Route = gin.New()
	server.Route.Use(StatBefore())
	server.Route.Use(StatAfter())
	server.Route.Use(Exception())
	server.Route.Use(AutoCommit())
	//prometheus
	//	server.Route.GET("/metrics", prometheusHandler())
	//	//consul
	//	server.Route.GET("/monitor", monitorHandler())
	//	node := &registry.Node{Id: "1", Address: "127.0.0.1", Port: 8080}
	//	service := &registry.Service{Name: name, Nodes: []*registry.Node{node}}
	//	registry.DefaultRegistry.Register(service, registry.RegisterTTL(time.Minute*5))
	//db
	db := OpenConnection()
	toolkit.X.Regist(toolkit.SQLE, db, "init")
	return server
}

func StatBefore() gin.HandlerFunc {
	return func(c *gin.Context) {
		curTime := time.Now().UnixNano() / 1000000
		c.Set(beginTime, curTime)
		reqID := toolkit.GenID(1)
		c.Set(REQID, reqID)
		logger.L.Http("[Start] reqID:", reqID, " "+c.Request.URL.RequestURI())
	}
}

func StatAfter() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			curTime := time.Now().UnixNano() / 1000000
			c.Set(endTime, curTime)
			ret, _ := c.Get(beginTime)
			bTime := ret.(int64)
			timeCost := curTime - bTime
			reqID, _ := c.Get(REQID)
			logger.L.Http("[End] reqID:", reqID, " "+c.Request.URL.RequestURI(), " time cost :", timeCost)
			//prometheus
			//			monitor.HttpUrlStat.WithLabelValues("200", c.Request.URL.RequestURI()).Add(1)
			//			monitor.HttpTimeStat.WithLabelValues(c.Request.URL.RequestURI()).Observe(float64(timeCost))
		}()
		c.Next()
	}
}

func AutoCommit() gin.HandlerFunc {
	return func(c *gin.Context) {
		toolkit.BeforeCommit()
		c.Next()
		toolkit.AfterCommit()
	}
}

func Exception() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				//log
				stack := toolkit.Stack(3)
				httprequest, _ := httputil.DumpRequest(c.Request, false)
				logger.L.Http("[Recovery] panic recovered:", string(httprequest), err, string(stack[:]))
				c.AbortWithStatus(http.StatusInternalServerError)
				//sentry
				//				flags := map[string]string{
				//					"endpoint": c.Request.URL.RequestURI(),
				//				}
				//				monitor.Report(flags, err, c.Errors)
			}
		}()
		c.Next()
	}
}

func OpenConnection() (db *orm.DB) {

	port, _ := config.App.Int("port")
	password, _ := config.App.Int("password")
	dbhost := config.App.String("host") + ":" + strconv.Itoa(port)
	dbhost = fmt.Sprintf("tcp(%v)", dbhost)
	db, err := orm.Open(config.App.String("driver"), fmt.Sprintf("%v:%v@%v/%v?charset=utf8&parseTime=True", config.App.String("username"), strconv.Itoa(password), dbhost, config.App.String("database")))
	if err != nil {
		panic("can't not open connection," + err.Error())
	}

	if os.Getenv("DEBUG") == "true" {
		db.LogMode(true)
	}
	db.DB().SetMaxIdleConns(10)
	return
}
