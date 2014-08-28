package gisp

import (
	"fmt"
	"reflect"
)

// Atom 类型表达基础的 atom 类型
type Atom struct {
	Name string
	Type reflect.Type
}

func (atom Atom) String() string {
	return fmt.Sprintf("%v::%v", atom.Name, atom.Type)
}

// Eval 方法实现 atom 实例的求值行为
func (atom Atom) Eval(env Env) (interface{}, error) {
	if value, ok := env.Lookup(atom.Name); ok {
		return value, nil
	}
	return nil, fmt.Errorf("value of atom %s not found", atom.Name)
}
