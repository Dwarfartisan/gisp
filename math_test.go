package gisp

import (
	px "github.com/Dwarfartisan/goparsec/parsex"
	"testing"
)

func TestAddx0(t *testing.T) {
	var data = []interface{}{0, 1, 2, 3, 4, 5, 6}
	st := px.NewStateInMemory(data)
	s, err := addx(st)
	if err != nil {
		t.Fatalf("except error is nil but %v", err)
	}
	if s.(Int) != 21 {
		t.Fatalf("except sum 0~6 is 21 but got %v", s)
	}
}

func TestAddx1(t *testing.T) {
	var data = []interface{}{0, 1, 2, 3.14, 4, 5, 6}
	st := px.NewStateInMemory(data)
	s, err := addx(st)
	if err != nil {
		t.Fatalf("except error is nil but %v", err)
	}
	if s.(Float) != 21.14 {
		t.Fatalf("except sum 0, 1, 2, 3.14, 4, 5, 6 is 21.14 but got %v", s)
	}
}

func TestAddExpr(t *testing.T) {
	gisp := NewGisp(map[string]Toolbox{
		"axioms": Axiom,
		"props":  Propositions,
	})
	adds, err := gisp.Parse("+")
	if err != nil {
		t.Fatalf("except add operator but error %v", err)
	}
	var expr = []interface{}{adds, 0, 1, 2, 3.14, 4, 5, 6}
	ret, err := gisp.Eval(List(expr))
	if err != nil {
		t.Fatalf("except add data %v but got error %v", expr[1:], err)
	}
	if ret.(Float) != 21.14 {
		t.Fatalf("except sum 0, 1, 2, 3.14, 4, 5, 6 is 21.14 but got %v", ret)
	}
}

func TestMulExpr(t *testing.T) {
	gisp := NewGisp(map[string]Toolbox{
		"axioms": Axiom,
		"props":  Propositions,
	})
	mulx, err := gisp.Parse("*")
	if err != nil {
		t.Fatalf("except add operator but error %v", err)
	}
	var expr = L(mulx, 1, 2, 3.14, 4, 5, 6)
	ret, err := gisp.Eval(expr)
	if err != nil {
		t.Fatalf("except add data %v but got error %v", expr[1], err)
	}
	if ret.(Float) != 753.6 {
		t.Fatalf("except multi 1, 2, 3.14, 4, 5, 6 is 753.6 but got %v", ret)
	}
	expr = L(mulx, 2, 3, 4, 5, 6)
	ret, err = gisp.Eval(expr)
	if err != nil {
		t.Fatalf("except add data %v but got error %v", expr[1:], err)

	}
	if ret.(Int) != 720 {
		t.Fatalf("except multi %v is %d but got %v", expr[1:], 720, ret)
	}
}
