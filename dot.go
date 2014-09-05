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

type Dot []Atom

func (dot Dot) Eval(env Env) (interface{}, error) {
	if len(dot) < 2 {
		return nil, fmt.Errorf("The dot %v too short.", dot)
	}
	first, err := Eval(env, dot[0])
	if err != nil {
		return nil, err
	}
	obj := reflect.ValueOf(first)
	if obj.Kind() != reflect.Struct {
		return nil, NotStructError{first}
	}
	data := make(List, len(dot))
	data[0] = obj
	for idx, item := range dot[1:] {
		data[idx+1] = item
	}
	return dot.eval(env, data)
}

func (dot Dot) eval(env Env, data List) (interface{}, error) {
	if len(data) == 1 {
		return data[0], nil
	}
	name := dot[1].Name
	obj := data[0].(reflect.Value)
	if field := obj.FieldByName(name); field.IsValid() {
		next := field
		d := append(List{next}, data[2:]...)
		return dot.eval(env, d)
	}
	if method := obj.MethodByName(name); method.IsValid() {
		next := method
		d := append(List{next}, data[2:]...)
		return dot.eval(env, d)
	}
	return nil, NameInvalid{name}
}

func DotParser(st p.ParseState) (interface{}, error) {
	data, err := p.SepBy1(atomNameParser, p.Rune('.'))(st)
	if err != nil {
		return nil, err
	}
	tokens := data.([]interface{})
	if len(tokens) == 1 {
		return nil, fmt.Errorf("dot expression except . at last but %v", data)
	}
	dot := make(Dot, len(tokens))
	for idx, name := range tokens {
		dot[idx] = AA(name.(string))
	}
	return dot, nil
}
