package gisp

import (
	"fmt"
)

type Let struct {
	Meta    map[string]interface{}
	Content List
}

// let => (let ((a, value), (b, value)...) ...)
func LetExpr(env Env) element {
	return func(args ...interface{}) (interface{}, error) {
		local := map[string]interface{}{}
		vars := args[0].(List)
		for _, v := range vars {
			declares := v.(List)
			name := (declares[0].(Atom)).Name
			value, err := eval(env, (declares[1]))
			if err != nil {
				return nil, err
			}
			local[name] = value
		}
		meta := map[string]interface{}{
			"local": local,
		}
		let := Let{meta, args}
		return let.Eval(env)
	}
}

func (this Let) Define(name string, value interface{}) error {
	if _, ok := this.Local(name); ok {
		return fmt.Errorf("local name %s is exists", name)
	}
	local := this.Meta["local"].(map[string]interface{})
	local[name] = value
	return nil
}

func (this Let) SetVar(name string, value interface{}) error {
	if _, ok := this.Local(name); ok {
		local := this.Meta["local"].(map[string]interface{})
		local[name] = value
		return nil
	} else {
		global := this.Meta["global"].(Env)
		return global.SetVar(name, value)
	}
}

func (this Let) Local(name string) (interface{}, bool) {
	local := this.Meta["local"].(map[string]interface{})
	value, ok := local[name]
	return value, ok
}

func (this Let) Lookup(name string) (interface{}, bool) {
	if value, ok := this.Local(name); ok {
		return value, true
	} else {
		return this.Global(name)
	}
}

func (this Let) Global(name string) (interface{}, bool) {
	global := this.Meta["global"].(Env)
	return global.Lookup(name)
}

func (this Let) Eval(env Env) (interface{}, error) {
	this.Meta["global"] = env
	l := len(this.Content)
	switch l {
	case 0:
		return nil, nil
	case 1:
		return eval(this, this.Content[0])
	default:
		for _, expr := range this.Content[:l-2] {
			eval(this, expr)
		}
		expr := this.Content[l-1]
		return eval(this, expr)
	}
}
