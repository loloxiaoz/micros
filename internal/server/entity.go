package server

import (
	"micros/pkg/toolkit"
	"time"
)

type Entity interface {
	Echo() string
}

type XEntity struct {
	ID        int64 `micros:"primary_key"`
	Ver       int   `sql:"version"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Parent Entity `json:"-" sql:"-"`
}

func (xEntity XEntity) Echo() string {
	return "xentity"
}

func (x *XEntity) InitTime() {
	x.ID = toolkit.GenID(1)
	x.Ver = 1
	x.CreatedAt = time.Now()
	x.UpdatedAt = time.Now()
}

func (x *XEntity) Create() {
}

func (x *XEntity) Update() {
}

func (x *XEntity) Del() {
}

func GetByID(table string, ID int64, out interface{}) {
}
