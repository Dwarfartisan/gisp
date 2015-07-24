package gisp

import (
	"fmt"
)

// Task 定义了可执行的通用结构
type Task struct {
	Meta    map[string]interface{}
	Content []interface{}
}

// Local 定义了 task 本地环境中的变量查找
func (task Task) Local(name string) (interface{}, bool) {
	my := task.Meta["my"].(map[string]Var)
	if slot, ok := my[name]; ok {
		return slot.Get(), true
	}
	if value, ok := task.ParameterValue(name); ok {
		return value, true
	}

	local := task.Meta["local"].(map[string]interface{})
	value, ok := local[name]
	return value, ok
}

// ParameterValue 获取指定参数
func (task Task) ParameterValue(name string) (interface{}, bool) {
	formals := task.Meta["formal parameters"].(List)
	actuals := task.Meta["actual parameters"].([]interface{})
	lastIdx := len(actuals) - 1
	for idx := range formals {
		formal := formals[idx].(Atom)
		if formal.Name == name {
			slot := actuals[idx]
			if idx == lastIdx {
				isVariadic := task.Meta["is variadic"].(bool)
				if isVariadic {
					slots := actuals[lastIdx].([]Var)
					value := make([]interface{}, len(slots))
					for idx, slot := range slots {
						value[idx] = slot.Get()
					}
					return value, true
				}
			}
			return slot.(Var).Get(), true
		}
	}
	return nil, false
}

// Global 在外部环境中查找
func (task Task) Global(name string) (interface{}, bool) {
	global := task.Meta["global"].(Env)
	return global.Lookup(name)
}

// Lookup 执行自内而外的查找
func (task Task) Lookup(name string) (interface{}, bool) {
	if value, ok := task.Local(name); ok {
		return value, true
	}
	return task.Global(name)
}

// Setvar 实现赋值行为
func (task Task) Setvar(name string, value interface{}) error {
	mine := task.Meta["my"].(map[string]Var)
	if _, ok := mine[name]; ok {
		mine[name].Set(value)
		return nil
	}
	local := task.Meta["local"].(map[string]Var)
	if _, ok := local[name]; ok {
		local[name].Set(value)
		return nil
	}
	return fmt.Errorf("can't found var named %s", name)
}

// Defvar 定义新变量
func (task Task) Defvar(name string, slot Var) error {
	mine := task.Meta["my"].(map[string]Var)
	if _, ok := mine[name]; ok {
		return fmt.Errorf("%s was exists", name)
	}
	mine[name] = slot
	return nil
}

// Defun 实现 Env.Defun
func (task Task) Defun(name string, functor Functor) error {
	if s, ok := task.Local(name); ok {
		switch slot := s.(type) {
		case Func:
			slot.Overload(functor)
		case Var:
			return fmt.Errorf("%s defined as a var", name)
		default:
			return fmt.Errorf("exists name %s isn't Expr", name)
		}
	}
	my := task.Meta["my"].(map[string]interface{})
	my[name] = NewFunction(name, task, functor)
	return nil
}

// Eval 实现求值逻辑
func (task Task) Eval(env Env) (interface{}, error) {
	formals := task.Meta["formal parameters"].(List)
	actuals := task.Meta["actual parameters"].([]interface{})
	values := make([]Var, len(actuals))
	for idx, atom := range formals {
		formal := atom.(Atom)
		slot := VarSlot(formal.Type)
		val, err := Eval(task, actuals[idx])
		if err != nil {
			return nil, err
		}
		slot.Set(val)
		values[idx] = slot
	}
	task.Meta["actual values"] = values

	task.Meta["global"] = env
	l := len(task.Content)
	switch l {
	case 0:
		return nil, nil
	case 1:
		return Eval(task, task.Content[0])
	default:
		for _, Expr := range task.Content[:l-2] {
			_, err := Eval(task, Expr)
			if err != nil {
				return nil, err
			}
		}
		return Eval(task, task.Content[l-1])
	}
}
