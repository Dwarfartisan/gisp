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
	g, err := NewGisp(map[string]Toolbox{
		"axioms": Axiom,
		"props":  Propositions,
	})
	if err != nil {
		t.Fatalf("except gisp parser but %v", err)
	}
	addx, err := g.Parse("+")
	if err != nil {
		t.Fatalf("except add operator but error %v", err)
	}
	var data = []interface{}{0, 1, 2, 3.14, 4, 5, 6}
	ret, err := addx.(Element)(data...)
	if err != nil {
		t.Fatalf("except add data %v but got error %v", data, err)
	}
	if ret.(Float) != 21.14 {
		t.Fatalf("except sum 0, 1, 2, 3.14, 4, 5, 6 is 21.14 but got %v", ret)
	}
}

func TestMulExpr(t *testing.T) {
	g, err := NewGisp(map[string]Toolbox{
		"axioms": Axiom,
		"props":  Propositions,
	})
	if err != nil {
		t.Fatalf("except gisp parser but %v", err)
	}
	mulx, err := g.Parse("*")
	if err != nil {
		t.Fatalf("except add operator but error %v", err)
	}
	var data = []interface{}{1, 2, 3.14, 4, 5, 6}
	ret, err := mulx.(Element)(data...)
	if err != nil {
		t.Fatalf("except add data %v but got error %v", data, err)
	}
	if ret.(Float) != 753.6 {
		t.Fatalf("except multi 1, 2, 3.14, 4, 5, 6 is 753.6 but got %v", ret)
	}
	data = []interface{}{2, 3, 4, 5, 6}
	ret, err = mulx.(Element)(data...)
	if err != nil {
		t.Fatalf("except add data %v but got error %v", data, err)

	}
	if ret.(Int) != 720 {
		t.Fatalf("except multi %v is %d but got %v", data, 720, ret)
	}
}
