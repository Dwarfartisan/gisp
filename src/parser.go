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
func NewGisp(buildins map[string]Environment) (*Gisp, error) {
	ret := Gisp{
		Meta: map[string]interface{}{
			"category": "gisp",
			"buildins": buildins,
		},
		Content: map[string]interface{}{},
	}
	return &ret, nil
}

// Define 实现 Env.Define
func (gisp Gisp) Define(name string, value interface{}) error {
	if _, ok := gisp.Content[name]; ok {
		return fmt.Errorf("var %s exists", name)
	}
	gisp.Content[name] = value
	return nil
}

// Set 实现 Env.Set 接口
func (gisp Gisp) SetVar(name string, value interface{}) error {
	if _, ok := gisp.Content[name]; ok {
		gisp.Content[name] = value
		return nil
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
	if value, ok := gisp.Global(name); ok {
		return value, true
	} else {
		return gisp.Local(name)
	}
}

// look up in buildins
func (gisp Gisp) Global(name string) (interface{}, bool) {
	buildins := gisp.Meta["buildins"].(map[string]Environment)
	for _, env := range buildins {
		if v, ok := env.Lookup(name); ok {
			return v, true
		}
	}
	return nil, false
}

func (gisp Gisp) Parse(code string) (interface{}, error) {
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
