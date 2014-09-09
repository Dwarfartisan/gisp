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

type Dot struct {
	obj  interface{}
	expr []Atom
}

func (dot Dot) Eval(env Env) (interface{}, error) {
	if len(dot.expr) < 1 {
		return nil, fmt.Errorf("The dot %v too short.", dot)
	}
	obj, err := Eval(env, dot.obj)
	if err != nil {
		return nil, err
	}
	val := reflect.ValueOf(obj)

	return dot.eval(env, val, dot.expr)
}

func (dot Dot) eval(env Env, obj reflect.Value, names []Atom) (interface{}, error) {
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
	data, err := p.Many1(p.Bind_(p.Rune('.'), atomNameParser))(st)
	if err != nil {
		return nil, err
	}
	tokens := data.([]interface{})
	expr := make([]Atom, len(tokens))
	for idx, name := range tokens {
		expr[idx] = AA(name.(string))
	}
	return expr, nil
}
