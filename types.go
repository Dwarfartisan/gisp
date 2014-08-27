package gisp

import (
	"reflect"
)

var (
	BOOL   = reflect.TypeOf((*bool)(nil)).Elem()
	STRING = reflect.TypeOf((*string)(nil)).Elem()
	INT    = reflect.TypeOf((*int64)(nil)).Elem()
	FLOAT  = reflect.TypeOf((*float64)(nil)).Elem()
	ANY    = reflect.TypeOf((*interface{})(nil)).Elem()
	ATOM   = reflect.TypeOf((*Atom)(nil)).Elem()
	QUOTE  = reflect.TypeOf((*Quote)(nil)).Elem()
)
