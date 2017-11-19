package route

import (
	"github.com/gin-gonic/gin"
	"micros/orm"
)

type Server struct {
	route *gin.Engine
}

func NewServer() *Server {
	server := new(Server)
	server.route = gin.New()
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
