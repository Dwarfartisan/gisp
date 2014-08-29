package gisp

import (
	"reflect"
)

type Var interface {
	Get() interface{}
	Set(interface{})
	Type() reflect.Type
}

type OptionVar struct {
	slot reflect.Value
}

func (this OptionVar) Get() interface{} {
	if this.slot.Elem().IsNil() {
		return nil
	}
	return reflect.Indirect(this.slot).Elem().Interface()
}

func (this *OptionVar) Set(value interface{}) {
	if value == nil {
		null := reflect.Zero(reflect.PtrTo(this.Type()))
		this.slot.Elem().Set(null)
		return
	}
	val := reflect.New(this.Type())
	val.Elem().Set(reflect.ValueOf(value))
	reflect.Indirect(this.slot).Set(val)
}

func (this OptionVar) Type() reflect.Type {
	return this.slot.Type().Elem().Elem()
}

func DefOption(typ reflect.Type) OptionVar {
	slot := reflect.New(reflect.PtrTo(typ))
	null := reflect.Zero(reflect.PtrTo(typ))
	slot.Elem().Set(null)
	return OptionVar{slot}
}

type StrictVar struct {
	slot reflect.Value
}

func (this StrictVar) Get() interface{} {
	if this.slot.IsNil() {
		return nil
	}
	return this.slot.Elem().Interface()
}

func (this *StrictVar) Set(value interface{}) {
	this.slot.Elem().Set(reflect.ValueOf(value))
}

func (this StrictVar) Type() reflect.Type {
	return this.slot.Type().Elem()
}

func DefStrict(typ reflect.Type) StrictVar {
	slot := reflect.New(typ)
	return StrictVar{slot}
}
