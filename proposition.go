package gisp

import (
	px "github.com/Dwarfartisan/goparsec/parsex"
)

var Propositions Toolkit = Toolkit{
	Meta: map[string]interface{}{
		"name":     "propositions",
		"category": "environment",
	},
	Content: map[string]Expr{
		"lambda": LambdaExpr,
		"let":    LetExpr,
		"+":      addExpr,
		"add":    addExpr,
		"-":      subExpr,
		"sub":    subExpr,
		"*":      mulExpr,
		"mul":    mulExpr,
		"/":      divExpr,
		"div":    divExpr,
	},
}

func addExpr(env Env) Element {
	return func(args ...interface{}) (interface{}, error) {
		st := px.NewStateInMemory(args)
		return addx(st)
	}
}

func subExpr(env Env) Element {
	return func(args ...interface{}) (interface{}, error) {
		st := px.NewStateInMemory(args)
		return subx(st)
	}
}

func mulExpr(env Env) Element {
	return func(args ...interface{}) (interface{}, error) {
		st := px.NewStateInMemory(args)
		return mulx(st)
	}
}

func divExpr(env Env) Element {
	return func(args ...interface{}) (interface{}, error) {
		st := px.NewStateInMemory(args)
		return divx(st)
	}
}
