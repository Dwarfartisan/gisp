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

	return dot.eval(env, obj, dot.expr)
}

func (dot Dot) eval(env Env, root interface{}, names []Atom) (interface{}, error) {
	if len(names) == 0 {
		return root, nil
	}
	switch obj := root.(type) {
	case Toolkit:
		return dot.evalToolkit(env, obj, names)
	case reflect.Value:
		if obj.Type() == reflect.TypeOf((*Toolkit)(nil)).Elem() && obj.IsValid() {
			return dot.evalToolkit(env, obj.Interface().(Toolkit), names)
		}
		return dot.evalValue(env, obj, names)
	default:
		val := reflect.ValueOf(obj)
		return dot.evalValue(env, val, names)
	}
}

func (dot Dot) evalToolkit(env Env, obj Toolkit, names []Atom) (interface{}, error) {
	name := names[0].Name
	if expr, ok := obj.Content[name]; ok {
		return dot.eval(env, expr, names[1:])
	}
	return nil, fmt.Errorf("Export expr %v from tookit %v but not found in dot %v.%v.",
		name, obj, obj, name)
}

func (dot Dot) evalValue(env Env, val reflect.Value, names []Atom) (interface{}, error) {
	name := names[0].Name
	if val.Kind() == reflect.Struct {
		if field := val.FieldByName(name); field.IsValid() {
			return dot.eval(env, field, names[1:])
		}
	}
	if method := val.MethodByName(name); method.IsValid() {
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
