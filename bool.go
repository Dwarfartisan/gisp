package gisp

import (
	"fmt"
	p "github.com/Dwarfartisan/goparsec"
)

type Bool bool

// BoolParser 解析 bool
var BoolParser = p.Bind(p.Choice(p.String("true"), p.String("false")), func(input interface{}) p.Parser {
	return func(st p.ParseState) (interface{}, error) {
		switch input.(string) {
		case "true":
			return Bool(true), nil
		case "false":
			return Bool(false), nil
		default:
			return nil, fmt.Errorf("Unexcept bool token %v", input)
		}
	}
})

// NilParser 解析 nil
var NilParser = p.Bind_(p.String("nil"), p.Return(nil))

type Nil struct {
}

func (n Nil) Eval(env Env) (interface{}, error) {
	return nil, nil
}
