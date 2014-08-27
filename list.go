package gisp

import (
	"fmt"
	"strings"
)

type List []interface{}

func (this List) String() string {
	frags := []string{}
	for _, item := range this {
		frags = append(frags, fmt.Sprintf("%v", item))
	}
	body := strings.Join(frags, " ")
	return fmt.Sprintf("(%s)", body)
}

func (this List) Eval(env Env) (interface{}, error) {
	l := len(this)
	if l == 0 {
		return nil, nil
	} else {
		var lisp interface{}
		switch fun := this[0].(type) {
		case Atom:
			var ok bool
			if lisp, ok = env.Lookup(fun.Name); !ok {
				return nil, fmt.Errorf("any callable named %s not found", fun.Name)
			}
		case List:
			var err error
			lisp, err = fun.Eval(env)
			if err != nil {
				return nil, err
			}
		}
		switch item := lisp.(type) {
		case function:
			return item(env)(this[1:]...)
		case Function:
			return item.Eval(env)
		case Lambda:
			fun := item.Call(this[1:]...)
			return fun.Eval(env)
		case Let:
			return item.Eval(env)
		default:
			return nil, fmt.Errorf("%v:%t is't callable", this[0], this[0])
		}
	}
}
