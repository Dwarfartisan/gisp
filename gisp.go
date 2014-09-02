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
	Defun(name string, fun Function) error
}

type element func(args ...interface{}) (interface{}, error)
type function func(Env) element
