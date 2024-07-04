package service

import (
	"micros/pkg/toolkit"
	"time"
)

// Entity 实体接口
type Entity interface {
	ToString() string
}

// XEntity 实体
type XEntity struct {
	ID        int64 `micros:"primary_key"`
	Ver       int   `sql:"version"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// ToString 实体名称
func (x *XEntity) ToString() string {
	return "xentity"
}

// Init 初始化
func (x *XEntity) Init() {
	x.ID = toolkit.GenID(1)
	x.Ver = 1
	x.CreatedAt = time.Now()
	x.UpdatedAt = time.Now()
}

// Student 实体
type Student struct {
	XEntity
	Name string `gorm:"index"`
	Age int `json:"age"`
}