package gisp

import (
	"fmt"
)

var Propositions Toolkit = Toolkit{
	Meta: map[string]interface{}{
		"name":     "propositions",
		"category": "environment",
	},
	Content: map[string]function{
		"lambda": LambdaExpr,
		"let":    LetExpr,
		"+":      addExpr,
		"add":    addExpr,
		"-":      subExpr,
		"sub":    subExpr,
		"*":      mulExpr,
		"mul":    mulExpr,
		"/":      divExpr,
		"div":    divExpr,
	},
}

func tofloat64(env Env, arg interface{}) (float64, error) {
	switch value := arg.(type) {
	case float64:
		return value, nil
	case float32:
		return float64(value), nil
	case int:
		return float64(value), nil
	case int8:
		return float64(value), nil
	case int16:
		return float64(value), nil
	case int32:
		return float64(value), nil
	case int64:
		return float64(value), nil
	case Lisp:
		v, err := value.Eval(env)
		if err != nil {
			return 0.0, err
		}
		return tofloat64(env, v)
	default:
		return 0.0, fmt.Errorf("%v isn't avalid number", arg)
	}
}

func addExpr(env Env) element {
	return func(args ...interface{}) (interface{}, error) {
		x := 0.0
		for _, arg := range args {
			v, err := tofloat64(env, arg)
			if err == nil {
				x += v
			} else {
				return nil, err
			}
		}
		return x, nil
	}
}

func subExpr(env Env) element {
	return func(args ...interface{}) (interface{}, error) {
		x, err := tofloat64(env, args[0])
		if err != nil {
			return nil, err
		}
		for _, arg := range args[1:] {
			value, err := tofloat64(env, arg)
			if err != nil {
				return nil, err
			}
			x -= value
		}
		return x, nil
	}
}

func mulExpr(env Env) element {
	return func(args ...interface{}) (interface{}, error) {
		x, err := tofloat64(env, args[0])
		if err != nil {
			return nil, err
		}
		for _, arg := range args[1:] {
			value, err := tofloat64(env, arg)
			if err != nil {
				return nil, err
			}
			x *= value
		}
		return x, nil
	}
}

func divExpr(env Env) element {
	return func(args ...interface{}) (interface{}, error) {
		x, err := tofloat64(env, args[0])
		if err != nil {
			return nil, err
		}
		for _, arg := range args[1:] {
			value, err := tofloat64(env, arg)
			if err != nil {
				return nil, err
			}
			x /= value
		}
		return x, nil
	}
}
