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

// AA 构造一个指定命名的Atom，类型为 ANYOPTION
func AA(name string) Atom {
	return Atom{Name: name, Type: ANYOPTION}
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
		case TaskExpr:
			return slot(env)
		default:
			return slot, nil
		}
	}
	return nil, fmt.Errorf("value of atom %s not found", atom.Name)
}

func atomNameParser(st p.ParseState) (interface{}, error) {
	ret, err := p.Bind(p.Many1(p.NoneOf("'[]() \t\r\n\".:")),
		p.ReturnString)(st)
	if err != nil {
		return nil, err
	}
	test := p.MemoryParseState(ret.(string))
	_, err = p.Bind_(p.Many1(p.Digit), p.Eof)(test)
	if err == nil {
		return nil, fmt.Errorf("atom name can't be a int like %s", ret.(string))
	}
	return ret, nil
}

// AtomParserExt 生成带扩展包的 Atom
func AtomParserExt(env Env) p.Parser {
	return func(st p.ParseState) (interface{}, error) {
		a, err := atomNameParser(st)
		if err != nil {
			return nil, err
		}
		t, err := p.Try(ExtTypeParser(env))(st)
		if err == nil {
			return Atom{a.(string), t.(Type)}, nil
		}
		return Atom{a.(string), ANYMUST}, nil
	}
}

// AtomParser 生成 Atom 对象，但是它不带扩展环境
func AtomParser(st p.ParseState) (interface{}, error) {
	a, err := atomNameParser(st)
	if err != nil {
		return nil, err
	}
	t, err := p.Try(TypeParser)(st)
	if err == nil {
		return Atom{a.(string), t.(Type)}, nil
	}
	return Atom{a.(string), ANYMUST}, nil
}
