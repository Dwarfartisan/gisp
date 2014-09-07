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

	return dot.eval(env, obj, dot[1:])
}

func (dot Dot) eval(env Env, obj reflect.Value, names Dot) (interface{}, error) {
	if len(names) == 0 {
		return obj, nil
	}
	name := names[0].Name
	if obj.Kind() == reflect.Struct {
		if field := obj.FieldByName(name); field.IsValid() {
			return dot.eval(env, field, names[1:])
		}
	}
	if method := obj.MethodByName(name); method.IsValid() {
		return dot.eval(env, method, names[1:])
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
