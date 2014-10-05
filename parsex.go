package gisp

import (
	"fmt"
	//	p "github.com/Dwarfartisan/goparsec/parsex"
	px "github.com/Dwarfartisan/goparsec/parsex"
	"reflect"
)

var Parsex Toolkit = Toolkit{
	Meta: map[string]interface{}{
		"name":     "parsex",
		"category": "package",
	},
	Content: map[string]interface{}{
		"state": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, ParsexSignErrorf("Parsex Arg Error:except args has 1 arg.")
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			switch data := param.(type) {
			case string:
				return Q(NewStringState(data)), nil
			case List:
				return Q(px.NewStateInMemory(data)), nil
			default:
				return nil, fmt.Errorf("Parsex Error: Except create a state from a string or List but %v", data)
			}
		},
		"anyone": ParsexBox(px.AnyOne),
		"one": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, ParsexSignErrorf("One Arg Error:except args has 1 arg.")
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			return ParsexBox(px.TheOne(param)), nil
		},
		"str": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, ParsexSignErrorf("One Arg Error:except args has 1 arg.")
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			return ParsexBox(px.String(param.(string))), nil
		},
		"rune": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, ParsexSignErrorf("Rune Arg Error:except args has 1 arg.")
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			return ParsexBox(px.Rune(param.(rune))), nil
		},
		"anyrune":    ParsexBox(px.AnyRune),
		"anyintx":    ParsexBox(px.AnyInt),
		"anyfloatx":  ParsexBox(px.AnyFloat64),
		"anystringx": ParsexBox(px.StringVal),
		"anyint":     ParsexBox(px.AnyInt),
		"anyfloat":   ParsexBox(px.AnyFloat64),
		"aint":       ParsexBox(px.Int),
		"afloat":     ParsexBox(px.Float),
		"astring":    ParsexBox(px.StringVal),
		"string": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, ParsexSignErrorf("string Arg Error:except args has 1 arg.")
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			var str string
			var ok bool
			if str, ok = param.(string); !ok {
				return nil, ParsexSignErrorf("stringx Arg Error:except 1 string arg.")
			}
			return ParsexBox(px.Str(str)), nil
		},
		"stringx": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, ParsexSignErrorf("stringx Arg Error:except args has 1 arg.")
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			var str string
			var ok bool
			if str, ok = param.(string); !ok {
				return nil, ParsexSignErrorf("stringx Arg Error:except 1 string arg.")
			}
			return ParsexBox(px.String(str)), nil
		},
		"digit": ParsexBox(px.Digit),
		"int": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, ParsexSignErrorf("int Arg Error:except args has 1 arg.")
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			var i Int
			var ok bool
			if i, ok = param.(Int); !ok {
				return nil, ParsexSignErrorf("int Arg Error:except 1 string arg.")
			}
			return ParsexBox(func(st px.ParsexState) (interface{}, error) {
				data, err := px.Int(st)
				if err != nil {
					return nil, st.Trap("gisp parsex error:except a int but error: %v", err)
				}
				if Int(data.(int)) != i {
					return nil, st.Trap("gisp parsex error:except a Int but %v", data)
				}
				return data, nil
			}), nil
		},
		"float": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, ParsexSignErrorf("float Arg Error:except args has 1 arg.")
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			var f Float
			var ok bool
			if f, ok = param.(Float); !ok {
				return nil, ParsexSignErrorf("float Arg Error:except 1 string arg.")
			}
			return ParsexBox(func(st px.ParsexState) (interface{}, error) {
				data, err := px.Float(st)
				if err != nil {
					return nil, st.Trap("gisp parsex error:except a float but error: %v", err)
				}
				if Float(data.(float64)) != f {
					return nil, st.Trap("gisp parsex error:except a Float but %v", data)
				}
				return data, nil
			}), nil
		},
		"eof":    ParsexBox(px.Eof),
		"nil":    ParsexBox(px.Nil),
		"atimex": ParsexBox(px.TimeVal),
		"try": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, ParsexSignErrorf("Parsex Parser Try Error: only accept one parsex parser as arg but %v", args)
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			switch parser := param.(type) {
			case Parsexer:
				return ParsexBox(px.Try(parser.Parser)), nil
			default:
				return nil, ParsexSignErrorf(
					"Try Arg Error:except 1 parser arg but %v.",
					reflect.TypeOf(param))
			}

		},
		"either": func(env Env, args ...interface{}) (Lisp, error) {
			ptype := reflect.TypeOf((px.Parser)(nil))
			params, err := GetArgs(env, px.UnionAll(TypeAs(ptype), TypeAs(ptype), px.Eof), args)
			if err != nil {
				return nil, err
			}
			return ParsexBox(px.Either(params[0].(Parsexer).Parser, params[1].(Parsexer).Parser)), nil
		},
		"return": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, ParsexSignErrorf("Parsex Parser Return Error: only accept one parsex parser as arg but %v", args)
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			return ParsexBox(px.Return(param)), nil
		},
		"option": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 2 {
				return nil, ParsexSignErrorf("Parsex Parser Option Error: only accept two parsex parser as arg but %v", args)
			}
			data, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			param, err := Eval(env, args[1])
			if err != nil {
				return nil, err
			}
			switch parser := param.(type) {
			case Parsexer:
				return ParsexBox(px.Option(data, parser.Parser)), nil
			default:
				return nil, ParsexSignErrorf(
					"Many Arg Error:except 1 parser arg but %v.",
					reflect.TypeOf(param))
			}
		},
		"many1": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, ParsexSignErrorf("Parsex Parser Many1 Erroparserr: only accept one parsex parser as arg but %v", args)
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			switch parser := param.(type) {
			case Parsexer:
				return ParsexBox(px.Many1(parser.Parser)), nil
			default:
				return nil, ParsexSignErrorf(
					"Many1 Arg Error:except 1 parser arg but %v.",
					reflect.TypeOf(param))
			}
		},
		"many": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, ParsexSignErrorf("Parsex Parser Many Error: only accept one parsex parser as arg but %v", args)
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			switch parser := param.(type) {
			case Parsexer:
				return ParsexBox(px.Many(parser.Parser)), nil
			default:
				return nil, ParsexSignErrorf(
					"Many Arg Error:except 1 parser arg but %v.",
					reflect.TypeOf(param))
			}
		},
		"failed": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, ParsexSignErrorf("Parsex Parser Failed Error: only accept one string as arg but %v", args)
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			var str string
			var ok bool
			if str, ok = param.(string); !ok {
				return nil, ParsexSignErrorf("Failed Arg Error:except 1 string arg.")
			}
			return ParsexBox(px.Fail(str)), nil
		},
		"oneof": func(env Env, args ...interface{}) (Lisp, error) {
			params, err := Evals(env, args...)
			if err != nil {
				return nil, err
			}
			return ParsexBox(px.OneOf(params...)), nil
		},
		"noneof": func(env Env, args ...interface{}) (Lisp, error) {
			params, err := Evals(env, args...)
			if err != nil {
				return nil, err
			}
			return ParsexBox(px.NoneOf(params)), nil
		},
		"between": func(env Env, args ...interface{}) (Lisp, error) {
			ptype := reflect.TypeOf((px.Parser)(nil))
			params, err := GetArgs(env, px.UnionAll(TypeAs(ptype), TypeAs(ptype), TypeAs(ptype), px.Eof), args)
			if err != nil {
				return nil, err
			}
			return ParsexBox(px.Between(params[0].(Parsexer).Parser, params[1].(Parsexer).Parser, params[2].(Parsexer).Parser)), nil
		},
		"sepby1": func(env Env, args ...interface{}) (Lisp, error) {
			ptype := reflect.TypeOf((px.Parser)(nil))
			params, err := GetArgs(env, px.UnionAll(TypeAs(ptype), TypeAs(ptype), px.Eof), args)
			if err != nil {
				return nil, err
			}
			return ParsexBox(px.SepBy1(params[0].(Parsexer).Parser, params[1].(Parsexer).Parser)), nil
		},
		"sepby": func(env Env, args ...interface{}) (Lisp, error) {
			ptype := reflect.TypeOf((px.Parser)(nil))
			params, err := GetArgs(env, px.UnionAll(TypeAs(ptype), TypeAs(ptype), px.Eof), args)
			if err != nil {
				return nil, err
			}
			return ParsexBox(px.SepBy(params[0].(Parsexer).Parser, params[1].(Parsexer).Parser)), nil
		},
		"manytil": func(env Env, args ...interface{}) (Lisp, error) {
			ptype := reflect.TypeOf((px.Parser)(nil))
			params, err := GetArgs(env, px.UnionAll(TypeAs(ptype), TypeAs(ptype), px.Eof), args)
			if err != nil {
				return nil, err
			}
			return ParsexBox(px.ManyTil(params[0].(Parsexer).Parser, params[1].(Parsexer).Parser)), nil
		},
		"maybe": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, ParsexSignErrorf("Parsex Parser Maybe Error: only accept one parsex parser as arg but %v", args)
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			switch parser := param.(type) {
			case Parsexer:
				return ParsexBox(px.Maybe(parser.Parser)), nil
			default:
				return nil, ParsexSignErrorf(
					"Manybe Arg Error:except 1 parser arg but %v.",
					reflect.TypeOf(param))
			}
		},
		"skip": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, ParsexSignErrorf("Parsex Parser Maybe Error: only accept one parsex parser as arg but %v", args)
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			switch parser := param.(type) {
			case Parsexer:
				return ParsexBox(px.Skip(parser.Parser)), nil
			default:
				return nil, ParsexSignErrorf(
					"Skip Arg Error:except 1 parser arg but %v.",
					reflect.TypeOf(param))
			}
		},
	},
}

func NewStringState(data string) *px.StateInMemory {
	buf := make([]interface{}, len(data))
	for idx, r := range data {
		buf[idx] = r
	}
	return px.NewStateInMemory(buf)
}

type Parsexer struct {
	Parser px.Parser
}

func (parsex Parsexer) Task(env Env, args ...interface{}) (Lisp, error) {
	if len(args) != 1 {
		return nil, ParsexSignErrorf(
			"Parsex Parser Exprission Error: only accept one parsex state as arg but %v",
			args[0])
	}
	param, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	var st px.ParsexState
	var ok bool
	if st, ok = param.(px.ParsexState); !ok {
		return nil, ParsexSignErrorf(
			"Parsex Parser Exprission Error: only accept one parsex state as arg but %v",
			reflect.TypeOf(args[0]))
	}
	return ParsexTask{parsex.Parser, st}, nil
}

func (parsex Parsexer) Eval(env Env) (interface{}, error) {
	return parsex, nil
}

func ParsexBox(parser px.Parser) Lisp {
	return Parsexer{parser}
}

type ParsexTask struct {
	Parser px.Parser
	State  px.ParsexState
}

func (pt ParsexTask) Eval(env Env) (interface{}, error) {
	return pt.Parser(pt.State)
}
