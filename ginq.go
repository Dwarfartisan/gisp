package gisp

import (
	"fmt"
	"reflect"
)

/*
Ginq methods:
 - select
 - where
 - all
 - any
 - max
 - min
 - count
 - sum
 - average
 - last
 - first
 - groupby
 - order
 - distinct
 - column
 - take
 － reverse
 - join
*/

// 这个 ginq 构造器效率是比较低的，每一个查询子句都会重新构造一个数据集，提升查询效率有赖
// 调用者调整查询结构。
// 未来应该将其内部逻辑改为构造一个查询语法树，尽可能的减少中间数据集的构造。
type Ginq struct {
	Meta    map[string]interface{}
	queries []interface{}
}

func NewGinq(queries ...interface{}) *Ginq {
	ginq := &Ginq{
		Meta: map[string]interface{}{
			"local": map[string]Var{},
			"ginq": map[string]interface{}{
				"select": LispExpr(func(env Env, args ...interface{}) (Lisp, error) {
					if len(args) != 1 {
						return nil, ParsexSignErrorf("ginq select args error: excpet one expression but: %v", args)
					}
					param, err := Eval(env, args[0])
					if err != nil {
						return nil, err
					}
					return Q(NewGinSelect(param)), nil
				}),
				"where": LispExpr(func(env Env, args ...interface{}) (Lisp, error) {
					if len(args) != 1 {
						return nil, ParsexSignErrorf("ginq where args error: excpet one bool expression but: %v", args)
					}
					param, err := Eval(env, args[0])
					if err != nil {
						return nil, err
					}
					return Q(NewGinWere(param)), nil
				}),
				"groupby": LispExpr(func(env Env, args ...interface{}) (Lisp, error) {
					params, err := Evals(env, args...)
					if err != nil {
						return nil, err
					}
					if len(params) != 2 {
						return nil, ParsexSignErrorf("ginq where args error: excpet one bool expression but: %v", params)
					}
					return Q(NewGinGroup(params[0], params[1])), nil
				}),
				"fs": LispExpr(func(env Env, args ...interface{}) (Lisp, error) {
					params, err := Evals(env, args...)
					if err != nil {
						return nil, err
					}
					return Q(NewGinFields(params...)), nil
				}),
				"sum": TaskExpr(func(env Env, args ...interface{}) (Tasker, error) {
					params, err := Evals(env, args...)
					if err != nil {
						return nil, err
					}
					if len(params) != 1 {
						return nil, ParsexSignErrorf("ginq sum args error: excpet one bool expression but: %v", params)
					}
					var l List
					var ok bool
					if l, ok = params[0].(List); !ok {
						return nil, ParsexSignErrorf("ginq sum args error: except a data List but: %v", params)
					}
					return func(env Env) (interface{}, error) {
						if len(l) == 0 {
							return nil, nil
						}
						if len(l) == 1 {
							return l[0], nil
						}
						add, _ := env.Lookup("+")
						root := l[0]
						for _, item := range l[1:] {
							call := L(add, root, item)
							data, err := Eval(env, call)
							if err != nil {
								return nil, err
							}
							root = data
						}
						return root, nil
					}, nil
				}),
				"max": TaskExpr(func(env Env, args ...interface{}) (Tasker, error) {
					params, err := Evals(env, args...)
					if err != nil {
						return nil, err
					}
					if len(params) != 1 {
						return nil, ParsexSignErrorf("ginq max args error: excpet one bool expression but: %v", params)
					}
					var l List
					var ok bool
					if l, ok = params[0].(List); !ok {
						return nil, ParsexSignErrorf("ginq sum max error: except a data List but: %v", params)
					}
					return func(env Env) (interface{}, error) {
						if len(l) == 0 {
							return nil, nil
						}
						if len(l) == 1 {
							return l[0], nil
						}
						lt, _ := env.Lookup("<")
						root := l[0]
						for _, item := range l[1:] {
							call := L(lt, root, item)
							data, err := Eval(env, call)
							if err != nil {
								return nil, err
							}
							if b, ok := data.(bool); ok {
								if b {
									root = item
								}
							} else {
								return nil, fmt.Errorf("ginq max error: except compare %v and %v got a bool but: %v",
									root, item, data)
							}
						}
						return root, nil
					}, nil
				}),
				"min": TaskExpr(func(env Env, args ...interface{}) (Tasker, error) {
					params, err := Evals(env, args...)
					if err != nil {
						return nil, err
					}
					if len(params) != 1 {
						return nil, ParsexSignErrorf("ginq sum args error: excpet one bool expression but: %v", params)
					}
					var l List
					var ok bool
					if l, ok = params[0].(List); !ok {
						return nil, ParsexSignErrorf("ginq sum args error: except a data List but: %v", params)
					}
					return func(env Env) (interface{}, error) {
						if len(l) == 0 {
							return nil, nil
						}
						if len(l) == 1 {
							return l[0], nil
						}
						lt, _ := env.Lookup("<")
						root := l[0]
						for _, item := range l[1:] {
							call := L(lt, root, item)
							data, err := Eval(env, call)
							if err != nil {
								return nil, err
							}
							if b, ok := data.(bool); ok {
								if b {
									root = item
								}
							} else {
								return nil, fmt.Errorf("ginq min error: except compare %v and %v got a bool but: %v",
									root, item, data)
							}
						}
						return root, nil
					}, nil
				}),
				"avg": TaskExpr(func(env Env, args ...interface{}) (Tasker, error) {
					params, err := Evals(env, args...)
					if err != nil {
						return nil, err
					}
					if len(params) != 1 {
						return nil, ParsexSignErrorf("ginq avg args error: excpet one bool expression but: %v", params)
					}
					var l List
					var ok bool
					if l, ok = params[0].(List); !ok {
						return nil, ParsexSignErrorf("ginq avg args error: except a data List but: %v", params)
					}
					return func(env Env) (interface{}, error) {
						if len(l) == 0 {
							return nil, nil
						}
						if len(l) == 1 {
							return l[0], nil
						}
						add, _ := env.Lookup("+")
						root := l[0]
						for _, item := range l[1:] {
							call := L(add, root, item)
							data, err := Eval(env, call)
							if err != nil {
								return nil, err
							}
							root = data
						}
						div, _ := env.Lookup("/")
						call := L(div, root, len(l))
						rev, err := Eval(env, call)
						if err != nil {
							return nil, err
						}
						return rev, nil
					}, nil
				}),
			},
		},
		queries: queries,
	}
	return ginq
}

func (ginq Ginq) Task(env Env, args ...interface{}) (Lisp, error) {
	if len(args) != 1 {
		return nil, ParsexSignErrorf("ginq avg args error: excpet one expression but: %v", args)
	}
	data, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	var l List
	var ok bool
	if l, ok = data.(List); !ok {
		return nil, ParsexSignErrorf("ginq run error: excpet arg eval got a list but %v", data)
	}
	meta := map[string]interface{}{}
	for k, v := range ginq.Meta {
		meta[k] = v
	}
	return GinQ{meta, ginq.queries, l}, nil
}

type GinQ struct {
	Meta    map[string]interface{}
	queries []interface{}
	data    List
}

func (ginq GinQ) Eval(env Env) (interface{}, error) {
	ginq.Meta["global"] = env
	var rel interface{} = ginq.data
	var err error
	for _, query := range ginq.queries {
		call := L(query, rel)
		rel, err = Eval(ginq, call)
		if err != nil {
			return nil, err
		}
	}
	return rel, nil
}

// Defvar 实现 Env.Defvar
func (ginq GinQ) Defvar(name string, slot Var) error {
	if _, ok := ginq.Local(name); ok {
		return fmt.Errorf("local name %s is exists", name)
	}
	local := ginq.Meta["local"].(map[string]Var)
	local[name] = slot
	return nil
}

// Defun 实现 Env.Defun
func (ginq GinQ) Defun(name string, functor Functor) error {
	if s, ok := ginq.Local(name); ok {
		switch slot := s.(type) {
		case Func:
			slot.Overload(functor)
		case Var:
			return fmt.Errorf("%s defined as a var", name)
		default:
			return fmt.Errorf("exists name %s isn't Expr", name)
		}
	}
	local := ginq.Meta["local"].(map[string]interface{})
	local[name] = NewFunction(name, ginq, functor)
	return nil
}

// Setvar 实现 Env.Setvar
func (ginq GinQ) Setvar(name string, value interface{}) error {
	if _, ok := ginq.Local(name); ok {
		local := ginq.Meta["local"].(map[string]Var)
		local[name].Set(value)
		return nil
	}
	global := ginq.Meta["global"].(Env)
	return global.Setvar(name, value)
}

// Local 实现 Env.Local
func (ginq GinQ) Local(name string) (interface{}, bool) {
	ginfun := ginq.Meta["ginq"].(map[string]interface{})
	if gf, ok := ginfun[name]; ok {
		return gf, true
	}
	local := ginq.Meta["local"].(map[string]Var)
	if slot, ok := local[name]; ok {
		return slot.Get(), true
	}
	return nil, false
}

// Lookup 实现 Env.Lookup
func (ginq GinQ) Lookup(name string) (interface{}, bool) {
	if value, ok := ginq.Local(name); ok {
		return value, true
	}
	return ginq.Global(name)

}

// Global 实现 Env.Global
func (ginq GinQ) Global(name string) (interface{}, bool) {
	global := ginq.Meta["global"].(Env)
	return global.Lookup(name)
}

type GinGroup struct {
	group interface{}
	by    interface{}
}

func NewGinGroup(by interface{}, group interface{}) GinGroup {
	return GinGroup{group, by}
}

func (group GinGroup) Task(env Env, args ...interface{}) (Lisp, error) {
	if len(args) != 1 {
		return nil, ParsexSignErrorf("ginq group by exec error: except group from a list but %v", args)
	}
	var l List
	var ok bool
	if l, ok = args[0].(List); !ok {
		return nil, ParsexSignErrorf("ginq group by exec error: except group from a list but %v", args[0])
	}
	pool := []List{}
	for _, data := range l {
		call := L(group.by, data)
		grp, err := Eval(env, call)
		if err != nil {
			return nil, err
		}
		flag := false
		// group is []List
		for _, gr := range pool {
			if reflect.DeepEqual(gr[0], grp) {
				flag = true
				grpr := gr[1].(List)
				gr[1] = append(grpr, grp)
				break
			}
		}
		if !flag {
			pool = append(pool, L(grp, L(data)))
		}
	}
	rel := make(List, len(pool))
	for idx, g := range pool {
		row := L(g[0])
		call := L(group.group, g[1])
		data, err := Eval(env, call)
		if err != nil {
			return nil, err
		}
		row = append(row, data)
		rel[idx] = row
	}
	return Q(rel), nil
}

type GinSelect struct {
	fun interface{}
}

func NewGinSelect(fun interface{}) GinSelect {
	return GinSelect{fun}
}

func (sel GinSelect) Task(env Env, args ...interface{}) (Lisp, error) {
	if len(args) != 1 {
		return nil, ParsexSignErrorf("ginq select args error: except select from a list but %v", args)
	}
	var l List
	var ok bool
	if l, ok = args[0].(List); !ok {
		return nil, ParsexSignErrorf("ginq select args error: except select from a list but %v", args[0])
	}
	rel := make(List, len(l))
	for idx, r := range l {
		call := L(sel.fun, r)
		data, err := Eval(env, call)
		if err != nil {
			return nil, err
		}
		rel[idx] = data
	}
	return Q(rel), nil
}

type GinWhere struct {
	expr interface{}
}

func NewGinWere(expr interface{}) GinWhere {
	return GinWhere{expr}
}

func (where GinWhere) Task(env Env, args ...interface{}) (Lisp, error) {
	if len(args) != 1 {
		return nil, ParsexSignErrorf("ginq where args error: except select from a list but %v", args)
	}
	var l List
	var ok bool
	if l, ok = args[0].(List); !ok {
		return nil, ParsexSignErrorf("ginq where args error: except select from a list but %v", args[0])
	}
	rel := List{}
	for _, r := range l {
		call := L(where.expr, Q(r))
		b, err := Eval(env, call)
		if err != nil {
			return nil, err
		}
		if t, ok := b.(bool); ok {
			if t {
				rel = append(rel, r)
			}
		} else {
			return nil, ParsexSignErrorf("ginq where exec error: except (%v %v) got a bool but %v",
				where.expr, r, b)
		}
	}
	return Q(rel), nil
}

type GinFields struct {
	funs List
}

func NewGinFields(args ...interface{}) GinFields {
	return GinFields{List(args)}
}

func (fs GinFields) Task(env Env, args ...interface{}) (Lisp, error) {
	if len(args) != 1 {
		return nil, ParsexSignErrorf("ginq fields args error: except fields from a list but %v", args)
	}
	row := make(List, len(fs.funs))
	for i, fun := range fs.funs {
		call := L(fun, Q(args[0]))
		data, err := Eval(env, call)
		if err != nil {
			return nil, err
		}
		row[i] = data
	}
	return Q(row), nil
}
