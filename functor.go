package gisp

import (
	px "github.com/Dwarfartisan/goparsec/parsex"
)

// 函数的求值总是延迟到一个 TaskBox 中，这样可以在需要的时候异步化

// Functor 定义一个可调用的函子接口
type Functor interface {
	Task(env Env, args ...interface{}) (Lisp, error)
}

// ArgsSignChecker 定义函数签名的验证器类型
type ArgsSignChecker func(args ...interface{}) error

// ParsexSignChecker 定义 Parsex 环境下的函数签名验证器
func ParsexSignChecker(parser px.Parser) ArgsSignChecker {
	return func(args ...interface{}) error {
		st := px.NewStateInMemory(args)
		_, err := parser(st)
		return err
	}
}

// EmptyFunc 定义一个空函数，它用于封装其内含的 functor。
type EmptyFunc struct {
	functor func(env Env, args ...interface{}) (Lisp, error)
}

// Task 实现 EmptyFunc 的求值逻辑，它总是调用其内含的函子
func (ef EmptyFunc) Task(env Env, args ...interface{}) (Lisp, error) {
	lisp, err := ef.functor(env, args...)
	if err != nil {
		return nil, err
	}
	return lisp, nil
}

// SimpleFunc 定一个简单的函数实现，它封装一个 tasker 构造函数
type SimpleFunc struct {
	tasker func(env Env, args ...interface{}) (Tasker, error)
}

// Task 对其内部的tasker求值，返回一个 taskbox
func (sf SimpleFunc) Task(env Env, args ...interface{}) (Lisp, error) {
	params, err := Evals(env, args...)
	if err != nil {
		return nil, err
	}
	tasker, err := sf.tasker(env, params...)
	if err != nil {
		return nil, err
	}
	return TaskBox{tasker}, nil
}

// SimpleBox 实现一个带参数类型检查的 Box
type SimpleBox struct {
	checker func(args ...interface{}) error
	task    func(args ...interface{}) Tasker
}

// Task 方法 实现 SimpleBox 的解析求值
func (sb SimpleBox) Task(env Env, args ...interface{}) (Lisp, error) {
	params, err := Evals(env, args...)
	if err != nil {
		return nil, err
	}
	err = sb.checker(params...)
	if err != nil {
		return nil, err
	}
	return TaskBox{sb.task(params...)}, nil
}

// ExprxBox 是一个带参数类型检查的 TaskerBox
type ExprxBox struct {
	TaskerBox
	checker ArgsSignChecker
}

// BoxExprx 构造一个 ExprxBox
func BoxExprx(asign ArgsSignChecker, expr TaskExpr) ExprxBox {
	return ExprxBox{TaskerBox{expr}, asign}
}

// Task 实现了 ExprxBox 的求值行为
func (box ExprxBox) Task(env Env, args ...interface{}) (Lisp, error) {
	params, err := Evals(env, args...)
	if err != nil {
		return nil, err
	}
	err = box.checker(params...)
	if err != nil {
		return nil, ParsexSignErrorf("Args Type Sign Error: pass %v got error: %v", args, err)
	}

	return box.TaskerBox.Task(env, args...)
}

// TaskerBox 实现 TasskExpr 的封装
type TaskerBox struct {
	functor TaskExpr
}

// BoxExpr 定义一个新的  TaskerBox
func BoxExpr(expr TaskExpr) TaskerBox {
	return TaskerBox{expr}
}

// Task 定义了 TaskerBox 的求值行为
func (box TaskerBox) Task(env Env, args ...interface{}) (Lisp, error) {
	task, err := box.functor(env, args...)
	if err != nil {
		return nil, ParsexSignErrorf("Args Type Sign Error: pass %v got error: %v", args, err)
	}
	return TaskBox{task}, nil
}

// EvalBox 定义了对一个 LispExpr 的求值封装
type EvalBox struct {
	functor LispExpr
}

// EvalExpr 构造一个新的 EvalBox
func EvalExpr(expr LispExpr) EvalBox {
	return EvalBox{expr}
}

// Task 实现了 EvalBox 的求值行为，它返回的是被封装的 LispExpr 的求值结果
func (box EvalBox) Task(env Env, args ...interface{}) (Lisp, error) {
	lisp, err := box.functor(env, args...)
	if err != nil {
		return nil, ParsexSignErrorf("Args Type Sign Error: pass %v got error: %v", args, err)
	}
	return lisp, nil
}
