package gisp

import (
	"fmt"
	"reflect"
	"sort"
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
				"count": GinCount{},
				"sum": TaskExpr(func(env Env, args ...interface{}) (Tasker, error) {
					if len(args) != 1 {
						return nil, ParsexSignErrorf("ginq sum args error: excpet one data list but: %v", args)
					}

					param := args[0]
					var l List
					var ok bool
					if l, ok = param.(List); !ok {
						return nil, ParsexSignErrorf("ginq sum args error: except a data List but: %v", param)
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
					if len(args) != 1 {
						return nil, ParsexSignErrorf("ginq max args error: excpet one data list but: %v", args)
					}

					param := args[0]
					var l List
					var ok bool
					if l, ok = param.(List); !ok {
						return nil, ParsexSignErrorf("ginq max args error: except a data List but: %v", param)
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
					if len(args) != 1 {
						return nil, ParsexSignErrorf("ginq min args error: excpet one data list but: %v", args)
					}

					param := args[0]
					var l List
					var ok bool
					if l, ok = param.(List); !ok {
						return nil, ParsexSignErrorf("ginq min args error: except a data List but: %v", param)
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
					if len(args) != 1 {
						return nil, ParsexSignErrorf("ginq avg args error: excpet one data list but: %v", args)
					}

					param := args[0]
					var l List
					var ok bool
					if l, ok = param.(List); !ok {
						return nil, ParsexSignErrorf("ginq avg args error: except a data List but: %v", param)
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
				"sums": LispExpr(func(env Env, args ...interface{}) (Lisp, error) {
					if len(args) != 1 {
						return nil, ParsexSignErrorf("ginq sum select args error: excpet one expression but: %v", args)
					}
					param, err := Eval(env, args[0])
					if err != nil {
						return nil, err
					}
					return Q(NewGinSumSelect(param)), nil
				}),
				"avgs": LispExpr(func(env Env, args ...interface{}) (Lisp, error) {
					if len(args) != 1 {
						return nil, ParsexSignErrorf("ginq avg select args error: excpet one expression but: %v", args)
					}
					param, err := Eval(env, args[0])
					if err != nil {
						return nil, err
					}
					return Q(NewGinAvgSelect(param)), nil
				}),
				"mins": LispExpr(func(env Env, args ...interface{}) (Lisp, error) {
					if len(args) != 1 {
						return nil, ParsexSignErrorf("ginq min select args error: excpet one expression but: %v", args)
					}
					param, err := Eval(env, args[0])
					if err != nil {
						return nil, err
					}
					return Q(NewGinMinSelect(param)), nil
				}),
				"maxs": LispExpr(func(env Env, args ...interface{}) (Lisp, error) {
					if len(args) != 1 {
						return nil, ParsexSignErrorf("ginq max select args error: excpet one expression but: %v", args)
					}
					param, err := Eval(env, args[0])
					if err != nil {
						return nil, err
					}
					return Q(NewGinMaxSelect(param)), nil
				}),
				"sort": TaskExpr(func(env Env, args ...interface{}) (Tasker, error) {
					if len(args) != 1 {
						return nil, ParsexSignErrorf("ginq sort args error: excpet one data list but: %v", args)
					}

					param := args[0]
					var l List
					var ok bool
					if l, ok = param.(List); !ok {
						return nil, ParsexSignErrorf("ginq sort args error: except a data List but: %v", param)
					}
					return func(env Env) (interface{}, error) {
						buf := make(List, len(l))
						copy(buf, l)
						s := GinSort{buf, env, nil}
						sort.Sort(&s)
						if s.err == nil {
							return buf, nil
						} else {
							return nil, s.err
						}
					}, nil
				}),
				"sortby": LispExpr(func(env Env, args ...interface{}) (Lisp, error) {
					if len(args) != 1 {
						return nil, ParsexSignErrorf("ginq sort by args error: excpet one expression but: %v", args)
					}
					param, err := Eval(env, args[0])
					if err != nil {
						return nil, err
					}
					return Q(NewGinSortBy(param)), nil
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
		call := L(group.by, Q(data))
		grp, err := Eval(env, call)
		if err != nil {
			return nil, fmt.Errorf("excpet group list:\n\t%v\nby %v but got error: \n\t%v",
				data, group.by, err)
		}
		flag := false
		for _, gr := range pool {
			if reflect.DeepEqual(gr[0], grp) {
				flag = true
				//group pool row
				grpr := gr[1].(List)
				gr[1] = append(grpr, data)
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
		call := L(group.group, Q(g[1]))
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
	return Selector{sel.fun, l}, nil
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
		call := L(fun, args[0])
		data, err := Eval(env, call)
		if err != nil {
			return nil, err
		}
		row[i] = data
	}
	return Q(row), nil
}

type Selector struct {
	fun  interface{}
	data List
}

func (sp Selector) Eval(env Env) (interface{}, error) {
	pool := make(List, len(sp.data))
	for idx, row := range sp.data {
		call := L(sp.fun, Q(row))
		rev, err := Eval(env, call)
		if err != nil {
			return nil, err
		}
		pool[idx] = rev
	}
	return pool, nil
}

type GinSumSelect struct {
	fun interface{}
}

func NewGinSumSelect(fun interface{}) GinSumSelect {
	return GinSumSelect{fun}
}

func (sel GinSumSelect) Task(env Env, args ...interface{}) (Lisp, error) {
	if len(args) != 1 {
		return nil, ParsexSignErrorf("ginq sum select data error: except select from a list but %v", args)
	}
	param, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	var l List
	var ok bool
	if l, ok = param.(List); !ok {
		return nil, ParsexSignErrorf("ginq sum select data error: except select from a list but %v", args[0])
	}
	return GinSumSelector{Selector{sel.fun, l}}, nil
}

type GinSumSelector struct {
	Selector
}

func (ss GinSumSelector) Eval(env Env) (interface{}, error) {
	p, err := ss.Selector.Eval(env)
	if err != nil {
		return nil, err
	}
	pool := p.(List)
	if len(pool) == 0 {
		return nil, nil
	}
	if len(pool) == 1 {
		return pool[0], nil
	}
	add, _ := env.Lookup("+")
	root := pool[0]
	for _, item := range pool[1:] {
		call := L(add, root, item)
		data, err := Eval(env, call)
		if err != nil {
			return nil, err
		}
		root = data
	}
	return root, nil
}

type GinMaxSelect struct {
	fun interface{}
}

func NewGinMaxSelect(fun interface{}) GinMaxSelect {
	return GinMaxSelect{fun}
}

func (sel GinMaxSelect) Task(env Env, args ...interface{}) (Lisp, error) {
	if len(args) != 1 {
		return nil, ParsexSignErrorf("ginq max select data error: except select from a list but %v", args)
	}
	param, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	var l List
	var ok bool
	if l, ok = param.(List); !ok {
		return nil, ParsexSignErrorf("ginq max select data error: except select from a list but %v", args[0])
	}
	return GinMaxSelector{Selector{sel.fun, l}}, nil
}

type GinMaxSelector struct {
	Selector
}

func (ms GinMaxSelector) Eval(env Env) (interface{}, error) {
	p, err := ms.Selector.Eval(env)
	if err != nil {
		return nil, err
	}
	pool := p.(List)
	if len(pool) == 0 {
		return nil, nil
	}
	if len(pool) == 1 {
		return pool[0], nil
	}
	lt, _ := env.Lookup("<")
	root := pool[0]
	for _, item := range pool[1:] {
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
}

type GinMinSelect struct {
	fun interface{}
}

func NewGinMinSelect(fun interface{}) GinMinSelect {
	return GinMinSelect{fun}
}

func (sel GinMinSelect) Task(env Env, args ...interface{}) (Lisp, error) {
	if len(args) != 1 {
		return nil, ParsexSignErrorf("ginq min select data error: except select from a list but %v", args)
	}
	param, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	var l List
	var ok bool
	if l, ok = param.(List); !ok {
		return nil, ParsexSignErrorf("ginq min select data error: except select from a list but %v", args[0])
	}
	return GinMinSelector{Selector{sel.fun, l}}, nil
}

type GinMinSelector struct {
	Selector
}

func (ms GinMinSelector) Eval(env Env) (interface{}, error) {
	p, err := ms.Selector.Eval(env)
	if err != nil {
		return nil, err
	}
	pool := p.(List)
	if len(pool) == 0 {
		return nil, nil
	}
	if len(pool) == 1 {
		return pool[0], nil
	}
	lt, _ := env.Lookup("<")
	root := pool[0]
	for _, item := range pool[1:] {
		call := L(lt, item, root)
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
}

type GinAvgSelect struct {
	fun interface{}
}

func NewGinAvgSelect(fun interface{}) GinAvgSelect {
	return GinAvgSelect{fun}
}

func (sel GinAvgSelect) Task(env Env, args ...interface{}) (Lisp, error) {
	if len(args) != 1 {
		return nil, ParsexSignErrorf("ginq avg select data error: except select from a list but %v", args)
	}
	param, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	var l List
	var ok bool
	if l, ok = param.(List); !ok {
		return nil, ParsexSignErrorf("ginq avg select data error: except select from a list but %v", args[0])
	}
	return GinAvgSelector{Selector{sel.fun, l}}, nil
}

type GinAvgSelector struct {
	Selector
}

func (as GinAvgSelector) Eval(env Env) (interface{}, error) {
	p, err := as.Selector.Eval(env)
	if err != nil {
		return nil, err
	}
	pool := p.(List)
	if len(pool) == 0 {
		return nil, nil
	}
	if len(pool) == 1 {
		return pool[0], nil
	}
	add, _ := env.Lookup("+")
	root := pool[0]
	for _, item := range pool[1:] {
		call := L(add, root, item)
		data, err := Eval(env, call)
		if err != nil {
			return nil, err
		}
		root = data
	}
	div, _ := env.Lookup("/")
	call := L(div, root, len(pool))
	rev, err := Eval(env, call)
	if err != nil {
		return nil, err
	}
	return rev, nil
}

type GinCount struct {
}

func (c GinCount) Task(env Env, args ...interface{}) (Lisp, error) {
	if len(args) != 1 {
		return nil, ParsexSignErrorf("ginq count data error: except a list but %v", args)
	}
	param, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	var l List
	var ok bool
	if l, ok = param.(List); !ok {
		return nil, ParsexSignErrorf("ginq count data error: except count a list but %v", args[0])
	}
	return Q(len(l)), nil
}

type GinSort struct {
	List
	env Env
	err error
}

func (ls *GinSort) Less(x, y int) bool {
	less, _ := ls.env.Lookup("<")
	call := L(less, Q(ls.List[x]), Q(ls.List[y]))
	b, err := Eval(ls.env, call)
	if err != nil {
		ls.err = err
		return false
	}
	if is, ok := b.(bool); ok {
		return is
	} else {
		//ls.err = fmt.Errorf("except (less x y) as (< %v %v) return true or false but error: %v", err)
		ls.err = err
		return false
	}
}
func (ls *GinSort) Len() int {
	return len(ls.List)
}
func (ls *GinSort) Swap(i, j int) {
	tmp := ls.List[i]
	ls.List[i] = ls.List[j]
	ls.List[j] = tmp
}

type GinSortBy struct {
	fun interface{}
}

func NewGinSortBy(fun interface{}) GinSortBy {
	return GinSortBy{fun}
}

func (gsb GinSortBy) Task(env Env, args ...interface{}) (Lisp, error) {
	if len(args) != 1 {
		return nil, ParsexSignErrorf("ginq sort data error: except sort one list but %v", args)
	}
	// param, err := Eval(env, args[0])
	// if err != nil {
	// 	return nil, err
	// }
	param := args[0]
	var l List
	var ok bool
	if l, ok = param.(List); !ok {
		return nil, ParsexSignErrorf("ginq sort data error: except sort a list but %v", args[0])
	}
	buf := make(List, len(l))
	copy(buf, l)
	return &GinSortListBy{buf, env, gsb.fun, nil}, nil
}

type GinSortListBy struct {
	List
	env Env
	fun interface{}
	err error
}

func (gsl *GinSortListBy) Less(x, y int) bool {
	call := L(gsl.fun, Q(gsl.List[x]), Q(gsl.List[y]))
	b, err := Eval(gsl.env, call)
	if err != nil {
		gsl.err = err
		return false
	}
	if is, ok := b.(bool); ok {
		return is
	} else {
		//ls.err = fmt.Errorf("except (less x y) as (< %v %v) return true or false but error: %v", err)
		gsl.err = err
		return false
	}
}
func (gsl *GinSortListBy) Len() int {
	return len(gsl.List)
}
func (gsl *GinSortListBy) Swap(i, j int) {
	tmp := gsl.List[i]
	gsl.List[i] = gsl.List[j]
	gsl.List[j] = tmp
}
func (gsl *GinSortListBy) Eval(env Env) (interface{}, error) {
	sort.Sort(gsl)
	if gsl.err == nil {
		return gsl.List, nil
	} else {
		return nil, gsl.err
	}
}
