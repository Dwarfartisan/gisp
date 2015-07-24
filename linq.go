package gisp

// T 是 interface{} 的简写
type T interface{}

// Linq 定义了 linq 接口
type Linq interface {
	From(input T) Linq
	Range(start, count int)
	All(f func(T) (bool, error)) (all bool, err error)
	Any() (exists bool, err error)
	AnyWith(f func(T) (bool, error)) (exists bool, err error)
	Average() (avg float64, err error)
	Count() (count int, err error)
	CountBy(f func(T) (bool, error)) (c int, err error)
	Distinct() Linq
	DistinctBy(f func(T, T) (bool, error)) Linq
	ElementAt(i int) (elem T, found bool, err error)
	Except(inputSlice T) Linq
	First() (elem T, found bool, err error)
	FirstBy(f func(T) (bool, error)) (elem T, found bool, err error)
	GroupBy(keySelector func(T) T, valueSelector func(T) T) (map[T][]T, error)
	GroupJoin(innerSlice T, outerKeySelector func(T) T, innerKeySelector func(T) T, resultSelector func(outer T, inners []T) T) Linq
	Intersect(inputSlice T) Linq
	Join(innerSlice T, outerKeySelector func(T) T, innerKeySelector func(T) T, resultSelector func(outer T, inner T) T) Linq
	Last() (elem T, found bool, err error)
	LastBy(f func(T) (bool, error)) (elem T, found bool, err error)
	Max() (max T, err error)
	Min() (min T, err error)
	OrderBy(less func(this T, that T) bool) Linq
	Order(func(x, y T) (bool, error)) Linq
	Results() (List, error)
	Reverse() Linq
	Select(f func(T) (T, error)) Linq
	Single(f func(T) (bool, error)) (single T, err error)
	Skip(n int) Linq
	SkipWhile(f func(T) (bool, error)) Linq
	Sum() (sum float64, err error)
	Take(n int) Linq
	TakeWhile(f func(T) (bool, error)) Linq
	Union(inputSlice T) Linq
	Where(f func(T) (bool, error)) Linq
}
