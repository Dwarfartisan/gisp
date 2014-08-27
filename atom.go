package gisp

import (
	"fmt"
	"reflect"
)

type Atom struct {
	Name string
	Type reflect.Type
}

func (this Atom) String() string {
	return fmt.Sprintf("%v::%v", this.Name, this.Type)
}

func (this Atom) Eval(env Env) (interface{}, error) {
	if value, ok := env.Lookup(this.Name); ok {
		return value, nil
	} else {
		return nil, fmt.Errorf("value of atom %s not found", this.Name)
	}
}

type Variable struct {
	Atom
	value interface{}
}

func DefineVariable(name string, typ reflect.Type, val interface{}) Variable {
	return Variable{Atom{name, typ}, val}
}

func (this *Variable) Set(value interface{}) error {

}
func (this Variable) Get() (interface{}, error) {

}
