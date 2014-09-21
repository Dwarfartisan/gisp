package gisp

import (
	"fmt"
	p "github.com/Dwarfartisan/goparsec"
	"reflect"
)

type NotStructError struct {
	data interface{}
}

func (err NotStructError) Error() string {
	return fmt.Sprintf("%v not a struct, can't be dot", err.data)
}

type NameInvalid struct {
	Name string
}

func (err NameInvalid) Error() string {
	return fmt.Sprintf("name %s is invalid", err.Name)
}

type DotQ struct {
	obj  interface{}
	expr Quote
}

func (dotq DotQ) Eval(env Env) (interface{}, error) {
	obj, err := Eval(env, dotq.obj)
	if err != nil {
		return nil, err
	}
	expr, _ := dotq.expr.Eval(env)
	fun := L(expr, obj)
	return Eval(env, fun)
}

type Dot struct {
	obj  interface{}
	expr Atom
}

func (dot Dot) Eval(env Env) (interface{}, error) {
	o, err := Eval(env, dot.obj)
	if err != nil {
		return nil, err
	}

	switch obj := o.(type) {
	case Toolkit:
		return dot.evalToolbox(env, obj, dot.expr)
	case reflect.Value:
		if obj.IsValid() {
			inter := obj.Interface()
			switch data := inter.(type) {
			case Toolbox:
				return dot.evalToolbox(env, data, dot.expr)
			}
		}
		return dot.evalValue(env, obj, dot.expr)
	default:
		val := reflect.ValueOf(obj)
		return dot.evalValue(env, val, dot.expr)
	}
}

func (dot Dot) evalToolbox(env Env, obj Toolbox, name Atom) (interface{}, error) {
	if expr, ok := obj.Lookup(name.Name); ok {
		return expr, nil
	}
	return nil, fmt.Errorf("Export expr %v from tookit %v but not found in dot %v.%v.",
		name, obj, obj, name)
}

func (dot Dot) evalValue(env Env, val reflect.Value, name Atom) (interface{}, error) {
	if val.Kind() == reflect.Struct {
		if field := val.FieldByName(name.Name); field.IsValid() {
			return field, nil
		}
	}
	if method := val.MethodByName(name.Name); method.IsValid() {
		return method, nil
	}
	return nil, NameInvalid{name.Name}
}

func DotParser(st p.ParseState) (interface{}, error) {
	name, err := p.Bind_(p.Rune('.'), atomNameParser)(st)
	if err != nil {
		return nil, err
	}
	return AA(name.(string)), nil
}
