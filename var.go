package gisp

import (
	"reflect"
)

type Var struct {
	slot interface{}
	typ  reflect.Type
}

func (this Var) Get() interface{} {
	return this.slot
}

func (this *Var) Set(value interface{}) {
	if value == nil {
		this.slot = nil
		return
	}
	pipe := reflect.New(this.Type())
	pipe.Elem().Set(reflect.ValueOf(value))
	this.slot = pipe.Interface()
}

func (this Var) Type() reflect.Type {
	return this.typ.Elem()
}

func DefVar(typ reflect.Type) Var {
	t := reflect.PtrTo(typ)
	return Var{nil, t}
}
