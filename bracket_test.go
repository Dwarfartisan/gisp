package gisp

import (
	"reflect"
	"testing"

	p "github.com/Dwarfartisan/goparsec"
)

func TestBracketParser(t *testing.T) {
	data := "dict[\"meta\"]"
	st := p.MemoryParseState(data)
	re, err := p.Bind(AtomParser, BracketSuffixParser)(st)
	if err != nil {
		t.Fatalf("except a Dot but error %v", err)
	}
	t.Log(re)
}

func TestBracketBasic(t *testing.T) {
	g := NewGisp(map[string]Toolbox{
		"axioms": Axiom,
		"props":  Propositions,
	})
	g.Defvar("entry", VarSlot(ANYOPTION))
	g.Setvar("entry", map[string]interface{}{
		"meta": "meta data",
	})
	data, err := g.Parse(`entry["meta"]`)
	if err != nil {
		t.Fatalf("excpet got meta from entry but error: %v", err)
	}
	if !reflect.DeepEqual(data, "meta data") {
		t.Fatalf(`excpet got "meta data" from entry["meta"] but got %v`, data)
	}
}

func TestBracketExpression(t *testing.T) {
	g := NewGisp(map[string]Toolbox{
		"axioms": Axiom,
		"props":  Propositions,
	})
	l := L(3.14, 1.414, 1.735, 2.718)
	g.DefAs("list", l)
	pi, err := g.Parse("([0] list)")
	if err != nil {
		t.Fatalf("except got pi but error: %v", err)
	}
	if pi.(Float) != Float(3.14) {
		t.Fatalf("excpet got pi as float 3.14 but %v", pi)
	}
}

func TestBracketExpressionMap(t *testing.T) {
	box := Box{
		"a": Quote{AA("a")},
		"b": Quote{AA("bb")},
		"c": Quote{AA("ccc")},
	}
	bv := reflect.ValueOf(box)
	get := bv.MethodByName("Get")
	res := get.Call([]reflect.Value{reflect.ValueOf("b")})
	if !reflect.DeepEqual(res[0].Interface(), box["b"]) {
		t.Fatalf("except %v but got %v", box["b"], res[0].Interface())
	}
	g := NewGisp(map[string]Toolbox{
		"axioms": Axiom,
		"props":  Propositions,
	})
	g.DefAs("box", box)
	c, err := g.Parse(`(["c"] box)`)
	if err != nil {
		t.Fatalf("excpet got b but error %v", err)
	}
	if !reflect.DeepEqual(c, box["c"]) {
		t.Fatalf("except %v but gt %v", box["c"], c)
	}

}
