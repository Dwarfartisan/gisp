package gisp

import (
	"fmt"
	. "github.com/Dwarfartisan/goparsec"
)

type Gisp struct {
	Meta    map[string]interface{}
	Content map[string]interface{}
}

// 给定若干可以组合的基准环境
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

func (this Gisp) Define(name string, value interface{}) error {
	if _, ok := this.Content[name]; ok {
		return fmt.Errorf("var %s exists", name)
	} else {
		this.Content[name] = value
	}
	return nil
}

func (this Gisp) SetVar(name string, value interface{}) error {
	if _, ok := this.Content[name]; ok {
		this.Content[name] = value
		return nil
	} else {
		return fmt.Errorf("Setable var %s not found", name)
	}
}

func (this Gisp) Local(name string) (interface{}, bool) {
	if value, ok := this.Content[name]; ok {
		return value, true
	} else {
		return nil, false
	}
}

func (this Gisp) Lookup(name string) (interface{}, bool) {
	if value, ok := this.Global(name); ok {
		return value, true
	} else {
		return this.Local(name)
	}
}

// look up in buildins
func (this Gisp) Global(name string) (interface{}, bool) {
	buildins := this.Meta["buildins"].(map[string]Environment)
	for _, env := range buildins {
		if v, ok := env.Lookup(name); ok {
			return v, true
		}
	}
	return nil, false
}

func (this Gisp) Parse(code string) (interface{}, error) {
	st := MemoryParseState(code)

	value, err := ValueParser(st)
	if err != nil {
		return nil, err
	}
	switch lisp := value.(type) {
	case Lisp:
		return lisp.Eval(this)
	default:
		return lisp, nil
	}
}
