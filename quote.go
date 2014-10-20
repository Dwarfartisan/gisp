package gisp

type Quote struct {
	Lisp interface{}
}

func (this Quote) Eval(env Env) (interface{}, error) {
	return this.Lisp, nil
}

func Q(x interface{}) Quote {
	return Quote{x}
}

func QL(args ...interface{}) Quote {
	return Q(L(args...))
}
