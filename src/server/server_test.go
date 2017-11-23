package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"micros/toolkit"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Person struct {
	XEntity `micros:"embedded"`
	Name    string `sql:"size:255"`
	Status  bool
}

func (p Person) Echo() string {
	return "person"
}

func CreatePerson(name string) *Person {
	person := &Person{Name: name}
	person.Create()
	person.parent = *person
	person.Save()
	return person
}

func execRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestORM(t *testing.T) {
	r := NewServer()
	db := toolkit.GetXBoxDB()
	db.DropTable(&Person{})
	db.CreateTable(&Person{})
	r.Route.POST("/users", func(c *gin.Context) {
		person := CreatePerson("micros")
		person.Name = "luopan"
		toolkit.CommitAndBegin()
		person.Update()
		var person2 Person
		GetByID(person2.Echo(), person.ID, person2)
		fmt.Println(person2)
		person.Del()
		c.String(http.StatusOK, "save micros")
	})
	w1 := execRequest(r.Route, "POST", "/users")
	toolkit.AssertEqual(w1.Code, 200, t)

}

//func TestORMException(t *testing.T) {
//	r := NewServer()
//	r.Route.POST("/users2", func(c *gin.Context) {
//		person := Person{Name: "micros exception"}
//		person.Save()
//		panic("exception")
//		c.String(http.StatusOK, "save micros exception")
//	})
//	w := execRequest(r.Route, "POST", "/users2")
//	toolkit.AssertEqual(w.Code, 500, t)
//}
