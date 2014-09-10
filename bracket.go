package gisp

import (
	"fmt"
	p "github.com/Dwarfartisan/goparsec"
	px "github.com/Dwarfartisan/goparsec/parsex"
	"reflect"
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
		return value.Interface()
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
		i, err := bracket.computeIndex(val, item)
		if err != nil {
			return nil, err
		}
		indexs[idx] = i
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

func BracketParser(st p.ParseState) (interface{}, error) {
	bracket := p.Between(p.Rune('['), p.Rune(']'),
		p.SepBy1(ValueParser, p.Rune(':')),
	)
	t, err := bracket(st)
	if err != nil {
		return nil, err
	}
	tokens := t.([]interface{})
	stx := px.NewStateInMemory(tokens)
	format := px.Choice(
		px.Binds_(px.StringVal, px.Eof),
		px.Binds_(IntVal, px.Eof),
		px.Binds_(IntVal, IntVal, px.Eof),
		px.Binds_(IntVal, IntVal, IntVal, px.Eof),
	)
	_, err = format(stx)
	if err != nil {
		return nil, err
	}

	return tokens, nil
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
