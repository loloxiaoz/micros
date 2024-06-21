package server

import (
	"net/http"
	"net/http/httptest"
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
	person.InitTime()
	person.Parent = *person
	person.Create()
	return person
}

func GetPerson(ID int64) *Person {
	person := &Person{}
	GetByID(person.Echo(), ID, person)
	person.Parent = person
	return person
}

func execRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

//func TestORM(t *testing.T) {
//	r := NewServer("micros")
//	db := toolkit.GetXBoxDB()
//	db.DropTable(&Person{})
//	db.CreateTable(&Person{})
//	r.Route.POST("/users", func(c *gin.Context) {
//		person := CreatePerson("micros")
//		toolkit.CommitAndBegin()
//
//		person2 := GetPerson(person.ID)
//		person2.Name = "xxx"
//		person2.Status = true
//		person2.Update()
//
//		person2.Del()
//		c.String(http.StatusOK, "save micros")
//	})
//	w1 := execRequest(r.Route, "POST", "/users")
//	toolkit.AssertEqual(w1.Code, 200, t)
//
//}

//func TestORMException(t *testing.T) {
//	r := NewServer("micros")
//	r.Route.POST("/users2", func(c *gin.Context) {
//		person := Person{Name: "micros exception"}
//		person.Save()
//		panic("exception")
//		c.String(http.StatusOK, "save micros exception")
//	})
//	w := execRequest(r.Route, "POST", "/users2")
//	toolkit.AssertEqual(w.Code, 500, t)
//}
