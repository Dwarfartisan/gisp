package gisp

import (
	"reflect"
	"testing"
)

func TestTypeFound(t *testing.T) {
	m := money{9.99, "USD"}
	g, _ := NewGisp(map[string]Toolbox{
		"axioms": Axiom,
		"props":  Propositions,
	})
	g.DefAs("money", reflect.TypeOf(m))
	_, err := g.Parse("(var bill::money)")
	if err != nil {
		t.Fatalf("except define a money var but error: %v", err)
	}
}
