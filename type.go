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

func stop(st p.ParseState) (interface{}, error) {
	pos := st.Pos()
	defer st.SeekTo(pos)
	r, err := p.Choice(
		p.Try(p.Space),
		p.Try(p.NewLine),
		p.Try(p.OneOf(":.()[]{}?")),
		p.Try(p.Eof),
	)(st)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func stopWord(x interface{}) p.Parser {
	return p.Bind_(stop, p.Return(x))
}

func typeName(word string) p.Parser {
	return p.Bind(p.String(word), stopWord)
}

var anyType = p.Bind(p.Bind(p.Many1(p.Either(p.Try(p.Digit), p.Letter)), stopWord), p.ReturnString)

func SliceTypeParserExt(env Env) p.Parser {
	return func(st p.ParseState) (interface{}, error) {
		t, err := p.Bind_(p.String("[]"), ExtTypeParser(env))(st)
		if err != nil {
			return nil, err
		}
		return reflect.SliceOf(t.(Type).Type), nil
	}
}

func MapTypeParserExt(env Env) p.Parser {
	return func(st p.ParseState) (interface{}, error) {
		key, err := p.Between(p.String("map["), p.Rune(']'), ExtTypeParser(env))(st)
		if err != nil {
			return nil, err
		}
		value, err := ExtTypeParser(env)(st)
		if err != nil {
			return nil, err
		}
		return reflect.MapOf(key.(Type).Type, value.(Type).Type), nil
	}
}

func MapTypeParser(st p.ParseState) (interface{}, error) {
	key, err := p.Between(p.String("map["), p.Rune(']'), TypeParser)(st)
	if err != nil {
		return nil, err
	}
	value, err := TypeParser(st)
	if err != nil {
		return nil, err
	}
	return reflect.MapOf(key.(Type).Type, value.(Type).Type), nil
}

func ExtTypeParser(env Env) p.Parser {
	return func(st p.ParseState) (interface{}, error) {
		_, err := p.String("::")(st)
		if err != nil {
			return nil, err
		}
		buildin := p.Choice(
			p.Try(p.Bind_(typeName("bool"), p.Return(BOOL))),
			p.Try(p.Bind_(typeName("float"), p.Return(FLOAT))),
			p.Try(p.Bind_(typeName("int"), p.Return(INT))),
			p.Try(p.Bind_(typeName("string"), p.Return(STRING))),
			p.Try(p.Bind_(typeName("time"), p.Return(TIME))),
			p.Try(p.Bind_(typeName("duration"), p.Return(DURATION))),
			p.Try(p.Bind_(typeName("any"), p.Return(ANY))),
			p.Try(p.Bind_(typeName("atom"), p.Return(ATOM))),
			p.Try(p.Bind_(p.String("list"), p.Return(LIST))),
			p.Try(p.Bind_(typeName("quote"), p.Return(QUOTE))),
			p.Try(p.Bind_(p.String("dict"), p.Return(DICT))),
			p.Try(MapTypeParserExt(env)),
		)
		ext := func(st p.ParseState) (interface{}, error) {
			n, err := anyType(st)
			if err != nil {
				return nil, err
			}
			t, ok := env.Lookup(n.(string))
			if !ok {
				return nil, st.Trap("type %v not found", n)
			}
			if typ, ok := t.(reflect.Type); ok {
				return typ, nil
			} else {
				return nil, st.Trap("var %v is't a type. It is %v", n, reflect.TypeOf(t))
			}
		}
		t, err := p.Either(buildin, ext)(st)
		if err != nil {
			return nil, err
		}
		_, err = p.Try(p.Rune('?'))(st)
		option := err == nil
		return Type{t.(reflect.Type), option}, nil
	}
}

func TypeParser(st p.ParseState) (interface{}, error) {
	t, err := p.Bind_(p.String("::"),
		p.Choice(
			p.Try(p.Bind_(p.String("bool"), p.Return(BOOL))),
			p.Try(p.Bind_(p.String("float"), p.Return(FLOAT))),
			p.Try(p.Bind_(p.String("int"), p.Return(INT))),
			p.Try(p.Bind_(p.String("string"), p.Return(STRING))),
			p.Try(p.Bind_(p.String("time"), p.Return(TIME))),
			p.Try(p.Bind_(p.String("duration"), p.Return(DURATION))),
			p.Try(p.Bind_(p.String("any"), p.Return(ANY))),
			p.Try(p.Bind_(p.String("atom"), p.Return(ATOM))),
			p.Try(p.Bind_(p.String("list"), p.Return(LIST))),
			p.Try(p.Bind_(p.String("quote"), p.Return(QUOTE))),
			p.Try(p.Bind_(p.String("dict"), p.Return(DICT))),
			MapTypeParser,
		))(st)
	if err != nil {
		return nil, err
	}
	_, err = p.Try(p.Rune('?'))(st)
	option := err == nil
	return Type{t.(reflect.Type), option}, nil
}
