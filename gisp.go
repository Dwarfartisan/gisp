package gisp

type Lisp interface {
	Eval(env Env) (interface{}, error)
}

type Functor interface {
	Task(env Env, args ...interface{}) (Lisp, error)
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

type Element func(args ...interface{}) (interface{}, error)
type Expr func(Env) Element

type Gearbox struct {
	Meta    map[string]interface{}
	Content map[string]interface{}
}

// Lookup 实现了基本的 Env.Lookup 策略：现在 Local 中查找，否则向上查找 Gobal
func (gb Gearbox) Lookup(name string) (interface{}, bool) {
	if v, ok := gb.Local(name); ok {
		return v, true
	}
	return gb.Global(name)
}

// Local 实现 Env.Local
func (gb Gearbox) Local(name string) (interface{}, bool) {
	if v, ok := gb.Content[name]; ok {
		return v, true
	}
	return nil, false

}

// Global 实现 Env.Global 。如果 Meta 中没有注册 global ，视作顶层环境，返回notfound
func (gb Gearbox) Global(name string) (interface{}, bool) {
	if o, ok := gb.Meta["global"]; ok {
		outer := o.(Env)
		return outer.Lookup(name)
	}
	return nil, false
}
