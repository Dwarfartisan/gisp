package gisp

import (
	"fmt"
)

type TypeSignError struct {
	Type  Type
	Value interface{}
}

func (err TypeSignError) Error() string {
	return fmt.Sprintf("%v can't match %v", err.Value, err.Type)
}

type Func interface {
	Functor
	Name() string
	Overload(functor ...Functor) error
	Content() []Functor
}

type Function struct {
	atom    Atom
	Global  Env
	content []Functor
}

func (fun Function) Name() string {
	return fun.atom.Name
}

func (fun Function) Task(args ...interface{}) (Lisp, error) {
	for _, functor := range fun.content {
		task, err := functor.Task(args...)
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

func (fun *Function) Overload(functor ...Functor) error {
	fun.content = append(functor, fun.content...)
	return nil
}

func (fun Function) Content() []Functor {
	return fun.content
}

func Defun(env Env, funName Atom, functor Functor) (*Function, error) {
	fun := Function{funName, env, []Functor{functor}}
	err := env.Defun(&fun)
	return &fun, err
}

func DefunExpr(env Env) Element {
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
				return nil, fmt.Errorf("%v is defined as no Expr", funName.Name)
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
