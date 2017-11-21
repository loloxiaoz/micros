package server

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"micros/orm"
)

type Context struct {
	gin.Context
}

//func (c *Context) GetDB() *orm.DB {
//	ret, err := c.Get(dbExecutor)
//	if err != true {
//		panic(err)
//	}
//	db := interface{}(ret).(*orm.DB)
//	return db
//}

func GetDB(c *gin.Context) *orm.DB {
	ret, err := c.Get(dbExecutor)
	if err != true {
		panic(err)
	}
	db := interface{}(ret).(*orm.DB)
	return db
}

func PrometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
