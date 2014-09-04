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
func LetExpr(env Env) Element {
	return func(args ...interface{}) (interface{}, error) {
		local := map[string]Var{}
		vars := args[0].(List)
		for _, v := range vars {
			declares := v.(List)
			varb := declares[0].(Atom)
			slot := VarSlot(varb.Type)
			value, err := eval(env, (declares[1]))
			if err != nil {
				return nil, err
			}
			slot.Set(value)
			local[varb.Name] = slot
		}
		meta := map[string]interface{}{
			"local": local,
		}
		let := Let{meta, args}
		return let.Eval(env)
	}
}

// Defvar 实现 Env.Defvar
func (let Let) Defvar(name string, slot Var) error {
	if _, ok := let.Local(name); ok {
		return fmt.Errorf("local name %s is exists", name)
	}
	local := let.Meta["local"].(map[string]interface{})
	local[name] = slot
	return nil
}

// Defun 实现 Env.Defun
func (let Let) Defun(name string, functor Functor) error {
	if s, ok := let.Local(name); ok {
		switch slot := s.(type) {
		case Func:
			slot.Overload(functor)
		case Var:
			return fmt.Errorf("%s defined as a var", name)
		default:
			return fmt.Errorf("exists name %s isn't Expr", name)
		}
	}
	local := let.Meta["local"].(map[string]interface{})
	local[name] = NewFunction(name, let, functor)
	return nil
}

// Setvar 实现 Env.Setvar
func (let Let) Setvar(name string, value interface{}) error {
	if _, ok := let.Local(name); ok {
		local := let.Meta["local"].(map[string]Var)
		local[name].Set(value)
		return nil
	}
	global := let.Meta["global"].(Env)
	return global.Setvar(name, value)
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
		for _, Expr := range let.Content[:l-2] {
			eval(let, Expr)
		}
		Expr := let.Content[l-1]
		return eval(let, Expr)
	}
}
