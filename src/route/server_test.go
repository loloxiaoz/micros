package route

import (
	"github.com/gin-gonic/gin"
	"micros/orm"
	"net/http"
	"testing"
)

func TestReflect(t *testing.T) {
	r := NewServer()
	r.route.GET("/", func(c *gin.Context) {
		db, _ := c.Get("db")
		tx := interface{}(db).(*orm.DB)
		u2 := orm.User{Name: "luopan"}
		tx.Save(u2)
		c.String(http.StatusOK, "save luopan")
	})

	r.route.Run()
}
