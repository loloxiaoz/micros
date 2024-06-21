package server

import (
	"net/http"
	"net/http/httputil"
	"time"

	"micros/internal/monitor"
	"micros/pkg/logger"
	"micros/pkg/registry"
	"micros/pkg/toolkit"

	"github.com/gin-gonic/gin"
)

const (
	beginTime = "beginTime"
	endTime   = "endTime"
)

type Server struct {
	Route *gin.Engine
}

func NewServer(name string) *Server {
	server := new(Server)
	//handler
	server.Route = gin.New()
	server.Route.Use(StatBefore())
	server.Route.Use(StatAfter())
	server.Route.Use(Exception())
	server.Route.Use(AutoCommit())
	//prometheus
	server.Route.GET("/metrics", prometheusHandler())
	//consul
	server.Route.GET("/monitor", monitorHandler())
	node := &registry.Node{Id: "1", Address: "127.0.0.1", Port: 8080}
	service := &registry.Service{Name: name, Nodes: []*registry.Node{node}}
	registry.DefaultRegistry.Register(service, registry.RegisterTTL(time.Minute*5))
	//db
	return server
}

func StatBefore() gin.HandlerFunc {
	return func(c *gin.Context) {
		curTime := time.Now().UnixNano() / 1000000
		c.Set(beginTime, curTime)
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
			//prometheus
			monitor.HttpUrlStat.WithLabelValues("200", c.Request.URL.RequestURI()).Add(1)
			monitor.HttpTimeStat.WithLabelValues(c.Request.URL.RequestURI()).Observe(float64(timeCost))
		}()
		c.Next()
	}
}

func AutoCommit() gin.HandlerFunc {
	return func(c *gin.Context) {
//		toolkit.BeforeCommit()
		c.Next()
//		toolkit.AfterCommit()
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
				//db
//				toolkit.Rollback()
				c.AbortWithStatus(http.StatusInternalServerError)
				//sentry
				flags := map[string]string{
					"endpoint": c.Request.URL.RequestURI(),
				}
//				monitor.Report(flags, err, c.Errors)
				monitor.Report(flags, err)
			}
		}()
		c.Next()
	}
}
