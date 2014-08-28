package gisp

import (
	"fmt"
	"reflect"
	"strconv"

	p "github.com/Dwarfartisan/goparsec"
)

// BoolParser 解析 bool
var BoolParser = p.Bind(p.Choice(p.String("true"), p.String("false")), func(input interface{}) p.Parser {
	return func(st p.ParseState) (interface{}, error) {
		switch input.(string) {
		case "true":
			return true, nil
		case "false":
			return false, nil
		default:
			return nil, fmt.Errorf("Unexcept bool token %v", input)
		}
	}
})

// NilParser 解析 nil
var NilParser = p.Bind_(p.String("nil"), p.Return(nil))

// FloatParser 解析浮点数
func FloatParser(st p.ParseState) (interface{}, error) {
	f, err := Try(p.Float)(st)
	if err == nil {
		return Float(strconv.ParseFloat(f.(string), 64))
	}
	return nil, err
}

// IntParser 解析整数
func IntParser(st p.ParseState) (interface{}, error) {
	i, err := p.Int(st)
	if err == nil {
		return Int(strconv.Atoi(i.(string)))
	}
	return nil, err

}

var EscapeChar = p.Bind_(p.Rune('\\'), func(st ParseState) (interface{}, error) {
	r, err := OneOf("nrt\"\\")(st)
	if err == nil {
		ru := r.(rune)
		switch ru {
		case 'r':
			return '\r', nil
		case 'n':
			return '\n', nil
		// FIXME:引号的解析偷懒了，单双引号的应该分开。
		case '\'':
			return '\'', nil
		case '"':
			return '"', nil
		case '\\':
			return '\\', nil
		case 't':
			return '\t', nil
		default:
			return nil, st.Trap("Unknown escape sequence \\%c", r)
		}
	} else {
		return nil, err
	}
})

var RuneParser = Bind(
	Between(Rune('\''), Rune('\''),
		Either(Try(EscapeChar), NoneOf("'"))),
	ReturnString)

var StringParser = Bind(
	Between(Rune('"'), Rune('"'),
		Many(Either(Try(EscapeChar), NoneOf("\"")))),
	ReturnString)

func TypeParser(st ParseState) (interface{}, error) {
	return Bind_(String("::"),
		Choice(
			Bind_(String("bool"), Return(BOOL)),
			Bind_(String("float"), Return(FLOAT)),
			Bind_(String("int"), Return(INT)),
			Bind_(String("string"), Return(STRING)),
			Bind_(String("any"), Return(ANY)),
			Bind_(String("atom"), Return(ATOM)),
			Bind_(String("quote"), Return(QUOTE)),
		))(st)
}

func AtomParser(st ParseState) (interface{}, error) {
	a, err := Bind(Many1(NoneOf("'() \t\r\n.:")),
		ReturnString)(st)
	if err != nil {
		return nil, err
	}
	t, err := Try(TypeParser)(st)
	if err == nil {
		return Atom{a.(string), t.(reflect.Type)}, nil
	} else {
		return Atom{a.(string), ANY}, nil
	}
}

func bodyParser(st ParseState) (interface{}, error) {
	value, err := SepBy(ValueParser, Many1(Space))(st)
	return value, err
}

func ListParser(st ParseState) (interface{}, error) {
	one := Bind(AtomParser, func(atom interface{}) Parser {
		return Bind_(Rune(')'), Return(List{atom}))
	})
	list, err := Either(Try(Bind_(Rune('('), one)),
		Between(Rune('('), Rune(')'), bodyParser))(st)
	if err == nil {
		return List(list.([]interface{})), nil
	} else {
		return nil, err
	}
}

func QuoteParser(st ParseState) (interface{}, error) {
	lisp, err := Bind_(Rune('\''), ValueParser)(st)
	if err == nil {
		return Quote{lisp}, nil
	} else {
		return nil, err
	}
}

func ValueParser(st ParseState) (interface{}, error) {
	value, err := Choice(StringParser,
		IntParser,
		FloatParser,
		QuoteParser,
		RuneParser,
		StringParser,
		BoolParser,
		NilParser,
		AtomParser,
		ListParser)(st)
	return value, err
}
