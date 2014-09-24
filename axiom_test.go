package gisp

import (
	"reflect"
	"testing"
	// tm "time"
)

func TestQuoteFound(t *testing.T) {
	ret, ok := Axiom.Lookup("quote")
	if !ok || ret == nil {
		t.Fatalf("except found quote in axioms")
	}
}

func TestQuoteCall(t *testing.T) {
	gisp := NewGisp(map[string]Toolbox{
		"axioms": Axiom,
	})
	list := L(Int(1), Int(1), Int(2), Int(3), Int(5), Int(8), Int(13), Int(21))
	q, ok := gisp.Lookup("quote")
	if !ok {
		t.Fatalf("except found quote in axioms")
	}
	quote := q.(Functor)
	quoted, err := quote.Task(gisp, list)
	if err != nil {
		t.Fatalf("except quote task but %v", err)
	}
	fb, err := quoted.Eval(gisp)
	if err != nil {
		t.Fatalf("except quote the list but %v", err)
	}
	fbq := Quote{list}
	if !reflect.DeepEqual(fb, fbq) {
		t.Fatalf("except quote (1 1 2 3 5 8 13 21) got %v but %v", fbq, fb)
	}
}

func TestQuoteEval(t *testing.T) {
	g := NewGisp(map[string]Toolbox{
		"axioms": Axiom,
	})
	gisp := *g
	list := L(Int(1), Int(1), Int(2), Int(3), Int(5), Int(8), Int(13), Int(21))
	re, err := gisp.Eval(Quote{list})
	if err != nil {
		t.Fatalf("except eval quote got a list but error %v", err)
	}
	if !reflect.DeepEqual(list, re) {
		t.Fatalf("except Eval (quote (1 1 2 3 5 8 13 21)) got %v but %v", list, re)
	}
}

func TestGetEval(t *testing.T) {
	gisp := NewGisp(map[string]Toolbox{
		"axioms": Axiom,
	})
	_, err := gisp.Parse("(var pi 3.14)")
	if err != nil {
		t.Fatalf("except var pi as 3.14 but error: %v", err)
	}
	pi, err := gisp.Eval(Atom{"pi", FLOATMUST})
	if err != nil {
		t.Fatalf("except got pi is 3.14 but error: %v", err)
	}
	if !reflect.DeepEqual(pi, Float(3.14)) {
		t.Fatalf("except got pi is 3.14 but %v", pi)
	}
}

func TestVarEval(t *testing.T) {
	g := NewGisp(map[string]Toolbox{
		"axioms": Axiom,
	})
	gisp := *g
	_, err := gisp.Parse("(var pi 3.14)")
	if err != nil {
		t.Fatalf("except var pi as 3.14 but error: %v", err)
	}
	pi, err := gisp.Eval(Atom{"pi", FLOATMUST})
	if err != nil {
		t.Fatalf("except got pi is 3.14 but error: %v", err)
	}
	if !reflect.DeepEqual(pi, Float(3.14)) {
		t.Fatalf("except got pi is 3.14 but %v", pi)
	}
}

func TestSetBracket(t *testing.T) {
	gisp := NewGisp(map[string]Toolbox{
		"axioms": Axiom,
	})
	err := gisp.DefAs("box", map[string]interface{}{
		"a": "abc",
		"b": "bcd",
		"c": 3.14,
		"d": 99,
	})
	if err != nil {
		t.Fatalf("except var pi as 3.14 but error: %v", err)
	}
	pi, err := gisp.Parse("box[\"c\"]")
	if err != nil {
		t.Fatalf("except got pi is 3.14 but error: %v", err)
	}
	if !reflect.DeepEqual(pi, 3.14) {
		t.Fatalf("except got pi is 3.14 but %v", pi)
	}
	_, err = gisp.Parse("(set 'box[\"c\"] 3.14)")
	if err != nil {
		t.Fatalf("except got pi is 3.14 but error: %v", err)
	}
	pi, err = gisp.Parse("box[\"c\"]")
	if err != nil {
		t.Fatalf("except got pi is 3.14 but error: %v", err)
	}
	if !reflect.DeepEqual(pi, Float(3.14)) {
		t.Fatalf("except got pi is 3.14 but %v", pi)
	}
}
