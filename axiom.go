package gisp

import (
	"fmt"
	px "github.com/Dwarfartisan/goparsec/parsex"
	"reflect"
)

// Axiom 是基本的 LISP 公理实现，尽可能贴近原始的 LISP 公理描述，但是部分实现对实际的 golang
// 环境做了妥协
var Axiom = Toolkit{
	Meta: map[string]interface{}{
		"name":     "axiom",
		"category": "package",
	},
	Content: map[string]Expr{
		"quote": func(env Env) Element {
			return func(args ...interface{}) (interface{}, error) {
				return Quote{args[0]}, nil
			}
		},
		"var": func(env Env) Element {
			return func(args ...interface{}) (interface{}, error) {
				first := args[0].(Atom)
				slot := VarSlot(first.Type)
				if len(args) == 1 {
					err := env.Defvar(first.Name, slot)
					return nil, err
				}
				value, err := Eval(env, args[1])
				if err != nil {
					return nil, err
				}
				slot.Set(value)
				err = env.Defvar(first.Name, slot)
				return nil, err
			}
		},
		"set": func(env Env) Element {
			return func(args ...interface{}) (interface{}, error) {
				value, err := Eval(env, args[1])
				if err != nil {
					return nil, err
				}
				return set(env, args[0], value)
			}
		},
		"equal": func(env Env) Element {
			return func(args ...interface{}) (interface{}, error) {
				if len(args) != 2 {
					return nil, fmt.Errorf("args error: equal need two args but only",
						args)
				}
				x, err := Eval(env, args[0])
				if err != nil {
					return nil, err
				}
				y, err := Eval(env, args[1])
				if err != nil {
					return nil, err
				}
				return reflect.DeepEqual(x, y), nil
			}
		},
		"cond": func(env Env) Element {
			return func(args ...interface{}) (interface{}, error) {
				cases := args[0].([]interface{})
				l := len(args)
				var els interface{}
				if l > 1 {
					els = args[1]
				} else {
					els = nil
				}

				for _, b := range cases { // FIXME: need a else
					branch := b.([]interface{})
					cond := branch[0].(List)
					result, err := Eval(env, cond)
					if err != nil {
						return nil, err
					}
					if ok := result.(bool); ok {
						return Eval(env, branch[1])
					}
				}
				// else branch
				if els != nil {
					return Eval(env, els)
				}
				return nil, nil
			}
		},
		"car": func(env Env) Element {
			return func(args ...interface{}) (interface{}, error) {
				params, err := GetArgs(
					env,
					px.Binds_(TypeAs(LIST), px.Eof),
					args,
				)
				if err != nil {
					return nil, err
				}
				lisp := params[0]
				return (lisp.(List))[0], nil
			}
		},
		"cdr": func(env Env) Element {
			return func(args ...interface{}) (interface{}, error) {
				params, err := GetArgs(
					env,
					px.Binds_(TypeAs(LIST), px.Eof),
					args,
				)
				if err != nil {
					return nil, err
				}
				lisp := params[0]
				return (lisp.(List))[1:], nil
			}
		},
		// atom while true both lisp atom or go value
		"atom": func(env Env) Element {
			return func(args ...interface{}) (interface{}, error) {
				arg := args[0]
				if l, ok := arg.(List); ok {
					return len(l) == 0, nil
				}
				return true, nil
			}
		},
		// 照搬 cons 运算符对于 golang 嵌入没有足够的收益，这里的 concat 是一个 cons 的变形，
		// 它总是返回包含所有参数的 List 。
		"concat": func(env Env) Element {
			return func(args ...interface{}) (interface{}, error) {
				return List(args), nil
			}
		},
	},
}

func set(env Env, slot, arg interface{}) (interface{}, error) {
	switch setter := slot.(type) {
	case Atom:
		err := env.Setvar(setter.Name, arg)
		if err == nil {
			return nil, err
		}
		return arg, nil
	case Bracket:
		return setter.SetItemBy(env, arg)
	case List:
		s, err := Eval(env, setter)
		if err != nil {
			return nil, err
		}
		return set(env, s, arg)
	default:
		return arg, fmt.Errorf("set error: set %v(%v) as %v is invalid",
			slot, reflect.TypeOf(slot), arg)
	}
}
