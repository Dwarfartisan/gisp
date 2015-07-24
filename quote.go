package gisp

// Quote 定义了 Lisp Quote
type Quote struct {
	Lisp interface{}
}

// Eval 实现了 Eval 行为
func (this Quote) Eval(env Env) (interface{}, error) {
	return this.Lisp, nil
}

// Q 得到一个 Quote
func Q(x interface{}) Quote {
	return Quote{x}
}

// QL 得到一个 Quote 后的列表
func QL(args ...interface{}) Quote {
	return Q(L(args...))
}
