package toolkit

import (
	"errors"
)

const (
	ROUTER = "router"
	SQLE   = "sqlExecuter"
)

var (
	X           = newXbox()
	ErrNotFound = errors.New("not found")
	ErrRegisted = errors.New("already regist")
)

type XBox struct {
	objs   map[string]interface{}
	wheres map[string]string
}

func newXbox() *XBox {
	ins := new(XBox)
	ins.objs = make(map[string]interface{})
	ins.wheres = make(map[string]string)
	return ins
}

func (x *XBox) Regist(name string, value interface{}, where string) error {
	if _, ok := x.objs[name]; ok {
		return ErrRegisted
	}
	x.objs[name] = value
	x.wheres[name] = where
	return nil
}

func (x *XBox) Get(name string) (obj interface{}, err error) {
	if _, ok := x.objs[name]; !ok {
		return nil, ErrNotFound
	}
	obj, _ = x.objs[name]
	return obj, nil
}
