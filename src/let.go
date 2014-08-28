package gisp

import (
	"fmt"
)

// Let 实现 let 环境
type Let struct {
	Meta    map[string]interface{}
	Content List
}

// LetExpr 将 let => (let ((a, value), (b, value)...) ...) 形式构造为一个 let 环境
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

// Define 实现 Env.Define
func (let Let) Define(name string, value interface{}) error {
	if _, ok := let.Local(name); ok {
		return fmt.Errorf("local name %s is exists", name)
	}
	local := let.Meta["local"].(map[string]interface{})
	local[name] = value
	return nil
}

// SetVar 实现 Env.SetVar
func (let Let) SetVar(name string, value interface{}) error {
	if _, ok := let.Local(name); ok {
		local := let.Meta["local"].(map[string]interface{})
		local[name] = value
		return nil
	}
	global := let.Meta["global"].(Env)
	return global.SetVar(name, value)

}

// Local 实现 Env.Local
func (let Let) Local(name string) (interface{}, bool) {
	local := let.Meta["local"].(map[string]interface{})
	value, ok := local[name]
	return value, ok
}

// Lookup 实现 Env.Lookup
func (let Let) Lookup(name string) (interface{}, bool) {
	if value, ok := let.Local(name); ok {
		return value, true
	}
	return let.Global(name)

}

// Global 实现 Env.Global
func (let Let) Global(name string) (interface{}, bool) {
	global := let.Meta["global"].(Env)
	return global.Lookup(name)
}

// Eval 实现 Lisp.Eval
func (let Let) Eval(env Env) (interface{}, error) {
	let.Meta["global"] = env
	l := len(let.Content)
	switch l {
	case 0:
		return nil, nil
	case 1:
		return eval(let, let.Content[0])
	default:
		for _, expr := range let.Content[:l-2] {
			eval(let, expr)
		}
		expr := let.Content[l-1]
		return eval(let, expr)
	}
}
