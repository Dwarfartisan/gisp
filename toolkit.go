package gisp

import (
	//"fmt"
	"reflect"
)

// Toolkit 实现了一个基本环境，它没有 Define 和 Set ，用于封装只读的环境。
type Toolkit struct {
	Meta    map[string]interface{}
	Content map[string]interface{}
}

// Lookup 实现了基本的 Env.Lookup 策略：现在 Local 中查找，否则向上查找 Gobal
func (tk Toolkit) Lookup(name string) (interface{}, bool) {
	if v, ok := tk.Local(name); ok {
		return v, true
	}
	return tk.Global(name)
}

// Local 实现 Env.Local
func (tk Toolkit) Local(name string) (interface{}, bool) {
	if value, ok := tk.Content[name]; ok {
		if slot, ok := value.(Var); ok {
			return slot.Get(), true
		}
		return value, true
	}
	return nil, false
}

// Global 实现 Env.Global 。如果 Meta 中没有注册 global ，视作顶层环境，返回notfound
func (tk Toolkit) Global(name string) (interface{}, bool) {
	if o, ok := tk.Meta["global"]; ok {
		outer := o.(Env)
		return outer.Lookup(name)
	}
	return nil, false
}

// Eval 函数在给定环境中执行 lisp 对象进行求值
func Eval(env Env, lisp interface{}) (interface{}, error) {
	// a lisp data or go value
	switch o := lisp.(type) {
	case Lisp:
		value, err := o.Eval(env)
		return value, err
	default:
		return Value(o), nil
	}
}

// Evals 执行多个 lisp 对象
func Evals(env Env, args ...interface{}) ([]interface{}, error) {
	data := make([]interface{}, len(args))
	for idx, arg := range args {
		ret, err := Eval(env, arg)
		if err != nil {
			return nil, err
		}
		data[idx] = ret
	}
	return data, nil
}

// Value 函数对 golang 基本类型做封装，目前仅止将整型转为 Int， 将浮点型转为 Float，
// 将 rune 转为 Rune， 其它不做处理
func Value(x interface{}) interface{} {
	switch v := x.(type) {
	case bool:
		return Bool(v)
	case float32:
		return Float(v)
	case float64:
		return Float(v)
	case int8:
		return Int(v)
	case int16:
		return Int(v)
	case int32:
		return Int(v)
	case int64:
		return Int(v)
	case int:
		return Int(v)
	default:
		return v
	}
}

//InReflects 从反射对象中取值
func InReflects(values []reflect.Value) ([]interface{}, error) {
	res := make([]interface{}, len(values))
	for idx, value := range values {
		if value.CanInterface() {
			res[idx] = value.Interface()
		}
		if !value.IsValid() || value.Kind() == reflect.Ptr && value.IsNil() {
			res[idx] = nil
		}
		res[idx] = value.Interface()
	}
	return res, nil
}
