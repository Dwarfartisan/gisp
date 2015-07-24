package gisp

import (
	"fmt"
	"reflect"

	p "github.com/Dwarfartisan/goparsec"
	px "github.com/Dwarfartisan/goparsec/parsex"
)

// Parsec 定义了 parsec 包的结构
var Parsec = Toolkit{
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
				return Q(p.MemoryParseState(data)), nil
			default:
				return nil, fmt.Errorf("Parsex Error: Except create a state from a string or List but %v", data)
			}
		},
		"s2str": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, ParsexSignErrorf("Slice to string Arg Error:except args has 1 arg.")
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			var (
				slice []interface{}
				ok    bool
			)
			if slice, ok = param.([]interface{}); !ok {
				return nil, ParsexSignErrorf("s2str Arg Error:except 1 []interface{} arg.")
			}
			return Q(p.ExtractString(slice)), nil
		},
		"str": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, ParsexSignErrorf("Str Arg Error:except args has 1 arg.")
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			var (
				data string
				ok   bool
			)
			if data, ok = param.(string); !ok {
				return nil, ParsexSignErrorf("Str Arg Error:except args has 1 string arg.")
			}
			return ParsecBox(p.String(data)), nil
		},
		"rune": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, ParsexSignErrorf("Rune Arg Error:except args has 1 arg.")
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			var (
				data Rune
				ok   bool
			)
			if data, ok = param.(Rune); !ok {
				return nil, ParsexSignErrorf("One Arg Error:except args has 1 rune arg but %v.", reflect.TypeOf(param))
			}
			return ParsecBox(p.Rune(rune(data))), nil
		},
		"anyone": ParsecBox(p.AnyRune),
		"int":    ParsecBox(p.Int),
		"float":  ParsecBox(p.Float),
		"digit":  ParsecBox(p.Digit),
		"eof":    ParsecBox(p.Eof),
		"try": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, ParsexSignErrorf("Parsec Parser Try Error: only accept one Parsec Parser as arg but %v", args)
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			switch parser := param.(type) {
			case Parsecer:
				return ParsecBox(p.Try(parser.Parser)), nil
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
			return ParsecBox(p.Either(params[0].(Parsecer).Parser, params[1].(Parsecer).Parser)), nil
		},
		"choice": func(env Env, args ...interface{}) (Lisp, error) {
			ptype := reflect.TypeOf((px.Parser)(nil))
			params, err := GetArgs(env, px.ManyTil(TypeAs(ptype), px.Eof), args)
			if err != nil {
				return nil, err
			}
			parsers := make([]p.Parser, len(params))
			for idx, prs := range params {
				if parser, ok := prs.(Parsecer); ok {
					parsers[idx] = parser.Parser
				}
				return nil, ParsexSignErrorf("Choice Args Error:except parsec parsers but %v is %v",
					prs, reflect.TypeOf(prs))
			}
			return ParsecBox(p.Choice(parsers...)), nil
		},
		"return": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, ParsexSignErrorf("Parsec Parser Return Error: only accept one Parsec Parser as arg but %v", args)
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			return ParsecBox(p.Return(param)), nil
		},
		"option": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 2 {
				return nil, ParsexSignErrorf("Parsec Parser Option Error: only accept two Parsec Parser as arg but %v", args)
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
			case Parsecer:
				return ParsecBox(p.Option(data, parser.Parser)), nil
			default:
				return nil, ParsexSignErrorf(
					"Many Arg Error:except 1 parser arg but %v.",
					reflect.TypeOf(param))
			}
		},
		"many1": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, ParsexSignErrorf("Parsec Parser Many1 Erroparserr: only accept one Parsec Parser as arg but %v", args)
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			switch parser := param.(type) {
			case Parsecer:
				return ParsecBox(p.Many1(parser.Parser)), nil
			default:
				return nil, ParsexSignErrorf(
					"Many1 Arg Error:except 1 parser arg but %v.",
					reflect.TypeOf(param))
			}
		},
		"many": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, ParsexSignErrorf("Parsec Parser Many Error: only accept one Parsec Parser as arg but %v", args)
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			switch parser := param.(type) {
			case Parsecer:
				return ParsecBox(p.Many(parser.Parser)), nil
			default:
				return nil, ParsexSignErrorf(
					"Many Arg Error:except 1 parser arg but %v.",
					reflect.TypeOf(param))
			}
		},
		"failed": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, ParsexSignErrorf("Parsec Parser Failed Error: only accept one string as arg but %v", args)
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
			return ParsecBox(p.Fail(str)), nil
		},
		"oneof": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, ParsexSignErrorf("OneOf Arg Error:except args has 1 arg.")
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			var (
				data string
				ok   bool
			)
			if data, ok = param.(string); !ok {
				return nil, ParsexSignErrorf("OneOf Arg Error:except args has 1 string arg.")
			}
			return ParsecBox(p.OneOf(data)), nil
		},
		"noneof": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, ParsexSignErrorf("NoneOf Arg Error:except args has 1 arg.")
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			var (
				data string
				ok   bool
			)
			if data, ok = param.(string); !ok {
				return nil, ParsexSignErrorf("NoneOf Arg Error:except args has 1 string arg.")
			}
			return ParsecBox(p.NoneOf(data)), nil
		},
		"between": func(env Env, args ...interface{}) (Lisp, error) {
			ptype := reflect.TypeOf((*Parsecer)(nil)).Elem()
			params, err := GetArgs(env, px.UnionAll(TypeAs(ptype), TypeAs(ptype), TypeAs(ptype), px.Eof), args)
			if err != nil {
				return nil, err
			}
			return ParsecBox(p.Between(params[0].(Parsecer).Parser, params[1].(Parsecer).Parser, params[2].(Parsecer).Parser)), nil
		},
		"bind": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 2 {
				return nil, ParsexSignErrorf("Bind Args Error:except 2 args.")
			}
			prs, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			var parser Parsecer
			var ok bool
			if parser, ok = prs.(Parsecer); !ok {
				return nil, ParsexSignErrorf("Bind Args Error:except first arg is a parsecer.")
			}
			f, err := Eval(env, args[1])
			if err != nil {
				return nil, err
			}
			switch fun := f.(type) {
			case func(interface{}) p.Parser:
				return ParsecBox(p.Bind(parser.Parser, fun)), nil
			case Functor:
				return ParsecBox(p.Bind(parser.Parser, func(x interface{}) p.Parser {
					tasker, err := fun.Task(env, x)
					if err != nil {
						return func(st p.ParseState) (interface{}, error) {
							return nil, err
						}
					}
					pr, err := tasker.Eval(env)
					if err != nil {
						return func(st p.ParseState) (interface{}, error) {
							return nil, err
						}
					}
					switch parser := pr.(type) {
					case p.Parser:
						return parser
					case Parsecer:
						return parser.Parser
					default:
						return func(st p.ParseState) (interface{}, error) {
							return nil, ParsexSignErrorf("excpet got a parser but %v", pr)
						}
					}
				})), nil
			default:
				return nil, ParsexSignErrorf("excpet got a parser but %v", prs)
			}
		},
		"bind_": func(env Env, args ...interface{}) (Lisp, error) {
			ptype := reflect.TypeOf((*Parsecer)(nil)).Elem()
			params, err := GetArgs(env, px.UnionAll(TypeAs(ptype), TypeAs(ptype), px.Eof), args)
			if err != nil {
				return nil, err
			}
			return ParsecBox(p.Bind_(params[0].(Parsecer).Parser, params[1].(Parsecer).Parser)), nil
		},
		"sepby1": func(env Env, args ...interface{}) (Lisp, error) {
			ptype := reflect.TypeOf((*Parsecer)(nil)).Elem()
			params, err := GetArgs(env, px.UnionAll(TypeAs(ptype), TypeAs(ptype), px.Eof), args)
			if err != nil {
				return nil, err
			}
			return ParsecBox(p.SepBy1(params[0].(Parsecer).Parser, params[1].(Parsecer).Parser)), nil
		},
		"sepby": func(env Env, args ...interface{}) (Lisp, error) {
			ptype := reflect.TypeOf((*Parsecer)(nil)).Elem()
			params, err := GetArgs(env, px.UnionAll(TypeAs(ptype), TypeAs(ptype), px.Eof), args)
			if err != nil {
				return nil, err
			}
			return ParsecBox(p.SepBy(params[0].(Parsecer).Parser, params[1].(Parsecer).Parser)), nil
		},
		"manytil": func(env Env, args ...interface{}) (Lisp, error) {
			ptype := reflect.TypeOf((*Parsecer)(nil)).Elem()
			params, err := GetArgs(env, px.UnionAll(TypeAs(ptype), TypeAs(ptype), px.Eof), args)
			if err != nil {
				return nil, err
			}
			return ParsecBox(p.ManyTil(params[0].(Parsecer).Parser, params[1].(Parsecer).Parser)), nil
		},
		"maybe": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, ParsexSignErrorf("Parsec Parser Maybe Error: only accept one parsec parser as arg but %v", args)
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			switch parser := param.(type) {
			case Parsecer:
				return ParsecBox(p.Maybe(parser.Parser)), nil
			default:
				return nil, ParsexSignErrorf(
					"Manybe Arg Error:except 1 parser arg but %v.",
					reflect.TypeOf(param))
			}
		},
		"skip": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, ParsexSignErrorf("Parsec Parser Skip Error: only accept one parsec parser as arg but %v", args)
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			switch parser := param.(type) {
			case Parsecer:
				return ParsecBox(p.Skip(parser.Parser)), nil
			default:
				return nil, ParsexSignErrorf(
					"Skip Arg Error:except 1 parser arg but %v.",
					reflect.TypeOf(param))
			}
		},
	},
}

// Parsecer 定了一个对 Parsec 解释器的封装
type Parsecer struct {
	Parser p.Parser
}

// Task 实现了其求值逻辑
func (parsec Parsecer) Task(env Env, args ...interface{}) (Lisp, error) {
	if len(args) != 1 {
		return nil, ParsexSignErrorf(
			"Parsec Parser Exprission Error: only accept one parsec state as arg but %v",
			args[0])
	}
	param, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	var st p.ParseState
	var ok bool
	if st, ok = param.(p.ParseState); !ok {
		return nil, ParsexSignErrorf(
			"Parsec Parser Exprission Error: only accept one parsec state as arg but %v",
			reflect.TypeOf(args[0]))
	}
	return ParsecTask{parsec.Parser, st}, nil
}

// Eval 实现了 Parsecer 的求值解析
func (parsec Parsecer) Eval(env Env) (interface{}, error) {
	return parsec, nil
}

// ParsecBox 返回一个封装的 paer
func ParsecBox(parser p.Parser) Lisp {
	return Parsecer{parser}
}

// ParsecTask 是延迟执行 Parsec 逻辑的封装
type ParsecTask struct {
	Parser p.Parser
	State  p.ParseState
}

// Eval 实现了求值
func (pt ParsecTask) Eval(env Env) (interface{}, error) {
	return pt.Parser(pt.State)
}
