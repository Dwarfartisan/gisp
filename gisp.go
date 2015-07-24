package gisp

// Lisp 基础数据结构封装，Gisp 总是调用这个接口来解析数据
type Lisp interface {
	Eval(env Env) (interface{}, error)
}

// Toolbox 是 gisp 的基础数据结构
type Toolbox interface {
	Lookup(name string) (interface{}, bool)
	Local(name string) (interface{}, bool)
	Global(name string) (interface{}, bool)
}

// Env 是带环境的 Toolbox ，也就是说它有变量作用域
type Env interface {
	Toolbox
	Defvar(name string, slot Var) error
	Setvar(name string, value interface{}) error
	Defun(name string, functor Functor) error
}

// Parser 是解释器的通用接口，用它将解释器与执行的逻辑正交分解。
type Parser interface {
	Parse(string) (interface{}, error)
	Eval(lisps ...interface{}) (interface{}, error)
}

// Tasker 定义了可执行的函数形式
type Tasker func(env Env) (interface{}, error)

// TaskExpr 是带参数的可执行定义， Tasker可以是包装后的 TaskerExpr
type TaskExpr func(env Env, args ...interface{}) (Tasker, error)

// LispExpr 是带参数的 Lisp 执行体
type LispExpr func(env Env, args ...interface{}) (Lisp, error)
