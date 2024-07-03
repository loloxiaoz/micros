package logger

import (
	"fmt"
	"strings"

	"github.com/spf13/cast"
)

// StackErr 堆栈错误信息
type StackErr struct {
	Filename   string
	Line       int
	Message    string //标准输出报错信息
	StackTrace string
	Code       int    //错误码
	Info       string //错误详情
	Position   string
	Level      int //0最高优先级 1-4 普通优先级 5 可不关注的异常
}

// ErrorInfo 错误信息
func (s *StackErr) ErrorInfo() string {
	return s.Info
}

// Error 错误
func (s *StackErr) Error() string {
	return fmt.Sprintf("%d|%s", s.Code, s.Message)
}

// Stack 堆栈错误信息
func (s *StackErr) Stack() string {
	return fmt.Sprintf("(%s:%d)%s\tStack: %s", s.Filename, s.Line, s.Info, s.StackTrace)
}

// Detail 错误详情
func (s *StackErr) Detail() string {
	return fmt.Sprintf("(%s:%d)%s", s.Filename, s.Line, s.Info)
}

// Format 格式化
func (s *StackErr) Format(tag ...string) (data string) {
	var strs []string
	strs = append(strs, cast.ToString(s.Code))
	strs = append(strs, s.Message)
	strs = append(strs, s.Filename)
	strs = append(strs, cast.ToString(s.Line))
	strs = append(strs, s.Info)
	data = strings.Join(strs, "\t")
	return
}

// SetLevel 设置优先级
func (s *StackErr) SetLevel(lvl int) {
	s.Level = lvl
}

// GetLevel 获取优先级
func (s *StackErr) GetLevel() int {
	return s.Level
}
