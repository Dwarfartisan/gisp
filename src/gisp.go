package gisp

type Lisp interface {
	Eval(env Env) (interface{}, error)
}

type Env interface {
	Define(name string, value interface{}) error
	SetVar(name string, value interface{}) error
	Lookup(name string) (interface{}, bool)
	Local(name string) (interface{}, bool)
	Global(name string) (interface{}, bool)
}

type element func(args ...interface{}) (interface{}, error)
type function func(Env) element
