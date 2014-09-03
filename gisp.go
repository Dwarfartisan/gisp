package gisp

type Lisp interface {
	Eval(env Env) (interface{}, error)
}

type Functor interface {
	Task(args ...interface{}) (Lisp, error)
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
	Defun(fun Func) error
}

type Element func(args ...interface{}) (interface{}, error)
type function func(Env) Element
