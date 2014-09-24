package gisp

import (
	px "github.com/Dwarfartisan/goparsec/parsex"
)

type Functor interface {
	Task(env Env, args ...interface{}) (Lisp, error)
}

type ArgsSignChecker func(args ...interface{}) error

func ParsexSignChecker(parser px.Parser) ArgsSignChecker {
	return func(args ...interface{}) error {
		st := px.NewStateInMemory(args)
		_, err := parser(st)
		return err
	}
}

type EmptyFunc struct {
	functor func(env Env, args ...interface{}) (Lisp, error)
}

func (ef EmptyFunc) Task(env Env, args ...interface{}) (Lisp, error) {
	lisp, err := ef.functor(env, args...)
	if err != nil {
		return nil, err
	}
	return lisp, nil
}

type SimpleFunc struct {
	tasker func(env Env, args ...interface{}) (Tasker, error)
}

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

type SimpleBox struct {
	checker func(args ...interface{}) error
	task    func(args ...interface{}) Tasker
}

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

type ExprxBox struct {
	TaskerBox
	checker ArgsSignChecker
}

func BoxExprx(asign ArgsSignChecker, expr TaskExpr) ExprxBox {
	return ExprxBox{TaskerBox{expr}, asign}
}

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

type TaskerBox struct {
	functor TaskExpr
}

func BoxExpr(expr TaskExpr) TaskerBox {
	return TaskerBox{expr}
}

func (box TaskerBox) Task(env Env, args ...interface{}) (Lisp, error) {
	task, err := box.functor(env, args...)
	if err != nil {
		return nil, ParsexSignErrorf("Args Type Sign Error: pass %v got error: %v", args, err)
	}
	return TaskBox{task}, nil
}

type EvalBox struct {
	functor LispExpr
}

func EvalExpr(expr LispExpr) EvalBox {
	return EvalBox{expr}
}

func (box EvalBox) Task(env Env, args ...interface{}) (Lisp, error) {
	lisp, err := box.functor(env, args...)
	if err != nil {
		return nil, ParsexSignErrorf("Args Type Sign Error: pass %v got error: %v", args, err)
	}
	return lisp, nil
}
