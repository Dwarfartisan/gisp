package gisp

import (
	px "github.com/Dwarfartisan/goparsec/parsex"
	"reflect"
	"testing"
)

// declare a struct for overload test
type money struct {
	Amount   Float
	Currency string
}

func (m money) Mul(ratio Float) money {
	return money{m.Amount * ratio, m.Currency}
}
func (m money) EqualTo(x money) bool {
	return m.Amount == x.Amount && m.Currency == x.Currency
}

type mrMul struct {
	argsParser px.Parser
}

func mrmul() mrMul {
	mchecker := func(st px.ParsexState) (interface{}, error) {
		x, err := st.Next(px.Always)
		if err != nil {
			return nil, err
		}
		switch m := x.(type) {
		case money:
			return m, nil
		default:
			return nil, moneySignError(x)
		}
	}
	fchecker := func(st px.ParsexState) (interface{}, error) {
		x, err := st.Next(px.Always)
		if err != nil {
			return nil, err
		}
		switch f := x.(type) {
		case Float:
			return f, nil
		default:
			return nil, floatSignError(x)
		}
	}
	return mrMul{px.Union(mchecker, px.Many1(fchecker), px.Eof)}
}

func (mrm mrMul) Task(env Env, args ...interface{}) (Lisp, error) {
	params, err := Evals(env, args...)
	if err != nil {
		return nil, err
	}
	st := px.NewStateInMemory(params)
	data, err := mrm.argsParser(st)
	if err != nil {
		return nil, err
	}
	return TaskBox{func(env Env) (interface{}, error) {
		vals := data.([]interface{})
		m := vals[0].(money)
		for _, r := range vals[1].([]interface{}) {
			m = m.Mul(r.(Float))
		}
		return m, nil
	}}, nil
}

func moneySignError(value interface{}) error {
	return TypeSignError{Type{reflect.TypeOf((*money)(nil)).Elem(), false}, value}
}

func floatSignError(value interface{}) error {
	return TypeSignError{Type{reflect.TypeOf((*Float)(nil)).Elem(), false}, value}
}

func TestMoneyMul(t *testing.T) {
	in := money{Float(30.9), "CNY"}
	ratio := Float(0.8)
	out := in.Mul(ratio)
	g, err := NewGisp(map[string]Toolbox{
		"axioms": Axiom,
		"props":  Propositions,
	})
	if err != nil {
		t.Fatalf("except gisp parser but %v", err)
	}
	g.Defun("*", mrmul())
	mulx, ok := g.Lookup("*")
	if !ok {
		t.Fatalf("except got overloaded function *")
	}

	ret, err := g.Eval(List{mulx, in, ratio})
	if err != nil {
		t.Fatalf("except %v * %v is %v but error %v", in, ratio, out, err)
	}
	if !reflect.DeepEqual(ret, out) {
		t.Fatalf("except %v * %v is %v but %v", in, ratio, out, ret)
	}
}

func TestMulAutoOverload(t *testing.T) {
	in := Float(30.9)
	ratio := Float(0.8)
	out := in * ratio
	g, err := NewGisp(map[string]Toolbox{
		"axioms": Axiom,
		"props":  Propositions,
	})
	if err != nil {
		t.Fatalf("except gisp parser but %v", err)
	}
	g.Defun("*", mrmul())
	mulx, ok := g.Lookup("*")
	if !ok {
		t.Fatalf("except got overloaded function *")
	}

	ret, err := g.Eval(List{mulx, in, ratio})
	if err != nil {
		t.Fatalf("except %v * %v is %v but error %v", in, ratio, out, err)
	}
	if !reflect.DeepEqual(ret, out) {
		t.Fatalf("except %v * %v is %v but %v", in, ratio, out, ret)
	}
}
