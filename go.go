package gisp

type Go struct {
	Task
}

func (this Go) Eval(env Env) (interface{}, error) {
	go this.Task.Eval(env)
	return nil, nil
}
