package gisp

import (
	"fmt"
	p "github.com/Dwarfartisan/goparsec"
	"reflect"
)

// Gisp 实现一个基本的 gisp 解释器
type Gisp struct {
	Meta    map[string]interface{}
	Content map[string]interface{}
}

// NewGisp 给定若干可以组合的基准环境，用于构造环境
func NewGisp(buildins map[string]Toolbox) *Gisp {
	ret := Gisp{
		Meta: map[string]interface{}{
			"category": "gisp",
			"buildins": buildins,
		},
		Content: map[string]interface{}{},
	}
	return &ret
}

func NewGispWith(buildins map[string]Toolbox, ext map[string]Toolbox) *Gisp {
	gisp := NewGisp(buildins)
	if ext == nil {
		return gisp
	}
	for k, v := range ext {
		gisp.DefAs(k, v)
	}
	return gisp
}

// def as = def var + set var
func (gisp *Gisp) DefAs(name string, value interface{}) error {
	t := Type{reflect.TypeOf(value), false}
	slot := VarSlot(t)
	slot.Set(value)
	return gisp.Defvar(name, slot)
}

// def option as  = def var? + set var
func (gisp *Gisp) DefOptAs(name string, value interface{}) error {
	t := Type{reflect.TypeOf(value), true}
	slot := VarSlot(t)
	slot.Set(value)
	return gisp.Defvar(name, slot)
}

// Defvar 实现 Env.Defvar
func (gisp *Gisp) Defvar(name string, slot Var) error {
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
		if slot, ok := value.(Var); ok {
			return slot.Get(), true
		}
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
	var v interface{}
	var e error
	for {
		Skip(st)
		_, err := p.Eof(st)
		if err == nil {
			break
		}
		value, err := ValueParserExt(gisp)(st)
		if err != nil {
			return nil, err
		}
		switch lisp := value.(type) {
		case Lisp:
			 v, e = lisp.Eval(gisp)
		default:
			 v = lisp
			 e = nil
		}
	}
	return v, e
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
