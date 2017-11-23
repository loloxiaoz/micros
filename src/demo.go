package main

import (
	"github.com/gin-gonic/gin"
	"micros/orm"
	"micros/server"
	"net/http"
)

func main() {
	r := server.NewServer()
	r.Route.POST("/test1", func(c *gin.Context) {
		db := server.GetDB(c)
		user := orm.User{Name: "test1"}
		db.Save(user)
		c.String(http.StatusOK, "test1 ok")
	})
	r.Route.GET("/test2", func(c *gin.Context) {
		db := server.GetDB(c)
		user := orm.User{Name: "test2"}
		db.Save(user)
		panic("exception")
		c.String(http.StatusOK, "test2 ok")
	})
	r.Route.GET("/test2", func(c *gin.Context) {
		db := server.GetDB(c)
		user := orm.User{Name: "test2"}
		db.Save(user)
		panic("exception")
		c.String(http.StatusOK, "test2 ok")
	})

	r.Route.Run()
}
