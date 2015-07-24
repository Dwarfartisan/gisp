package gisp

import (
	"fmt"
	"io"
	"reflect"
	tm "time"

	px "github.com/Dwarfartisan/goparsec/parsex"
)

// FalseIfHasNil 实现一个 is nil 判断
func FalseIfHasNil(st px.ParsexState) (interface{}, error) {
	for {
		val, err := px.AnyOne(st)
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("False If Nil Error: Found EOF.")
			}
			return nil, err
		}
		if val == nil {
			return false, err
		}
	}
}

// LessThanNil 实现三值 Less
func LessThanNil(x interface{}) px.Parser {
	return func(st px.ParsexState) (interface{}, error) {
		val, _ := px.AnyOne(st)
		if x == nil || val == nil {
			return false, nil
		}
		return nil, ParsexSignErrorf("except nil value but: %v", val)
	}
}

// ListValue 检查 state 中下一个值是否是列表
func ListValue(st px.ParsexState) (interface{}, error) {
	val, err := px.AnyOne(st)
	if err == nil {
		if _, ok := val.(List); ok {
			return val, nil
		}
		return nil, ParsexSignErrorf("except a list value but %v ", val)
	}
	return nil, ParsexSignErrorf("except a list value but error: %v", err)
}

// LessThanList 从最近的环境中找到 < 的实现并调用其进行比较，这样用户可以自己实现特化的比较
func LessThanList(env Env) func(x interface{}) px.Parser {
	lessp, ok := env.Lookup("<")
	return func(x interface{}) px.Parser {
		return func(st px.ParsexState) (interface{}, error) {
			if !ok {
				return nil, fmt.Errorf("Less Than List Error: opreator < not found")
			}
			y, err := ListValue(st)
			if err != nil {
				return nil, err
			}
			for _, item := range ZipLess(x.(List), y.(List)) {
				b, err := Eval(env, L(lessp, item.(List)[0], item.(List)[1]))
				if err != nil {
					return nil, err
				}
				if b.(bool) {
					return true, nil
				}
			}
			return len(x.(List)) < len(y.(List)), nil
		}
	}
}

// LessThanListOption 允许返回 nil
func LessThanListOption(env Env) func(x interface{}) px.Parser {
	lessp, ok := env.Lookup("<?")
	return func(x interface{}) px.Parser {
		return func(st px.ParsexState) (interface{}, error) {
			if !ok {
				return nil, fmt.Errorf("Less Than List Option Error: <? opreator not found")
			}
			y, err := ListValue(st)
			if err != nil {
				return nil, err
			}
			for _, item := range ZipLess(x.(List), y.(List)) {
				b, err := Eval(env, L(lessp, item.(List)[0], item.(List)[1]))
				if err != nil {
					return nil, err
				}
				if b.(bool) {
					return true, nil
				}
			}
			return len(x.(List)) < len(y.(List)), nil
		}
	}
}

// TimeValue 判断 state 中下一个元素是否为 time.Time
func TimeValue(st px.ParsexState) (interface{}, error) {
	val, err := px.AnyOne(st)
	if err == nil {
		if _, ok := val.(tm.Time); ok {
			return val, nil
		}
		return nil, ParsexSignErrorf("except a time value but: %v", err)
	}
	return nil, ParsexSignErrorf("except a time value but error: %v", err)
}

// LessThanTime 对 Time 值进行比较
func LessThanTime(x interface{}) px.Parser {
	return func(st px.ParsexState) (interface{}, error) {
		y, err := TimeValue(st)
		if err == nil {
			return x.(tm.Time).Before(y.(tm.Time)), nil
		}
		return nil, err
	}
}

// StringValue 判断 state 中下一个值是否为 String
func StringValue(st px.ParsexState) (interface{}, error) {
	val, err := px.StringVal(st)
	if err == nil {
		return val, nil
	}
	return nil, ParsexSignErrorf("except a string value but error: %v", err)
}

// LessThanInt 实现整数的比较
func LessThanInt(x interface{}) px.Parser {
	return func(st px.ParsexState) (interface{}, error) {
		y, err := IntValue(st)
		if err == nil {
			return x.(Int) < y.(Int), nil
		}
		return nil, err
	}
}

// LessThanFloat 实现浮点数的比较
func LessThanFloat(x interface{}) px.Parser {
	return func(st px.ParsexState) (interface{}, error) {
		y, err := FloatValue(st)
		if err == nil {
			switch val := x.(type) {
			case Float:
				return val < y.(Float), nil
			case Int:
				return Float(val) < y.(Float), nil
			default:
				return nil, ParsexSignErrorf("unknown howto compoare %v < %v", x, y)
			}
		}
		return nil, err
	}
}

// LessThanNumber 实现数值的比较
func LessThanNumber(x interface{}) px.Parser {
	return func(st px.ParsexState) (interface{}, error) {
		pos := st.Pos()
		cmp, err := LessThanInt(x)(st)
		if err == nil {
			return cmp, nil
		}
		st.SeekTo(pos)
		return LessThanFloat(x)(st)
	}
}

// LessThanString 实现字符串的比较
func LessThanString(x interface{}) px.Parser {
	return func(st px.ParsexState) (interface{}, error) {
		y, err := px.StringVal(st)
		if err == nil {
			return x.(string) < y.(string), nil
		}
		return nil, ParsexSignErrorf("Except less compare string %v and %v but error: %v",
			x, y, err)
	}
}

func lessListIn(env Env, x, y List) (interface{}, error) {
	lessp, ok := env.Lookup("<")
	if !ok {
		return nil, fmt.Errorf("Less Than List Error: < opreator not found")
	}
	for _, item := range ZipLess(x, y) {
		b, err := Eval(env, L(lessp, item.(List)[0], item.(List)[1]))
		if err != nil {
			return nil, err
		}
		if b.(bool) {
			return true, nil
		}
	}
	return len(x) < len(y), nil
}

func lessListOptIn(env Env, x, y List) (interface{}, error) {
	lessp, ok := env.Lookup("<?")
	if !ok {
		return nil, fmt.Errorf("Less Than Option List Error: opreator <? not found")
	}
	for _, item := range ZipLess(x, y) {
		b, err := Eval(env, L(lessp, item.(List)[0], item.(List)[1]))
		if err != nil {
			return nil, err
		}
		if b.(bool) {
			return true, nil
		}
	}
	return len(x) < len(y), nil
}

func less(env Env) px.Parser {
	return func(st px.ParsexState) (interface{}, error) {
		l, err := px.Bind(px.Choice(
			px.Try(px.Bind(IntValue, LessThanNumber)),
			px.Try(px.Bind(NumberValue, LessThanFloat)),
			px.Try(px.Bind(px.StringVal, LessThanString)),
			px.Try(px.Bind(TimeValue, LessThanTime)),
			px.Bind(ListValue, LessThanList(env)),
		), func(l interface{}) px.Parser {
			return func(st px.ParsexState) (interface{}, error) {
				_, err := px.Eof(st)
				if err != nil {
					return nil, ParsexSignErrorf("less args sign error: except eof")
				}
				return l, nil
			}
		})(st)
		if err == nil {
			return l, nil
		}
		return nil, ParsexSignErrorf("Except two lessable values compare but error %v", err)
	}
}

// return false, true or nil
func lessOption(env Env) px.Parser {
	return func(st px.ParsexState) (interface{}, error) {
		l, err := px.Bind(px.Choice(
			px.Try(px.Bind(IntValue, LessThanNumber)),
			px.Try(px.Bind(NumberValue, LessThanFloat)),
			px.Try(px.Bind(px.StringVal, LessThanString)),
			px.Try(px.Bind(TimeValue, LessThanTime)),
			px.Try(px.Bind(ListValue, LessThanListOption(env))),
			px.Bind(px.AnyOne, LessThanNil),
		), func(l interface{}) px.Parser {
			return func(st px.ParsexState) (interface{}, error) {
				_, err := px.Eof(st)
				if err != nil {
					return nil, ParsexSignErrorf("less args sign error: except eof")
				}
				return l, nil
			}
		})(st)
		if err == nil {
			return l, nil
		}
		return nil, ParsexSignErrorf("Except two lessable values or nil compare but error: %v", err)
	}
}

func cmpInt(x, y Int) Int {
	if x < y {
		return Int(1)
	}
	if y < x {
		return Int(-1)
	}
	if x == y {
		return Int(0)
	}
	return Int(0)
}

func cmpFloat(x, y Float) Int {
	if x < y {
		return Int(1)
	}
	if y < x {
		return Int(-1)
	}
	if x == y {
		return Int(0)
	}
	return Int(0)
}

func cmpString(x, y string) Int {
	if x < y {
		return Int(1)
	}
	if y < x {
		return Int(-1)
	}
	if x == y {
		return Int(0)
	}
	return Int(0)
}

func cmpTime(x, y tm.Time) Int {
	if x.Before(y) {
		return Int(1)
	}
	if x.After(y) {
		return Int(-1)
	}
	return Int(0)
}

func cmpListIn(env Env, x, y List) (interface{}, error) {
	ret, err := lessListIn(env, x, y)
	if err != nil {
		return nil, err
	}
	if ret.(bool) {
		return -1, nil
	}
	ret, err = lessListIn(env, y, x)
	if err != nil {
		return nil, err
	}
	if ret.(bool) {
		return 1, nil
	}
	if reflect.DeepEqual(x, y) {
		return 0, nil
	}
	return nil, fmt.Errorf("Compare Error: Unknown howto copmare %v and %v", x, y)
}

// CmpInt 实现两个整数的三向比较
func CmpInt(x interface{}) px.Parser {
	return func(st px.ParsexState) (interface{}, error) {
		y, err := IntValue(st)
		if err == nil {
			return cmpInt(x.(Int), y.(Int)), nil
		}
		return nil, err
	}
}

// CmpFloat 实现两个浮点数的三向比较
func CmpFloat(x interface{}) px.Parser {
	return func(st px.ParsexState) (interface{}, error) {
		y, err := FloatValue(st)
		if err == nil {
			switch val := x.(type) {
			case Float:
				return cmpFloat(val, y.(Float)), nil
			case Int:
				return cmpFloat(Float(val), y.(Float)), nil
			default:
				return nil, ParsexSignErrorf("unknown howto compoare %v < %v", x, y)
			}
		}
		return nil, err
	}
}

// CmpNumber 实现两个数值的三向比较
func CmpNumber(x interface{}) px.Parser {
	return func(st px.ParsexState) (interface{}, error) {
		pos := st.Pos()
		cmp, err := CmpInt(x)(st)
		if err == nil {
			return cmp, nil
		}
		st.SeekTo(pos)
		return CmpFloat(x)(st)
	}
}

// CmpString 实现两个字符串的三向比较
func CmpString(x interface{}) px.Parser {
	return func(st px.ParsexState) (interface{}, error) {
		y, err := px.StringVal(st)
		if err == nil {
			return cmpString(x.(string), y.(string)), nil
		}
		return nil, ParsexSignErrorf("Except less compare string %v and %v but error: %v",
			x, y, err)
	}
}

// CmpTime 实现两个Time的三向比较
func CmpTime(x interface{}) px.Parser {
	return func(st px.ParsexState) (interface{}, error) {
		y, err := TimeValue(st)
		if err == nil {
			return cmpTime(x.(tm.Time), y.(tm.Time)), nil
		}
		return nil, ParsexSignErrorf("Except less compare string %v and %v but error: %v",
			x, y, err)
	}
}

func compare(st px.ParsexState) (interface{}, error) {
	l, err := px.Bind(px.Choice(
		px.Bind(IntValue, LessThanNumber),
		px.Bind(NumberValue, LessThanFloat),
		px.Bind(px.StringVal, LessThanString),
		px.Bind(TimeValue, LessThanTime),
	), func(l interface{}) px.Parser {
		return func(st px.ParsexState) (interface{}, error) {
			_, err := px.Eof(st)
			if err != nil {
				return nil, ParsexSignErrorf("less args sign error: except eof")
			}
			return l, nil
		}
	})(st)
	if err == nil {
		return l, nil
	}
	return nil, ParsexSignErrorf("Except two lessable values compare but error %v", err)
}

func equals(st px.ParsexState) (interface{}, error) {
	return px.Bind(px.AnyOne, eqs)(st)
}
func eqs(x interface{}) px.Parser {
	return func(st px.ParsexState) (interface{}, error) {
		y, err := st.Next(px.Always)
		if err != nil {
			if reflect.DeepEqual(err, io.EOF) {
				return true, nil
			}
			return nil, err
		}
		if reflect.DeepEqual(x, y) {
			return eqs(x)(st)
		}
		return false, nil
	}
}

func equalsOption(st px.ParsexState) (interface{}, error) {
	return px.Bind(px.AnyOne, eqsOption)(st)
}

func eqsOption(x interface{}) px.Parser {
	return func(st px.ParsexState) (interface{}, error) {
		y, err := st.Next(px.Always)
		if err != nil {
			if reflect.DeepEqual(err, io.EOF) {
				return true, nil
			}
			return nil, err
		}
		if x == nil || y == nil {
			return false, nil
		}
		if reflect.DeepEqual(x, y) {
			return eqsOption(x)(st)
		}
		return false, nil
	}
}

func notEquals(st px.ParsexState) (interface{}, error) {
	return px.Bind(px.AnyOne, neqs)(st)
}

func neqs(x interface{}) px.Parser {
	return func(st px.ParsexState) (interface{}, error) {
		y, err := st.Next(px.Always)
		if err != nil {
			if reflect.DeepEqual(err, io.EOF) {
				return false, nil
			}
			return nil, err
		}
		if x == nil || y == nil {
			return false, nil
		}
		if !reflect.DeepEqual(x, y) {
			return neqs(x)(st)
		}
		return false, nil
	}
}

// not equals function, NotEqual or !=, if anyone is nil, return false
func neqsOption(st px.ParsexState) (interface{}, error) {
	x, err := st.Next(px.Always)
	if err != nil {
		return nil, err
	}
	if x == nil {
		return false, nil
	}
	for {
		y, err := st.Next(px.Always)
		if err != nil {
			if reflect.DeepEqual(err, io.EOF) {
				return false, nil
			}
			return nil, err
		}
		if y == nil {
			return false, nil
		}
		if !reflect.DeepEqual(x, y) {
			return true, nil
		}
	}
}

// String2Values 将两个 StringValue 串为 List
var String2Values = px.Bind(StringValue, func(x interface{}) px.Parser {
	return func(st px.ParsexState) (interface{}, error) {
		y, err := StringValue(st)
		if err != nil {
			return nil, err
		}
		return []interface{}{x, y}, nil
	}
})

//TimeValue 将两个 Time 值串为 List
var Time2Values = px.Bind(TimeValue, func(x interface{}) px.Parser {
	return func(st px.ParsexState) (interface{}, error) {
		y, err := TimeValue(st)
		if err != nil {
			return nil, err
		}
		return []interface{}{x, y}, nil
	}
})

//ListValue 将两个 Time 值串为 List
var List2Values = px.Bind(ListValue, func(x interface{}) px.Parser {
	return func(st px.ParsexState) (interface{}, error) {
		y, err := ListValue(st)
		if err != nil {
			return nil, err
		}
		return []interface{}{x, y}, nil
	}
})
