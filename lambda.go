package gisp

import (
	"fmt"
	p "github.com/Dwarfartisan/goparsec"
	px "github.com/Dwarfartisan/goparsec/parsex"
)

// Lambda 实现基本的 Lambda 行为
type Lambda struct {
	Meta    map[string]interface{}
	Content List
}

// DeclareLambda 构造 Lambda 表达式 (lambda (args...) body)
func DeclareLambda(env Env, args List, lisps ...interface{}) (*Lambda, error) {
	ret := Lambda{map[string]interface{}{
		"category": "lambda",
		"local":    map[string]interface{}{},
	}, List{}}
	ret.prepareArgs(args)
	prepare := map[string]bool{}
	for _, lisp := range lisps {
		err := ret.prepare(env, prepare, lisp)
		if err != nil {
			return nil, err
		}
	}
	return &ret, nil
}

func (lambda *Lambda) prepareArgs(args List) {
	l := len(args)
	formals := make(List, len(args))
	if l == 0 {
		lambda.Meta["parameters parsex"] = []Var{}
		return
	}
	lidx := l - 1
	last := args[l-1].(Atom)
	name := last.Name
	isVariadic := false
	namelen := len(name)
	if namelen >= 4 {
		st := p.MemoryParseState(name[namelen-4:])
		_, err := p.Binds_(p.NoneOf("."), p.String("..."), p.Eof)(st)
		if err == nil {
			isVariadic = true
		}
	}
	lambda.Meta["is variadic"] = isVariadic
	ps := make([]px.Parser, l+1)
	for idx, arg := range args[:lidx] {
		ps[idx] = argParser(arg.(Atom))
		formals[idx] = arg
	}
	if isVariadic {
		ps[lidx] = px.Many(argParser(last))
		larg := Atom{name[namelen-3:], last.Type}
		formals[lidx] = larg
	} else {
		ps[lidx] = argParser(last)
		formals[lidx] = last
	}
	ps[l] = px.Eof
	lambda.Meta["formals parameters"] = formals
	lambda.Meta["parameter parsexs"] = ps
}

func (lambda *Lambda) prepare(env Env, prepare map[string]bool, content interface{}) error {
	next := map[string]bool{}
	for key := range prepare {
		next[key] = true
	}
	var err error
	switch lisp := content.(type) {
	case Atom:
		err = lambda.prepareAtom(env, next, lisp)
		return err
	case List:
		err = lambda.prepareList(env, next, lisp)
	}
	if err == nil {
		lambda.Content = append(lambda.Content, content)
	}
	return err
}

func (lambda Lambda) prepareAtom(env Env, prepare map[string]bool, one Atom) error {
	if _, ok := prepare[one.Name]; ok {
		return nil
	}
	next := map[string]bool{}
	for key := range prepare {
		next[key] = true
	}

	for _, arg := range lambda.Meta["formal parameters"].(List) {
		if (arg.(Atom)).Name == one.Name {
			return nil
		}
	}
	if _, ok := prepare[one.Name]; !ok {
		if v, ok := env.Lookup(one.Name); ok {
			local := (lambda.Meta["local"]).(map[string]interface{})
			local[one.Name] = v
		} else {
			return fmt.Errorf("%s not found", one.Name)
		}
	}
	return nil
}

func (lambda Lambda) prepareList(env Env, prepare map[string]bool, content List) error {
	next := map[string]bool{}
	for key := range prepare {
		next[key] = true
	}
	var err error = nil
	fun := content[0].(Atom)
	switch fun.Name {
	case "var":
		name := content[1].(string)
		if err != nil {
			return err
		} else {
			next[name] = true
		}
	case "lambda":
		args := content[1].(List)
		for _, a := range args {
			arg := a.(Atom)
			next[arg.Name] = true
		}
	case "let":
		for _, def := range content[1].(List) {
			define := def.(List)
			name := define[0].(string)
			next[name] = true
		}
	}
	for _, l := range content {
		switch lisp := l.(type) {
		case List:
			err = lambda.prepareList(env, next, lisp)
		case Atom:
			err = lambda.prepareAtom(env, next, lisp)
		}
	}
	return err
}

// 类型签名
func (lambda Lambda) TypeSign() []Type {
	formals := lambda.Meta["formal parameters"].(List)
	types := make([]Type, len(formals))
	for idx, formal := range formals {
		types[idx] = formal.(Atom).Type
	}
	return types
}

func (lambda Lambda) MatchArgsSign(args ...interface{}) (interface{}, error) {
	pxs := lambda.Meta["parameter parsexs"].([]px.Parser)
	st := px.NewStateInMemory(args)
	return px.Union(pxs...)(st)
}

// create a lambda s-expr can be eval
func (lambda Lambda) Task(args ...interface{}) (Lisp, error) {
	meta := map[string]interface{}{}
	for k, v := range lambda.Meta {
		meta[k] = v
	}
	actuals, err := lambda.MatchArgsSign(args...)
	if err != nil {
		return Nil{}, err
	}
	meta["actual parameters"] = actuals
	meta["my"] = map[string]Var{}
	l := len(lambda.Content)
	content := make([]interface{}, l)
	for idx, data := range lambda.Content {
		content[idx] = data
	}
	return &Task{meta, content}, nil
}

func LambdaExpr(env Env) Element {
	return func(args ...interface{}) (interface{}, error) {
		_args := args[0].(List)
		ret, err := DeclareLambda(env, _args, args[1:]...)
		if err == nil {
			return *ret, nil
		} else {
			return nil, err
		}
	}
}
