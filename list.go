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
	case Functor:
		lisp = fun
	case TaskExpr:
		lisp = fun
	case LispExpr:
		lisp = fun
	case Dot:
		lisp = fun
	case reflect.Value:
		if fun.Kind() == reflect.Func {
			lisp = fun
		}
	}
	switch item := lisp.(type) {
	case TaskExpr:
		task, err := item(env, list[1:]...)
		if err != nil {
			return nil, err
		}
		return task(env)
	case LispExpr:
		lisp, err := item(env, list[1:]...)
		if err != nil {
			return nil, err
		}
		return lisp.Eval(env)
	case Task:
		return item.Eval(env)
	case Functor:
		task, err := item.Task(env, list[1:]...)
		if err != nil {
			return nil, err
		}
		return task.Eval(env)
	case Let:
		return item.Eval(env)
	case Dot:
		v, err := item.Eval(env)
		if err != nil {
			return nil, err
		}
		if expr, ok := v.(TaskExpr); ok {
			tasker, err := expr(env, list[1:]...)
			if err != nil {
				return nil, err
			}
			return tasker(env)
		}
		//if expr, ok := v.(LispExpr); ok {
		if expr, ok := v.(func(Env, ...interface{}) (Lisp, error)); ok {
			lisp, err := expr(env, list[1:]...)
			if err != nil {
				return nil, err
			}
			return lisp.Eval(env)
		}
		if functor, ok := v.(Functor); ok {
			tasker, err := functor.Task(env, list[1:]...)
			if err != nil {
				return nil, err
			}
			return tasker.Eval(env)
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
			var revs []reflect.Value
			if value.Type().IsVariadic() {
				revs = value.CallSlice(values)
			} else {
				revs = value.Call(values)
			}
			res, err := InReflects(revs)
			if err != nil {
				return nil, err
			}
			if len(res) == 1 {
				return res[0], nil
			}
			return res, nil
		}
	case reflect.Value:
		args, err := Evals(env, list[1:]...)
		if err != nil {
			return nil, err
		}
		values := make([]reflect.Value, len(args))
		for idx, arg := range args {
			values[idx] = reflect.ValueOf(arg)
		}
		var revs []reflect.Value
		if item.Type().IsVariadic() {
			revs = item.CallSlice(values)
		} else {
			revs = item.Call(values)
		}
		res, err := InReflects(revs)
		if err != nil {
			return nil, err
		}
		data, err := Evals(env, res...)
		if err != nil {
			return nil, err
		}
		if len(data) == 1 {
			return data[0], nil
		}
		return data, nil
	}
	return nil, fmt.Errorf("List %v Eval Error: %v(%v):%v is't callable",
		list, list[0], lisp, reflect.TypeOf(lisp))
}

func (list List) indexn(index Int) interface{} {
	idx, err := list.Anchor(index)
	if err == nil {
		return list[idx.(int)]
	}
	return nil
}

// Anchor 是 IndexIs 的 parser 版本，可以用于 Eval 或 gisp code
func (list List) Anchor(index Int) (interface{}, error) {
	idx := list.IndexIs(index)
	if idx == -1 {
		return nil, fmt.Errorf("List Index Anchor Error: %v out range (0, %d)",
			index, len(list))
	}
	return idx, nil
}

// IndexIs 将负索引正规化， 返回正负索引对应的正规化索引[0, length)。如果索引index不在[-length, length)的范围内，返回－1
func (list List) IndexIs(index Int) int {
	idx := int(index)
	ll := len(list)
	if 0 <= idx || idx < ll {
		return idx
	}
	if -ll <= idx || idx < 0 {
		return ll + idx
	}
	return -1
}

// Index 实现基本的索引操作
func (list List) Index(index Int) (interface{}, error) {
	idx, err := list.Anchor(index)
	if err != nil {
		return nil, err
	}
	return list[idx.(int)], nil
}

// Zip 合并两个 List
func Zip(x, y List) List {
	xlen := len(x)
	ylen := len(y)
	max := MaxInts(Int(xlen), Int(ylen))
	min := MinInts(Int(xlen), Int(ylen))
	ret := ZipLess(x, y)
	nils := make(List, max-min)
	for idx := range nils {
		nils[idx] = nil
	}
	if xlen > int(min) {
		ret = append(ret, ZipLess(x[min:], nils)...)
	}
	if ylen > int(min) {
		ret = append(ret, ZipLess(nils, y[min:])...)
	}
	return ret
}

// ZipLess 在 Zip 时放弃超出的部分，类似 inner join 的效果
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

// L 构造一个 List
func L(data ...interface{}) List {
	l := make(List, len(data))
	for idx, item := range data {
		l[idx] = item
	}
	return l
}
