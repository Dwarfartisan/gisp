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

type moneyMulti struct {
}

func (mm moneyMulti) Task(args ...interface{}) (Lisp, error) {
	err := mmSignChecker(args...)
	if err != nil {
		return nil, err
	}
	return mrMul{args[0], args[1:]}, nil
}

func mmSignChecker(args ...interface{}) error {
	st := px.NewStateInMemory(args)
	mchecker := func(st px.ParsexState) (interface{}, error) {
		x, err := st.Next(px.Always)
		if err != nil {
			return nil, err
		}
		switch m := x.(type) {
		case money:
			return nil, nil
		case Atom:
			if m.Type.Type == ANY {
				return nil, nil
			} else {
				return nil, moneySignError(m)
			}
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
			return nil, nil
		case Atom:
			if f.Type.Type == FLOAT {
				return nil, nil
			} else {
				return nil, floatSignError(f)
			}
		default:
			return nil, floatSignError(x)
		}
	}
	_, err := px.Binds_(mchecker, px.Many1(fchecker), px.Eof)(st)
	return err
}

type mrMul struct {
	m  interface{}
	rs []interface{}
}

func (mrm mrMul) Eval(env Env) (interface{}, error) {
	m, err := eval(env, mrm.m)
	if err != nil {
		return nil, err
	}
	for _, r := range mrm.rs {
		f, err := eval(env, r)
		if err != nil {
			return nil, err
		}
		m = m.(money).Mul(f.(Float))
	}
	return m, nil
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
	g.Defun("*", moneyMulti{})
	mulx, ok := g.Lookup("*")
	if !ok {
		t.Fatalf("except got overloaded function *")
	}

	ret, err := g.Eval(List{mulx, in, Float(0.8)})
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
	g.Defun("*", moneyMulti{})
	mulx, ok := g.Lookup("*")
	if !ok {
		t.Fatalf("except got overloaded function *")
	}

	ret, err := g.Eval(List{mulx, in, Float(0.8)})
	if err != nil {
		t.Fatalf("except %v * %v is %v but error %v", in, ratio, out, err)
	}
	if !reflect.DeepEqual(ret, out) {
		t.Fatalf("except %v * %v is %v but %v", in, ratio, out, ret)
	}
}
