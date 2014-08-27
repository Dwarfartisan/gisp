package gisp

import (
	"reflect"
)

var Axiom = Environment{
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
				value, err := eval(env, args[1])
				if err != nil {
					return nil, err
				}
				err = env.Define((args[0].(Atom)).Name, value)
				return nil, err
			}
		},
		"set": func(env Env) element {
			return func(args ...interface{}) (interface{}, error) {
				value, err := eval(env, args[1])
				if err != nil {
					return nil, err
				}
				err = env.SetVar((args[0].(Atom)).Name, value)
				if err == nil {
					return nil, err
				} else {
					return value, nil
				}
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
				} else {
					return nil, nil
				}
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
				} else {
					return true, nil
				}
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
