package gisp

import (
	"fmt"
	"reflect"
	"strings"
)

// List 实现基本的 List 类型
type List []interface{}

func (list List) String() string {
	frags := []string{}
	for _, item := range list {
		frags = append(frags, fmt.Sprintf("%v", item))
	}
	body := strings.Join(frags, " ")
	return fmt.Sprintf("(%s)", body)
}

// Eval 实现 Lisp.Eval 方法
func (list List) Eval(env Env) (interface{}, error) {
	l := len(list)
	if l == 0 {
		return nil, nil
	}
	var lisp interface{}
	switch fun := list[0].(type) {
	case Atom:
		var ok bool
		if lisp, ok = env.Lookup(fun.Name); !ok {
			return nil, fmt.Errorf("any callable named %s not found", fun.Name)
		}
	case List:
		var err error
		lisp, err = fun.Eval(env)
		if err != nil {
			return nil, err
		}
	case Func:
		lisp = fun
	case Expr:
		lisp = fun
	case Dot:
		lisp = fun
	}
	switch item := lisp.(type) {
	case Expr:
		return item(env)(list[1:]...)
	case Task:
		return item.Eval(env)
	case Go:
		return item.Eval(env)
	case Lambda:
		lisp, err := item.Task(env, list[1:]...)
		if err != nil {
			return nil, err
		}
		return lisp.Eval(env)
	case Func:
		lisp, err := item.Task(env, list[1:]...)
		if err != nil {
			return nil, err
		}
		return lisp.Eval(env)
	case Let:
		return item.Eval(env)
	case Dot:
		v, err := item.Eval(env)
		if err != nil {
			return nil, err
		}
		if expr, ok := v.(Expr); ok {
			return expr(env)(list[1:]...)
		}
		value := v.(reflect.Value)
		if value.Kind() == reflect.Func {
			args, err := Evals(env, list[1:]...)
			if err != nil {
				return nil, err
			}
			values := make([]reflect.Value, len(args))
			for idx, arg := range args {
				values[idx] = reflect.ValueOf(arg)
			}
			res, err := InReflects(value.Call(values))
			if err != nil {
				return nil, err
			}
			data, err := Evals(env, res...)
			if err != nil {
				return nil, err
			}
			if len(data) == 1 {
				return data[0], nil
			} else {
				return data, nil
			}
		}

	}
	return nil, fmt.Errorf("%v:%v is't callable", list[0], reflect.TypeOf(list[0]))
}

func (l List) indexn(index Int) interface{} {
	idx, err := l.Anchor(index)
	if err == nil {
		return l[idx.(int)]
	}
	return nil
}

// IndexIs 的 parser 版本，可以用于 Eval 或 gisp code
func (l List) Anchor(index Int) (interface{}, error) {
	idx := l.IndexIs(index)
	if idx == -1 {
		return nil, fmt.Errorf("List Index Anchor Error: %v out range (0, %d)",
			index, len(l))
	}
	return idx, nil
}

// 将负索引正规化， 返回正负索引对应的正规化索引[0, length)。如果索引index不在[-length, length)的范围内，返回－1
func (l List) IndexIs(index Int) int {
	idx := int(index)
	ll := len(l)
	if 0 <= idx || idx < ll {
		return idx
	}
	if -ll <= idx || idx < 0 {
		return ll + idx
	}
	return -1
}

func (l List) Index(index Int) (interface{}, error) {
	idx, err := l.Anchor(index)
	if err != nil {
		return nil, err
	}
	return l[idx.(int)], nil
}

func Zip(x, y List) List {
	xlen := len(x)
	ylen := len(y)
	l := MaxInts(Int(xlen), Int(ylen))
	ret := make(List, l)
	for i := 0; i < int(l); i++ {
		idx := Int(i)
		ret[i] = L(x.indexn(idx), y.indexn(idx))
	}
	return ret
}

func ZipLess(x, y List) List {
	xlen := len(x)
	ylen := len(y)
	l := MinInts(Int(xlen), Int(ylen))
	ret := make(List, l)
	for i := 0; i < int(l); i++ {
		ret[i] = L(x[i], y[i])
	}
	return ret
}

func L(data ...interface{}) List {
	l := make(List, len(data))
	for idx, item := range data {
		l[idx] = item
	}
	return l
}
