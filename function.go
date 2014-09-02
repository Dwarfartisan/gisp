package gisp

import (
	"fmt"
)

type Function struct {
	atom    Atom
	Global  Env
	Content []Lambda
}

func (fun Function) Name() string {
	return fun.atom.Name
}

func (fun Function) Task(args ...interface{}) (*Task, error) {
	for _, lambda := range fun.Content {
		task, err := lambda.Task(args...)
		if err == nil {
			return task, nil
		}
	}
	if f, ok := fun.Global.Lookup(fun.Name()); ok {
		if foo := f.(Function); ok {
			return foo.Task(args...)
		}
	}

	return nil, fmt.Errorf("not found args type sign for %v", args)
}

func (fun *Function) Overload(lambda ...Lambda) error {
	fun.Content = append(lambda, fun.Content...)
	return nil
}

func Defun(env Env, funName Atom, lambda Lambda) (*Function, error) {
	fun := Function{funName, env, []Lambda{lambda}}
	err := env.Defun(funName.Name, fun)
	return &fun, err
}

func DefunExpr(env Env) element {
	return func(args ...interface{}) (interface{}, error) {
		funName := args[0].(Atom)
		_args := args[1].(List)
		lambda, err := DeclareLambda(env, _args, args[2:]...)
		if err != nil {
			return nil, err
		}
		if f, ok := env.Local(funName.Name); ok {
			if fun, ok := f.(Function); ok {
				err := fun.Overload(*lambda)
				if err == nil {
					return fun, nil
				} else {
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("%v is defined as no function", funName.Name)
			}
		} else {
			ret, err := Defun(env, funName, *lambda)
			if err == nil {
				return *ret, nil
			} else {
				return nil, err
			}
		}
	}
}
