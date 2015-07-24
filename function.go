package gisp

import (
	"fmt"
)

// TypeSignError 定义签名错误
type TypeSignError struct {
	Type  Type
	Value interface{}
}

func (err TypeSignError) Error() string {
	return fmt.Sprintf("%v can't match %v", err.Value, err.Type)
}

// ParsexSignError 定义了在解析过程中遇到的签名错误
type ParsexSignError struct {
	message string
	args    []interface{}
}

// ParsexSignErrorf 方法构造指定的 ParsexSignError
func ParsexSignErrorf(message string, args ...interface{}) ParsexSignError {
	return ParsexSignError{message, args}
}

func (err ParsexSignError) Error() string {
	return fmt.Sprintf(err.message, err.args...)
}

// Func 接口定义可以 Overload 的 Gisp 函数类型
type Func interface {
	Functor
	Name() string
	Overload(functor Functor) error
	Content() []Functor
}

// TaskBox 定义了可执行的 task 结构
type TaskBox struct {
	task func(env Env) (interface{}, error)
}

// Eval 实现了 Task 的求值行为
func (tb TaskBox) Eval(env Env) (interface{}, error) {
	return tb.task(env)
}

// Function 定义了 Gisp 函数实现
type Function struct {
	atom    Atom
	Global  Env
	content []Functor
}

// NewFunction 构造一个新的 Function 对象
func NewFunction(name string, global Env, functor Functor) *Function {
	return &Function{
		atom:    Atom{name, Type{ANY, false}},
		Global:  global,
		content: []Functor{functor},
	}
}

// Name 给出函数名
func (fun Function) Name() string {
	return fun.atom.Name
}

// Task 实现了 Function 对象的求值逻辑
func (fun Function) Task(env Env, args ...interface{}) (Lisp, error) {
	for _, functor := range fun.content {
		task, err := functor.Task(env, args...)
		if err == nil {
			return task, nil
		}
	}

	if f, ok := fun.Global.Global(fun.Name()); ok {
		switch foo := f.(type) {
		case Functor:
			return foo.Task(env, args...)
		case TaskExpr:
			task, err := foo(env, args...)
			if err != nil {
				return nil, err
			}
			return TaskBox{task}, nil
		case LispExpr:
			lisp, err := foo(env, args...)
			if err != nil {
				return nil, err
			}
			return lisp, nil
		}
	}
	return nil, fmt.Errorf("not found args type sign for %v", args)
}

// Overload 实现了 Function 的 Overload 行为
func (fun *Function) Overload(functor Functor) error {
	fun.content = append([]Functor{functor}, fun.content...)
	return nil
}

// Content 返回函数体
func (fun Function) Content() []Functor {
	return fun.content
}

// DefunExpr 是构造 Function 的表达式
func DefunExpr(env Env, args ...interface{}) (Tasker, error) {
	funName := args[0].(Atom)
	_args := args[1].(List)
	var err error
	lambda, err := DeclareLambda(env, _args, args[2:]...)
	if err != nil {
		return nil, err
	}
	if f, ok := env.Local(funName.Name); ok {
		if fun, ok := f.(Function); ok {
			err := fun.Overload(*lambda)
			if err == nil {
				return Q(fun).Eval, nil
			}
			return nil, err
		}
		return nil, fmt.Errorf("%v is defined as no Expr", funName.Name)
	}

	err = env.Defun(funName.Name, *lambda)
	if err == nil {
		return nil, nil
	}
	return nil, err
}
