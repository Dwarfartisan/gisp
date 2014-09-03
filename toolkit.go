package gisp

// Toolkit 实现了一个基本环境，它没有 Define 和 Set ，用于封装只读的环境。
type Toolkit struct {
	Meta    map[string]interface{}
	Content map[string]Expr
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
	if v, ok := tk.Content[name]; ok {
		return v, true
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

func eval(env Env, lisp interface{}) (interface{}, error) {
	// a lisp data or go value
	switch o := lisp.(type) {
	case Lisp:
		value, err := o.Eval(env)
		return value, err
	case bool:
		return Bool(o), nil
	case float32:
		return Float(o), nil
	case float64:
		return Float(o), nil
	case int8:
		return Int(o), nil
	case int16:
		return Int(o), nil
	case int32:
		return Int(o), nil
	case int64:
		return Int(o), nil
	case Float, Int, Bool, nil:
		return o, nil
	default:
		return lisp, nil
	}
}

func evals(env Env, args ...interface{}) (interface{}, error) {
	data := make([]interface{}, len(args))
	for idx, arg := range args {
		ret, err := eval(env, arg)
		if err != nil {
			return nil, err
		}
		data[idx] = ret
	}
	return data, nil
}
