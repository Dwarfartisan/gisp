package gisp

import (
	"fmt"
)

// Lambda 实现基本的 Lambda 行为
type Lambda struct {
	Meta    map[string]interface{}
	Content List
}

// DeclareLambda 构造 Lambda 表达式 (lambda (args...) body)
func DeclareLambda(env Env, args List, lisps ...interface{}) (*Lambda, error) {
	ret := Lambda{map[string]interface{}{
		"category":          "lambda",
		"formal parameters": args,
		"local":             map[string]interface{}{},
	}, List{}}
	prepare := map[string]bool{}
	for _, lisp := range lisps {
		err := ret.prepare(env, prepare, lisp)
		if err != nil {
			return nil, err
		}
	}
	return &ret, nil
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

// create a lambda s-expr can be eval
func (lambda Lambda) Call(args ...interface{}) Task {
	meta := map[string]interface{}{}
	for k, v := range lambda.Meta {
		meta[k] = v
	}
	meta["actual parameters"] = List(args)
	meta["my"] = map[string]interface{}{}
	l := len(lambda.Content)
	content := make([]interface{}, l)
	for idx, data := range lambda.Content {
		content[idx] = data
	}
	return Task{meta, content}
}

func LambdaExpr(env Env) element {
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

type Task struct {
	Meta    map[string]interface{}
	Content []interface{}
}

func (lambda Task) Local(name string) (interface{}, bool) {
	my := lambda.Meta["my"].(map[string]interface{})
	if value, ok := my[name]; ok {
		return value, true
	}
	if value, ok := lambda.Parameter(name); ok {
		return value, true
	}

	local := lambda.Meta["local"].(map[string]interface{})
	value, ok := local[name]
	return value, ok
}

func (lambda Task) Parameter(name string) (interface{}, bool) {
	formals := lambda.Meta["formal parameters"].(List)
	actuals := lambda.Meta["actual parameters"].(List)
	for idx := range formals {
		formal := formals[idx].(Atom)
		if formal.Name == name {
			return actuals[idx], true
		}
	}
	return nil, false
}

func (lambda Task) Global(name string) (interface{}, bool) {
	global := lambda.Meta["global"].(Env)
	return global.Lookup(name)
}

func (lambda Task) Lookup(name string) (interface{}, bool) {
	if value, ok := lambda.Local(name); ok {
		return value, true
	} else {
		return lambda.Global(name)
	}
}

func (lambda Task) Set(name string, value interface{}) error {
	mine := lambda.Meta["my"].(map[string]Var)
	if _, ok := mine[name]; ok {
		mine[name].Set(value)
		return nil
	} else {
		local := lambda.Meta["local"].(map[string]Var)
		if _, ok := local[name]; ok {
			local[name].Set(value)
			return nil
		} else {
			global := lambda.Meta["global"].(Env)
			return global.Set(name, value)
		}
	}
}

func (lambda Task) Define(name string, slot Var) error {
	mine := lambda.Meta["my"].(map[string]Var)
	if _, ok := mine[name]; ok {
		return fmt.Errorf("%s was exists.", name)
	} else {
		mine[name] = slot
		return nil
	}
}

func (lambda Task) Eval(env Env) (interface{}, error) {
	args := lambda.Meta["actual parameters"].(List)
	actual := make(List, len(args))
	for idx, arg := range args {
		a, err := eval(env, arg)
		if err == nil {
			actual[idx] = a
		} else {
			return nil, err
		}
	}
	lambda.Meta["actual parameters"] = args
	lambda.Meta["global"] = env
	l := len(lambda.Content)
	switch l {
	case 0:
		return nil, nil
	case 1:
		return eval(lambda, lambda.Content[0])
	default:
		for _, expr := range lambda.Content[:l-2] {
			_, err := eval(lambda, expr)
			if err != nil {
				return nil, err
			}
		}
		return eval(lambda, lambda.Content[l-1])
	}
}
