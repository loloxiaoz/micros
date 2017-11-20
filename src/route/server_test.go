package route

import (
	"github.com/gin-gonic/gin"
	"micros/orm"
	"micros/toolkit"
	"net/http"
	"net/http/httptest"
	"testing"
)

func execRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestORM(t *testing.T) {
	r := NewServer()
	r.route.POST("/users", func(c *gin.Context) {
		db, _ := c.Get("db")
		tx := interface{}(db).(*orm.DB)
		user := orm.User{Name: "micros 22"}
		tx.Save(user)
		c.String(http.StatusOK, "save micros")
	})
	w := execRequest(r.route, "POST", "/users")
	toolkit.AssertEqual(w.Code, 200, t)
}

func TestORMException(t *testing.T) {
	r := NewServer()
	r.route.POST("/users2", func(c *gin.Context) {
		db, _ := c.Get("db")
		tx := interface{}(db).(*orm.DB)
		user := orm.User{Name: "micros exception"}
		tx.Save(user)
		panic("exception")
		c.String(http.StatusOK, "save micros exception")
	})
	w := execRequest(r.route, "POST", "/users2")
	toolkit.AssertEqual(w.Code, 500, t)
}
