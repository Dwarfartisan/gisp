package gisp

import (
	"reflect"
)

// Var 定义了变量的结构
type Var interface {
	Get() interface{}
	Set(interface{})
	Type() reflect.Type
}

// OptionVar 定义了可空变量
type OptionVar struct {
	slot reflect.Value
}

// Get 提供变量值或者 nul
func (optVar OptionVar) Get() interface{} {
	if optVar.slot.Elem().IsNil() {
		return nil
	}
	return reflect.Indirect(optVar.slot).Elem().Interface()
}

// Set 实现了赋值逻辑
func (optVar *OptionVar) Set(value interface{}) {
	if value == nil {
		null := reflect.Zero(reflect.PtrTo(optVar.Type()))
		optVar.slot.Elem().Set(null)
		return
	}
	val := reflect.New(optVar.Type())
	val.Elem().Set(reflect.ValueOf(value))
	reflect.Indirect(optVar.slot).Set(val)
}

// Type 实现了类型访问
func (optVar OptionVar) Type() reflect.Type {
	return optVar.slot.Type().Elem().Elem()
}

// DefOption 构造一个可空的变量
func DefOption(typ reflect.Type) OptionVar {
	slot := reflect.New(reflect.PtrTo(typ))
	null := reflect.Zero(reflect.PtrTo(typ))
	slot.Elem().Set(null)
	return OptionVar{slot}
}

// StrictVar 定了非 Option 的变量，它简单的封装了 reflect.Value
type StrictVar struct {
	slot reflect.Value
}

// Get 提供了变量值
func (svar StrictVar) Get() interface{} {
	if svar.slot.IsNil() {
		return nil
	}
	return svar.slot.Elem().Interface()
}

// Set 实现了赋值行为
func (svar *StrictVar) Set(value interface{}) {
	svar.slot.Elem().Set(reflect.ValueOf(value))
}

// Type 给出变量给类型
func (svar StrictVar) Type() reflect.Type {
	return svar.slot.Type().Elem()
}

// DefStrict 构造 Strict 变量
func DefStrict(typ reflect.Type) StrictVar {
	slot := reflect.New(typ)
	return StrictVar{slot}
}

// StrictVarAs 按照值的类型构造对应的 slot
func StrictVarAs(x interface{}) StrictVar {
	slot := DefStrict(reflect.TypeOf(x))
	slot.Set(x)
	return slot
}

// VarSlot 构造一个指定类型的 slot
func VarSlot(typ Type) Var {
	if typ.Option() {
		ret := DefOption(typ.Type)
		return &ret
	}
	ret := DefStrict(typ.Type)
	return &ret
}
