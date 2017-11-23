package server

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"micros/config"
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

func GetCtxDB(ctx *gin.Context) *orm.DB {
	ret, err := ctx.Get(dbExecutor)
	if err != true {
		panic(err)
	}
	db := interface{}(ret).(*orm.DB)
	return db
}

func GetXBoxDB() *orm.DB {
	ret, err := config.X.Get(config.SQLE)
	if err != nil {
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
