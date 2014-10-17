package gisp

import (
	"fmt"
	"reflect"

	p "github.com/Dwarfartisan/goparsec"
	px "github.com/Dwarfartisan/goparsec/parsex"
)

type Bracket struct {
	obj  interface{}
	expr []interface{}
}

func (bracket Bracket) Eval(env Env) (interface{}, error) {
	obj, err := Eval(env, bracket.obj)
	if err != nil {
		return nil, err
	}
	val := reflect.ValueOf(obj)
	switch val.Kind() {
	case reflect.Slice, reflect.String, reflect.Array:
		switch len(bracket.expr) {
		case 1:
			return bracket.evalIndex(env, val)
		case 2, 3:
			return bracket.evalSlice(env, val)
		default:
			return nil, fmt.Errorf("Unknow howto index or slice:%v[%v]",
				bracket.obj, bracket.expr)
		}

	case reflect.Map:
		if len(bracket.expr) == 1 {
			key := reflect.ValueOf(bracket.expr[0])
			v := val.MapIndex(key)
			return bracket.inter(v), nil
		}
		return nil, fmt.Errorf("Unknow howto parse map %v[%v]",
			bracket.obj, bracket.expr)
	}
	return nil, fmt.Errorf("Unknow howto parse formal as %v[%v]",
		bracket.obj, bracket.expr)
}

func (bracket Bracket) inter(value reflect.Value) interface{} {
	if value.IsValid() {
		return Value(value.Interface())
	}
	return nil
}

func (bracket Bracket) evalIndex(env Env, val reflect.Value) (interface{}, error) {
	i, err := Eval(env, bracket.expr[0])
	if err != nil {
		return nil, err
	}
	if idx, ok := i.(Int); ok {
		v := val.Index(int(idx))
		return bracket.inter(v), nil
	}
	return nil, fmt.Errorf("Index for slice %v[%v]  is invalid data: %v",
		bracket.obj, bracket.expr, i)
}

func (bracket Bracket) evalSlice(env Env, val reflect.Value) (interface{}, error) {
	items, err := Evals(env, bracket.expr...)
	if err != nil {
		return nil, err
	}
	indexs, err := bracket.computeIndexs(val, items)
	if err != nil {
		return nil, err
	}
	switch len(indexs) {
	case 2:
		v := val.Slice(indexs[0], indexs[1])
		return bracket.inter(v), nil
	case 3:
		v := val.Slice3(indexs[0], indexs[1], indexs[2])
		return bracket.inter(v), nil
	}

	return nil, fmt.Errorf("Index for slice %v[%v]  is invalid",
		bracket.obj, bracket.expr)
}
func (bracket Bracket) computeIndexs(val reflect.Value, input []interface{}) ([]int, error) {
	indexs := make([]int, len(input))
	for idx, item := range input {
		if idx < 2 {
			i, err := bracket.computeIndex(val, item)
			if err != nil {
				return nil, err
			}
			indexs[idx] = i
		}
	}
	return indexs, nil
}

func (bracket Bracket) computeIndex(val reflect.Value, input interface{}) (int, error) {
	l := val.Len()
	if index, ok := input.(Int); ok {
		i := int(index)
		if i < 0 {
			i = l - i
		}
		if i < 0 || i > l-1 {
			return 0, fmt.Errorf("Try to slice %v[%v] but %v out range",
				bracket.obj, bracket.expr, index)
		}
		return i, nil
	}
	return 0, fmt.Errorf("Try to slice %v[%v] but %v is invalid",
		bracket.obj, bracket.expr, input)
}

func (bracket Bracket) SetItemBy(env Env, item interface{}) (interface{}, error) {
	obj, err := Eval(env, bracket.obj)
	if err != nil {
		return nil, err
	}
	val := reflect.ValueOf(obj)
	switch val.Kind() {
	case reflect.Map:
		return bracket.SetMapIndex(val, env, item)
	case reflect.Slice:
		return bracket.SetSliceIndex(val, env, item)
	default:
		return nil, fmt.Errorf("excpet %v[%v]=%v but %v is neither slice nor map.",
			obj, item, bracket.expr, obj)
	}
}

func (bracket Bracket) SetMapIndex(val reflect.Value, env Env, item interface{}) (interface{}, error) {
	if len(bracket.expr) != 1 {
		return nil, fmt.Errorf("excpet %v[%v]=%v but %v has error items(only accept one key).",
			val.Interface(), bracket.expr, item, bracket.expr)
	}
	k, err := Eval(env, bracket.expr[0])
	if err != nil {
		return nil, err
	}
	key := reflect.ValueOf(k)
	value := reflect.ValueOf(item)
	val.SetMapIndex(key, value)
	return val.Interface(), nil
}

func (bracket Bracket) SetSliceIndex(val reflect.Value, env Env, item interface{}) (interface{}, error) {
	if len(bracket.expr) < 1 {
		return nil, fmt.Errorf("excpet %v[%v]=%v but %v has error items(only accept one or two key).",
			val.Interface(), bracket.expr, item, bracket.expr)
	}
	item, err := Eval(env, bracket.expr[0])
	if err != nil {
		return nil, err
	}
	index, err := bracket.computeIndex(val, item)
	if err != nil {
		return nil, err
	}
	value := reflect.ValueOf(item)
	val.Index(index).Set(value)
	return val.Interface(), nil
}

func IntVal(st px.ParsexState) (interface{}, error) {
	x, err := st.Next(px.Always)
	if err != nil {
		return nil, err
	}
	if _, ok := x.(Int); ok {
		return x, nil
	}
	return nil, fmt.Errorf("except a Int value but got %v", x)
}

func BracketParser(st p.ParseState) (interface{}, error) {
	return p.Between(p.Rune('['), p.Rune(']'),
		p.SepBy1(ValueParser, p.Rune(':')),
	)(st)
}

func BracketParserExt(env Env) p.Parser {
	return p.Between(p.Rune('['), p.Rune(']'),
		p.SepBy1(ValueParserExt(env), p.Rune(':')),
	)
}

type BracketExpr struct {
	Expr []interface{}
}

func (be BracketExpr) Task(env Env, args ...interface{}) (Lisp, error) {
	if len(args) != 1 {
		return nil, ParsexSignErrorf("Bracket Expression Args Error: except a arg but %v", args)
	}
	return Bracket{args[0], be.Expr}, nil
}

func BracketExprParserExt(env Env) p.Parser {
	return func(st p.ParseState) (interface{}, error) {
		expr, err := BracketParserExt(env)(st)
		if err != nil {
			return nil, err
		}
		return BracketExpr{expr.([]interface{})}, nil
	}
}
