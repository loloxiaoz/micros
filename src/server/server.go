package server

import (
	"github.com/gin-gonic/gin"
	"micros/config"
	"micros/logger"
	"micros/monitor"
	"micros/orm"
	"micros/registry"
	"micros/toolkit"
	"net/http"
	"net/http/httputil"
	"time"
)

const (
	beginTime  = "beginTime"
	endTime    = "endTime"
	dbExecutor = "dbExecutor"
)

type Server struct {
	Route *gin.Engine
}

func NewServer() *Server {
	server := new(Server)
	server.Route = gin.New()
	server.Route.Use(StatBefore())
	server.Route.Use(StatAfter())
	server.Route.Use(Exception())
	server.Route.Use(AutoCommit())
	server.Route.GET("/metrics", PrometheusHandler())
	service := &registry.Service{Name: "micros"}
	err := registry.DefaultRegistry.Register(service)
	if err != nil {
	}
	db := orm.OpenConnection()
	config.X.Regist(config.SQLE, db, "init")
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
			monitor.HttpUrlStat.WithLabelValues("200", c.Request.URL.RequestURI()).Add(1)
			monitor.HttpTimeStat.WithLabelValues(c.Request.URL.RequestURI()).Observe(float64(timeCost))
		}()
		c.Next()
	}
}

func AutoCommit() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := GetXBoxDB().New()
		tx := db.Begin()
		c.Set(dbExecutor, tx)
		c.Next()
		tx.Commit()
	}
}

func Exception() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := toolkit.Stack(3)
				httprequest, _ := httputil.DumpRequest(c.Request, false)
				logger.L.Http("[Recovery] panic recovered:", string(httprequest), err, string(stack[:]))
				tx := GetCtxDB(c)
				tx.Rollback()
				c.AbortWithStatus(http.StatusInternalServerError)
				flags := map[string]string{
					"endpoint": c.Request.URL.RequestURI(),
				}
				monitor.Report(flags, err, c.Errors)
			}
		}()
		c.Next()
	}
}
