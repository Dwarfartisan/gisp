package gisp

import (
	"fmt"
	px "github.com/Dwarfartisan/goparsec/parsex"
	"reflect"
)

func typeis(x Atom) func(int, interface{}) (interface{}, error) {
	return func(pos int, data interface{}) (interface{}, error) {
		if data == nil {
			if x.Type.Option() {
				return data, nil
			} else {
				return nil, fmt.Errorf("%v's type not match %v", data, x.Type)
			}
		}
		if reflect.DeepEqual(x.Type.Type, ANY) {
			return data, nil
		}
		if reflect.DeepEqual(x.Type.Type, reflect.TypeOf(data)) {
			return data, nil
		} else {
			return data, fmt.Errorf("%v's type not match %v", data, x.Type)
		}
	}
}

// argParser 构造一个 parsex 解析器，判断输入数据是否与给定类型一致，如果判断成功，构造对应的
// Var。
func argParser(atom Atom) px.Parser {
	one := func(st px.ParsexState) (interface{}, error) {
		if data, err := st.Next(typeis(atom)); err == nil {
			slot := VarSlot(atom.Type)
			slot.Set(data)
			return slot, nil
		} else {
			return nil, err
		}
	}
	if atom.Name == "..." {
		return px.Many(one)
	} else {
		return one
	}
}

// argRing 组成参数解析链的的后续逻辑，供 parsex.Binds 调用
func argRing(atom Atom) func(interface{}) px.Parser {
	return func(x interface{}) px.Parser {
		return func(st px.ParsexState) (interface{}, error) {
			ring, err := argParser(atom)(st)
			if err == nil {
				return append(x.([]Var), ring.([]Var)...), nil
			} else {
				return nil, err
			}
		}
	}
}
