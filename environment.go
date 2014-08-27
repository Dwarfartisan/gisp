package gisp

type Environment struct {
	Meta    map[string]interface{}
	Content map[string]function
}

func (this Environment) Lookup(name string) (interface{}, bool) {
	if v, ok := this.Local(name); ok {
		return v, true
	} else {
		return this.Global(name)
	}
}

func (this Environment) Local(name string) (interface{}, bool) {
	if v, ok := this.Content[name]; ok {
		return v, true
	} else {
		return nil, false
	}
}

func (this Environment) Global(name string) (interface{}, bool) {
	if o, ok := this.Meta["global"]; ok {
		outer := o.(Env)
		return outer.Lookup(name)
	} else {
		return nil, false
	}
}

func eval(env Env, lisp interface{}) (interface{}, error) {
	// a lisp data or go value
	if l, ok := lisp.(Lisp); ok {
		value, err := l.Eval(env)
		return value, err
	} else {
		return lisp, nil
	}
}
