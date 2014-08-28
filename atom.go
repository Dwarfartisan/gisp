package gisp

import (
	"fmt"
	p "github.com/Dwarfartisan/goparsec"
	"reflect"
)

// Atom 类型表达基础的 atom 类型
type Atom struct {
	Name string
	Type reflect.Type
}

func (atom Atom) String() string {
	return fmt.Sprintf("%v::%v", atom.Name, atom.Type)
}

// Eval 方法实现 atom 实例的求值行为
func (atom Atom) Eval(env Env) (interface{}, error) {
	if value, ok := env.Lookup(atom.Name); ok {
		return value, nil
	}
	return nil, fmt.Errorf("value of atom %s not found", atom.Name)
}

func TypeParser(st p.ParseState) (interface{}, error) {
	return p.Bind_(p.String("::"),
		p.Choice(
			p.Bind_(p.String("bool"), p.Return(BOOL)),
			p.Bind_(p.String("float"), p.Return(FLOAT)),
			p.Bind_(p.String("int"), p.Return(INT)),
			p.Bind_(p.String("string"), p.Return(STRING)),
			p.Bind_(p.String("any"), p.Return(ANY)),
			p.Bind_(p.String("atom"), p.Return(ATOM)),
			p.Bind_(p.String("quote"), p.Return(QUOTE)),
		))(st)
}

func AtomParser(st p.ParseState) (interface{}, error) {
	a, err := p.Bind(p.Many1(p.NoneOf("'() \t\r\n.:")),
		p.ReturnString)(st)
	if err != nil {
		return nil, err
	}
	t, err := p.Try(TypeParser)(st)
	if err == nil {
		return Atom{a.(string), t.(reflect.Type)}, nil
	} else {
		return Atom{a.(string), ANY}, nil
	}
}
