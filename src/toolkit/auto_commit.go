package toolkit

import (
	"micros/orm"
)

var db *orm.DB

func GetXBoxDB() *orm.DB {
	ret, err := X.Get(SQLE)
	if err != nil {
		panic(err)
	}
	db := interface{}(ret).(*orm.DB)
	return db
}

func GetCtxDB() *orm.DB {
	return db
}

func BeforeCommit() {
	db = nil
	ndb := GetXBoxDB().New()
	db = interface{}(ndb).(*orm.DB)
	db.Begin()
}

func AfterCommit() {
	db.Commit()
	db = nil
}

func Rollback() {
	db.Rollback()
	db = nil
}

func CommitAndBegin() {
	db.Commit()
	db.Begin()
}
