package gisp

import (
	p "github.com/Dwarfartisan/goparsec"
	"reflect"
	"testing"
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
	g, _ := NewGisp(map[string]Toolbox{
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
