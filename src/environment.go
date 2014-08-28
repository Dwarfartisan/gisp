package gisp

// Environment 实现了一个基本环境
type Environment struct {
	Meta    map[string]interface{}
	Content map[string]function
}

// Lookup 实现了基本的 Env.Lookup 策略：现在 Local 中查找，否则向上查找 Gobal
func (env Environment) Lookup(name string) (interface{}, bool) {
	if v, ok := env.Local(name); ok {
		return v, true
	}
	return env.Global(name)

}

// Local 实现 Env.Local
func (env Environment) Local(name string) (interface{}, bool) {
	if v, ok := env.Content[name]; ok {
		return v, true
	}
	return nil, false

}

// Global 实现 Env.Global 。如果 Meta 中没有注册 global ，视作顶层环境，返回notfound
func (env Environment) Global(name string) (interface{}, bool) {
	if o, ok := env.Meta["global"]; ok {
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
	case Float, Int:
		return o, nil
	default:
		return lisp, nil
	}
}
