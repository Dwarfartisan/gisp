package gisp

type Lisp interface {
	Eval(env Env) (interface{}, error)
}

type Toolbox interface {
	Lookup(name string) (interface{}, bool)
	Local(name string) (interface{}, bool)
	Global(name string) (interface{}, bool)
}

type Env interface {
	Toolbox
	Defvar(name string, slot Var) error
	Setvar(name string, value interface{}) error
	Defun(name string, functor Functor) error
}

type Parser interface {
	Parse(string) (interface{}, error)
	Eval(lisps ...interface{}) (interface{}, error)
}

type Tasker func(env Env) (interface{}, error)
type TaskExpr func(env Env, args ...interface{}) (Tasker, error)
type LispExpr func(env Env, args ...interface{}) (Lisp, error)
