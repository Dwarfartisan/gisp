package gisp

import (
	px "github.com/Dwarfartisan/goparsec/parsex"
)

var Propositions Toolkit = Toolkit{
	Meta: map[string]interface{}{
		"name":     "propositions",
		"category": "package",
	},
	Content: map[string]Expr{
		"lambda": LambdaExpr,
		"let":    LetExpr,
		"+":      addExpr,
		"add":    addExpr,
		"-":      subExpr,
		"sub":    subExpr,
		"*":      mulExpr,
		"mul":    mulExpr,
		"/":      divExpr,
		"div":    divExpr,
		"cmp":    cmpExpr,
		"less":   lessExpr,
		"<":      lessExpr,
		"<?":     lsoExpr,
		"<=":     leExpr,
		"<=?":    leoExpr,
		">":      greatExpr,
		">?":     gtoExpr,
		">=":     geExpr,
		">=?":    gtoExpr,
		"==":     eqsExpr,
		"==?":    eqsoExpr,
		"!=":     neqsExpr,
		"!=?":    neqsoExpr,
	},
}

func ParsexExpr(pxExpr px.Parser) Expr {
	return func(env Env) Element {
		return func(args ...interface{}) (interface{}, error) {
			data, err := Evals(env, args...)
			if err != nil {
				return nil, err
			}
			st := px.NewStateInMemory(data)
			return pxExpr(st)
		}
	}
}

func ExtExpr(extExpr func(Env) px.Parser) Expr {
	return func(env Env) Element {
		return func(args ...interface{}) (interface{}, error) {
			data, err := Evals(env, args...)
			if err != nil {
				return nil, err
			}
			st := px.NewStateInMemory(data)
			return extExpr(env)(st)
		}
	}
}

func NotParsex(pxExpr px.Parser) px.Parser {
	return func(st px.ParsexState) (interface{}, error) {
		b, err := pxExpr(st)
		if err != nil {
			return nil, err
		}
		if boolean, ok := b.(bool); ok {
			return !boolean, nil
		} else {
			return nil, ParsexSignErrorf("Unknow howto not %v", b)
		}
	}
}

func ParsexReverseExpr(pxExpr px.Parser) Expr {
	return func(env Env) Element {
		return func(args ...interface{}) (interface{}, error) {
			data, err := Evals(env, args...)
			if err != nil {
				return nil, err
			}
			l := len(data)
			last := l - 1
			datax := make([]interface{}, l)
			for idx, item := range data {
				datax[last-idx] = item
			}
			st := px.NewStateInMemory(data)
			return pxExpr(st)
		}
	}
}

func NotExpr(expr Expr) Expr {
	return func(env Env) Element {
		element := expr(env)
		return func(args ...interface{}) (interface{}, error) {
			ret, err := element(args...)
			if err != nil {
				return ret, err
			}
			if b, ok := ret.(bool); ok {
				return !b, nil
			} else {
				return nil, ParsexSignErrorf("Unknow howto not %v", b)
			}
		}
	}
}

func OrExpr(x, y px.Parser) Expr {
	return func(env Env) Element {
		return func(args ...interface{}) (interface{}, error) {
			data, err := Evals(env, args...)
			if err != nil {
				return nil, err
			}
			st := px.NewStateInMemory(data)
			rex, err := x(st)
			if err != nil {
				return nil, err
			}
			if b, ok := rex.(bool); ok {
				if b {
					return true, nil
				}
				st.SeekTo(0)
				return y(st)
			} else {
				return nil, ParsexSignErrorf("Unknow howto combine %v or %v for %v", x, y, args)
			}
		}
	}
}

func OrExtExpr(x, y func(Env) px.Parser) Expr {
	return func(env Env) Element {
		return OrExpr(x(env), y(env))(env)
	}
}

func OrExtRExpr(x px.Parser, y func(Env) px.Parser) Expr {
	return func(env Env) Element {
		return OrExpr(x, y(env))(env)
	}
}

func ExtReverseExpr(expr func(Env) px.Parser) Expr {
	return func(env Env) Element {
		return ParsexReverseExpr(expr(env))(env)
	}
}

var addExpr = ParsexExpr(addx)
var subExpr = ParsexExpr(subx)
var mulExpr = ParsexExpr(mulx)
var divExpr = ParsexExpr(divx)
var lessExpr = ExtExpr(less)
var lsoExpr = ExtExpr(lessOption)
var leExpr = OrExtRExpr(equals, less)
var leoExpr = OrExtRExpr(equalsOption, lessOption)
var cmpExpr = ParsexExpr(compare)
var greatExpr = ExtReverseExpr(less)
var gtoExpr = ExtReverseExpr(lessOption)
var geExpr = OrExtRExpr(equals, less)
var geoExpr = func(env Env) Element {
	return func(args ...interface{}) (interface{}, error) {
		st := px.NewStateInMemory(args)
		return px.Choice(px.Try(NotParsex(less(env))), FalseIfHasNil)(st)
	}
}
var eqsExpr = ParsexExpr(equals)
var eqsoExpr = ParsexExpr(equalsOption)
var neqsExpr = NotExpr(eqsExpr)
var neqsoExpr = ParsexExpr(neqsOption)
