package gisp

import (
	px "github.com/Dwarfartisan/goparsec/parsex"
)

var Propositions Toolkit = Toolkit{
	Meta: map[string]interface{}{
		"name":     "propositions",
		"category": "package",
	},
	Content: map[string]interface{}{
		"lambda": BoxExpr(LambdaExpr),
		"let":    BoxExpr(LetExpr),
		"+":      EvalExpr(ParsexExpr(addx)),
		"add":    EvalExpr(ParsexExpr(addx)),
		"-":      EvalExpr(ParsexExpr(subx)),
		"sub":    EvalExpr(ParsexExpr(subx)),
		"*":      EvalExpr(ParsexExpr(mulx)),
		"mul":    EvalExpr(ParsexExpr(mulx)),
		"/":      EvalExpr(ParsexExpr(divx)),
		"div":    EvalExpr(ParsexExpr(divx)),
		"cmp":    EvalExpr(cmpExpr),
		"less":   EvalExpr(lessExpr),
		"<":      EvalExpr(lessExpr),
		"<?":     EvalExpr(lsoExpr),
		"<=":     EvalExpr(leExpr),
		"<=?":    EvalExpr(leoExpr),
		">":      EvalExpr(greatExpr),
		">?":     EvalExpr(gtoExpr),
		">=":     EvalExpr(geExpr),
		">=?":    EvalExpr(geoExpr),
		"==":     EvalExpr(eqsExpr),
		"==?":    EvalExpr(eqsoExpr),
		"!=":     EvalExpr(neqsExpr),
		"!=?":    EvalExpr(neqsoExpr),
	},
}

func ParsexExpr(pxExpr px.Parser) LispExpr {
	return func(env Env, args ...interface{}) (Lisp, error) {
		data, err := Evals(env, args...)
		if err != nil {
			return nil, err
		}
		st := px.NewStateInMemory(data)
		ret, err := pxExpr(st)
		if err != nil {
			return nil, err
		}
		return Q(ret), nil
	}
}

func ExtExpr(extExpr func(env Env) px.Parser) LispExpr {
	return func(env Env, args ...interface{}) (Lisp, error) {
		data, err := Evals(env, args...)
		if err != nil {
			return nil, err
		}
		st := px.NewStateInMemory(data)
		ret, err := extExpr(env)(st)
		if err != nil {
			return nil, err
		}
		return Q(ret), nil
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

func ParsexReverseExpr(pxExpr px.Parser) LispExpr {
	return func(env Env, args ...interface{}) (Lisp, error) {
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
		lisp, err := pxExpr(st)
		if err != nil {
			return nil, err
		}
		return lisp.(Lisp), nil
	}
}

func NotExpr(expr LispExpr) LispExpr {
	return func(env Env, args ...interface{}) (Lisp, error) {
		element, err := expr(env, args...)
		if err != nil {
			return nil, err
		}
		ret, err := element.Eval(env)
		if err != nil {
			return nil, err
		}
		if b, ok := ret.(bool); ok {
			return Q(!b), nil
		} else {
			return nil, ParsexSignErrorf("Unknow howto not %v", b)
		}
	}
}

func OrExpr(x, y px.Parser) LispExpr {
	return func(env Env, args ...interface{}) (Lisp, error) {
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
				return Q(true), nil
			}
			st.SeekTo(0)
			rex, err = y(st)
			if err != nil {
				return nil, err
			}
			return Q(rex), nil
		} else {
			return nil, ParsexSignErrorf("Unknow howto combine %v or %v for %v", x, y, data)
		}
	}
}

func OrExtExpr(x, y func(Env) px.Parser) LispExpr {
	return func(env Env, args ...interface{}) (Lisp, error) {
		return OrExpr(x(env), y(env))(env, args...)
	}
}

func OrExtRExpr(x px.Parser, y func(Env) px.Parser) LispExpr {
	return func(env Env, args ...interface{}) (Lisp, error) {
		return OrExpr(x, y(env))(env, args...)
	}
}

func ExtReverseExpr(expr func(Env) px.Parser) LispExpr {
	return func(env Env, args ...interface{}) (Lisp, error) {
		return ParsexReverseExpr(expr(env))(env, args...)
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
var geoExpr = func(env Env, args ...interface{}) (Lisp, error) {
	st := px.NewStateInMemory(args)
	ret, err := px.Choice(px.Try(NotParsex(less(env))), FalseIfHasNil)(st)
	if err != nil {
		return nil, err
	}
	return Q(ret), nil
}
var eqsExpr = ParsexExpr(equals)
var eqsoExpr = ParsexExpr(equalsOption)
var neqsExpr = NotExpr(eqsExpr)
var neqsoExpr = ParsexExpr(neqsOption)
