package logger

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
)

type DefaultLogger struct {
	*log.Logger
}

func (l *DefaultLogger) Debug(v ...interface{}) {
	l.Output(calldepth, header("DEBUG", fmt.Sprint(v...)))
}

func (l *DefaultLogger) Debugf(format string, v ...interface{}) {
	l.Output(calldepth, header("DEBUG", fmt.Sprintf(format, v...)))
}

func (l *DefaultLogger) Info(v ...interface{}) {
	l.Output(calldepth, header(color.GreenString("INFO "), fmt.Sprint(v...)))
}

func (l *DefaultLogger) Infof(format string, v ...interface{}) {
	l.Output(calldepth, header(color.GreenString("INFO "), fmt.Sprintf(format, v...)))
}

func (l *DefaultLogger) Warn(v ...interface{}) {
	l.Output(calldepth, header(color.YellowString("WARN "), fmt.Sprint(v...)))
}

func (l *DefaultLogger) Warnf(format string, v ...interface{}) {
	l.Output(calldepth, header(color.YellowString("WARN "), fmt.Sprintf(format, v...)))
}

func (l *DefaultLogger) Error(v ...interface{}) {
	l.Output(calldepth, header(color.RedString("ERROR"), fmt.Sprint(v...)))
}

func (l *DefaultLogger) Errorf(format string, v ...interface{}) {
	l.Output(calldepth, header(color.RedString("ERROR"), fmt.Sprintf(format, v...)))
}

func (l *DefaultLogger) Fatal(v ...interface{}) {
	l.Output(calldepth, header(color.MagentaString("FATAL"), fmt.Sprint(v...)))
	os.Exit(1)
}

func (l *DefaultLogger) Fatalf(format string, v ...interface{}) {
	l.Output(calldepth, header(color.MagentaString("FATAL"), fmt.Sprintf(format, v...)))
	os.Exit(1)
}

func (l *DefaultLogger) Panic(v ...interface{}) {
	l.Logger.Panic(v)
}

func (l *DefaultLogger) Panicf(format string, v ...interface{}) {
	l.Logger.Panicf(format, v...)
}

func header(lvl, msg string) string {
	return fmt.Sprintf("%s: %s", lvl, msg)
}
