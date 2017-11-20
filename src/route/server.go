package route

import (
	"github.com/gin-gonic/gin"
	"micros/logger"
	"micros/orm"
	"micros/toolkit"
	"net/http/httputil"
)

type Server struct {
	route *gin.Engine
}

func NewServer() *Server {
	server := new(Server)
	server.route = gin.New()
	server.route.Use(Exception())
	server.route.Use(AutoCommit())
	return server
}

func AutoCommit() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := orm.OpenConnection()
		tx := db.Begin()
		c.Set("db", tx)
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
				logger.GetIns().Http("[Recovery] panic recovered:", string(httprequest), err, string(stack[:]))
				db, _ := c.Get("db")
				tx := interface{}(db).(*orm.DB)
				tx.Rollback()
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}
