package gisp

import (
	"fmt"
	px "github.com/Dwarfartisan/goparsec/parsex"
	"reflect"
)

type TypeMatchError struct {
	Value interface{}
	Type  reflect.Type
}

func (err TypeMatchError) Error() string {
	return fmt.Sprintf("%v not match type %v", err.Value, err.Type)
}

type NotIntError struct {
	Value interface{}
}

func (err NotIntError) Error() string {
	return fmt.Sprintf("%v is't a valid Int", err.Value)
}

type NotFloatError struct {
	Value interface{}
}

func (err NotFloatError) Error() string {
	return fmt.Sprintf("%v is't a valid Float", err.Value)
}

type NotNumberError struct {
	Value interface{}
}

func (err NotNumberError) Error() string {
	return fmt.Sprintf("%v is't a valid Number", err.Value)
}

// IntValue 将所有整型处理为 Int ，其它类型不接受
func IntValue(st px.ParsexState) (interface{}, error) {
	v, err := st.Next(px.Always)
	if err != nil {
		return nil, err
	}
	switch val := v.(type) {
	case int:
		return Int(val), nil
	case int8:
		return Int(val), nil
	case int16:
		return Int(val), nil
	case int32:
		return Int(val), nil
	case int64:
		return Int(val), nil
	case Int:
		return val, nil
	default:
		return nil, NotIntError{v}
	}
}

// FloatValue 将所有浮点型处理为 Float ，其它类型不接受
func FloatValue(st px.ParsexState) (interface{}, error) {
	v, err := st.Next(px.Always)
	if err != nil {
		return nil, err
	}
	switch val := v.(type) {
	case float32:
		return Float(val), nil
	case float64:
		return Float(val), nil
	case Float:
		return val, nil
	default:
		return nil, NotFloatError{v}
	}
}

// NumberValue 将所有整型和浮点型处理为 Float ，其它类型不接受
func NumberValue(st px.ParsexState) (interface{}, error) {
	v, err := st.Next(px.Always)
	if err != nil {
		return nil, err
	}
	switch val := v.(type) {
	case int:
		return Float(val), nil
	case int8:
		return Float(val), nil
	case int16:
		return Float(val), nil
	case int32:
		return Float(val), nil
	case int64:
		return Float(val), nil
	case Int:
		return Float(val), nil
	case float32:
		return Float(val), nil
	case float64:
		return Float(val), nil
	case Float:
		return val, nil
	default:
		return nil, NotNumberError{v}
	}
}

// addx 实现一个parsex累加解析器，精度向上适配。我一直觉得应该有一个简单的高效版本，不需要回溯的
// 但是目前还没有找到。
func addx(st px.ParsexState) (interface{}, error) {
	ints, err := px.Try(px.ManyTil(IntValue, px.Eof))(st)
	if err == nil {
		root := Int(0)
		for _, x := range ints.([]interface{}) {
			root += x.(Int)
		}
		return root, nil
	}
	numbers, err := px.ManyTil(NumberValue, px.Eof)(st)
	if err == nil {
		root := Float(0)
		for _, x := range numbers.([]interface{}) {
			root += x.(Float)
		}
		return root, nil
	}

	if nerr, ok := err.(NotNumberError); ok {
		return nil, TypeSignError{Type: FLOATMUST, Value: nerr.Value}
	}
	return nil, err
}

// subx 实现一个左折叠的 parsex 连减解析器，精度向上适配。
func subx(st px.ParsexState) (interface{}, error) {
	data, err := px.Try(px.ManyTil(IntValue, px.Eof))(st)
	if err == nil {
		ints := data.([]interface{})
		root := ints[0].(Int)
		for _, x := range ints[1:] {
			root -= x.(Int)
		}
		return root, nil
	}
	data, err = px.ManyTil(NumberValue, px.Eof)(st)
	if err == nil {
		numbers := data.([]interface{})
		root := numbers[0].(Float)
		for _, x := range numbers[1:] {
			root -= x.(Float)
		}
		return root, nil
	}

	if nerr, ok := err.(NotNumberError); ok {
		return nil, TypeSignError{Type: Type{FLOAT, false}, Value: nerr.Value}
	}
	return nil, err
}

// mulx 实现一个 parsex 累乘解析器，精度向上适配。
func mulx(st px.ParsexState) (interface{}, error) {
	data, err := px.Try(px.ManyTil(IntValue, px.Eof))(st)
	if err == nil {
		ints := data.([]interface{})
		root := ints[0].(Int)
		for _, x := range ints[1:] {
			root *= x.(Int)
		}
		return root, nil
	}
	data, err = px.ManyTil(NumberValue, px.Eof)(st)
	if err == nil {
		numbers := data.([]interface{})
		root := numbers[0].(Float)
		for _, x := range numbers[1:] {
			root *= x.(Float)
		}
		return root, nil
	}

	if nerr, ok := err.(NotNumberError); ok {
		return nil, TypeSignError{Type: Type{FLOAT, false}, Value: nerr.Value}
	}
	return nil, err
}

// divx 实现一个左折叠的 parsex 连除解析器，精度向上适配。
func divx(st px.ParsexState) (interface{}, error) {
	data, err := px.Try(px.ManyTil(IntValue, px.Eof))(st)
	if err == nil {
		ints := data.([]interface{})
		root := ints[0].(Int)
		for _, x := range ints[1:] {
			root /= x.(Int)
		}
		return root, nil
	}
	data, err = px.ManyTil(NumberValue, px.Eof)(st)
	if err == nil {
		numbers := data.([]interface{})
		root := numbers[0].(Float)
		for _, x := range numbers[1:] {
			root /= x.(Float)
		}
		return root, nil
	}

	if nerr, ok := err.(NotNumberError); ok {
		return nil, TypeSignError{Type: Type{FLOAT, false}, Value: nerr.Value}
	}
	return nil, err
}
