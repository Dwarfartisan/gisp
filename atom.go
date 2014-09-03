package gisp

import (
	"fmt"
	p "github.com/Dwarfartisan/goparsec"
)

// Atom 类型表达基础的 atom 类型
type Atom struct {
	Name string
	Type Type
}

func (atom Atom) String() string {
	return fmt.Sprintf("%v::%v", atom.Name, atom.Type)
}

// Eval 方法实现 atom 实例的求值行为
func (atom Atom) Eval(env Env) (interface{}, error) {
	if s, ok := env.Lookup(atom.Name); ok {
		switch slot := s.(type) {
		case Var:
			value := slot.Get()
			return value, nil
		case Expr:
			return slot(env), nil
		default:
			return slot, nil
		}
	}
	return nil, fmt.Errorf("value of atom %s not found", atom.Name)
}

func AtomParser(st p.ParseState) (interface{}, error) {
	a, err := p.Bind(p.Many1(p.NoneOf("'() \t\r\n.:")),
		p.ReturnString)(st)
	if err != nil {
		return nil, err
	}
	t, err := p.Try(TypeParser)(st)
	if err == nil {
		return Atom{a.(string), t.(Type)}, nil
	} else {
		return Atom{a.(string), Type{ANY, false}}, nil
	}
}
