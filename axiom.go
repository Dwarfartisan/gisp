package gisp

import (
	"reflect"
)

// Axiom 是基本的 LISP 公理实现，尽可能贴近原始的 LISP 公理描述，但是部分实现对实际的 golang
// 环境做了妥协
var Axiom = Toolkit{
	Meta: map[string]interface{}{
		"name":     "axiom",
		"category": "environment",
	},
	Content: map[string]function{
		"quote": func(env Env) element {
			return func(args ...interface{}) (interface{}, error) {
				return Quote{args[0]}, nil
			}
		},
		"var": func(env Env) element {
			return func(args ...interface{}) (interface{}, error) {
				first := args[0].(Atom)
				slot := VarSlot(first.Type)
				if len(args) == 1 {
					err := env.Defvar(first.Name, slot)
					return nil, err
				}
				value, err := eval(env, args[1])
				if err != nil {
					return nil, err
				}
				slot.Set(value)
				err = env.Defvar(first.Name, slot)
				return nil, err
			}
		},
		"set": func(env Env) element {
			return func(args ...interface{}) (interface{}, error) {
				value, err := eval(env, args[1])
				if err != nil {
					return nil, err
				}
				err = env.Setvar((args[0].(Atom)).Name, value)
				if err == nil {
					return nil, err
				}
				return value, nil

			}
		},
		"equal": func(env Env) element {
			return func(args ...interface{}) (interface{}, error) {
				x, err := eval(env, args[0])
				if err != nil {
					return nil, err
				}
				y, err := eval(env, args[1])
				if err != nil {
					return nil, err
				}
				return reflect.DeepEqual(x, y), nil
			}
		},
		"cond": func(env Env) element {
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
					result, err := eval(env, cond)
					if err != nil {
						return nil, err
					}
					if ok := result.(bool); ok {
						return eval(env, branch[1])
					}
				}
				// else branch
				if els != nil {
					return eval(env, els)
				}
				return nil, nil
			}
		},
		"car": func(env Env) element {
			return func(args ...interface{}) (interface{}, error) {
				// FIXME: out range error
				lisp, err := eval(env, args[0])
				if err != nil {
					return nil, err
				}
				return (lisp.(List))[0], nil
			}
		},
		"cdr": func(env Env) element {
			return func(args ...interface{}) (interface{}, error) {
				// FIXME: out range error
				lisp, err := eval(env, args[0])
				if err != nil {
					return nil, err
				}
				return (lisp.(List))[1:], nil
			}
		},
		// atom while true both lisp atom or go value
		"atom": func(env Env) element {
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
		"concat": func(env Env) element {
			return func(args ...interface{}) (interface{}, error) {
				return List(args), nil
			}
		},
	},
}
