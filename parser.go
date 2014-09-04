package gisp

import (
	"fmt"
	p "github.com/Dwarfartisan/goparsec"
)

// Gisp 实现一个基本的 gisp 解释器
type Gisp struct {
	Meta    map[string]interface{}
	Content map[string]interface{}
}

// NewGisp 给定若干可以组合的基准环境，用于构造环境
func NewGisp(buildins map[string]Toolbox) (*Gisp, error) {
	ret := Gisp{
		Meta: map[string]interface{}{
			"category": "gisp",
			"buildins": buildins,
		},
		Content: map[string]interface{}{},
	}
	return &ret, nil
}

// Defvar 实现 Env.Defvar
func (gisp Gisp) Defvar(name string, slot Var) error {
	if _, ok := gisp.Content[name]; ok {
		return fmt.Errorf("var %s exists", name)
	}
	gisp.Content[name] = slot
	return nil
}

// Defun 实现 Env.Defun
func (gisp *Gisp) Defun(name string, functor Functor) error {
	if s, ok := gisp.Local(name); ok {
		switch slot := s.(type) {
		case Func:
			slot.Overload(functor)
		case Var:
			return fmt.Errorf("%s defined as a var")
		default:
			return fmt.Errorf("exists name %s isn't Expr", name)
		}
	}
	gisp.Content[name] = &Function{
		Atom{name, Type{ANY, false}},
		gisp,
		[]Functor{functor},
	}
	return nil
}

// Set 实现 Env.Set 接口
func (gisp *Gisp) Setvar(name string, value interface{}) error {
	if s, ok := gisp.Content[name]; ok {
		switch slot := s.(type) {
		case Var:
			slot.Set(value)
			return nil
		case Function:
			return fmt.Errorf("%v is a Expr", name)
		default:
			return fmt.Errorf("%v is't a var canbe set", name)
		}
	} else {
		return fmt.Errorf("Setable var %s not found", name)
	}
}

func (gisp Gisp) Local(name string) (interface{}, bool) {
	if value, ok := gisp.Content[name]; ok {
		return value, true
	} else {
		return nil, false
	}
}

func (gisp Gisp) Lookup(name string) (interface{}, bool) {
	if value, ok := gisp.Local(name); ok {
		return value, true
	} else {
		return gisp.Global(name)
	}
}

// look up in buildins
func (gisp Gisp) Global(name string) (interface{}, bool) {
	buildins := gisp.Meta["buildins"].(map[string]Toolbox)
	for _, env := range buildins {
		if v, ok := env.Lookup(name); ok {
			return v, true
		}
	}
	return nil, false
}

func (gisp *Gisp) Parse(code string) (interface{}, error) {
	st := p.MemoryParseState(code)

	value, err := ValueParser(st)
	if err != nil {
		return nil, err
	}
	switch lisp := value.(type) {
	case Lisp:
		return lisp.Eval(gisp)
	default:
		return lisp, nil
	}
}

func (gisp *Gisp) Eval(lisps ...interface{}) (interface{}, error) {
	var ret interface{}
	var err error
	for _, l := range lisps {
		switch lisp := l.(type) {
		case Lisp:
			ret, err = lisp.Eval(gisp)
			if err != nil {
				return nil, err
			}
		default:
			ret = lisp
		}
	}
	return ret, nil
}
