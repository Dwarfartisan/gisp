package gisp

import (
	"reflect"
)

type Var struct {
	slot reflect.Value
}

func (this Var) Get() interface{} {
	return this.slot
}

func (this *Var) Set(value interface{}) {
	if value == nil {
		typ := this.Type()
		zero := reflect.Zero(reflect.PtrTo(typ))
		this.slot.Set(zero)
		return
	}
	this.slot.Elem().Set(reflect.ValueOf(value))
}

func (this Var) Type() reflect.Type {
	return this.slot.Type().Elem()
}

func DefVar(typ reflect.Type) Var {
	return Var{reflect.New(typ)}
}
