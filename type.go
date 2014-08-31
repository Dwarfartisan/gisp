package gisp

import (
	p "github.com/Dwarfartisan/goparsec"
	"reflect"
)

type Type struct {
	reflect.Type
	option bool
}

func (typ Type) String() string {
	str := typ.Type.String()
	if typ.option {
		return str + "?"
	} else {
		return str
	}
}

func (typ Type) Option() bool {
	return typ.option
}

func TypeParser(st p.ParseState) (interface{}, error) {
	t, err := p.Bind_(p.String("::"),
		p.Choice(
			p.Bind_(p.String("bool"), p.Return(BOOL)),
			p.Bind_(p.String("float"), p.Return(FLOAT)),
			p.Bind_(p.String("int"), p.Return(INT)),
			p.Bind_(p.String("string"), p.Return(STRING)),
			p.Bind_(p.String("any"), p.Return(ANY)),
			p.Bind_(p.String("atom"), p.Return(ATOM)),
			p.Bind_(p.String("quote"), p.Return(QUOTE)),
		))(st)
	if err != nil {
		return nil, err
	}
	_, err = p.Try(p.Rune('?'))(st)
	option := err == nil
	return Type{t.(reflect.Type), option}, nil
}
