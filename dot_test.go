package gisp

import (
	//"fmt"
	p "github.com/Dwarfartisan/goparsec"
	"reflect"
	"testing"
	"time"
)

func TestDotTime(t *testing.T) {
	now := time.Now()
	g, _ := NewGisp(map[string]Toolbox{
		"axioms": Axiom,
		"props":  Propositions,
	})
	slot := VarSlot(ANYOPTION)
	slot.Set(now)
	g.Defvar("now", slot)
	year := Int(now.Year())
	y, err := g.Parse("(now.Year)") //g.Eval(List{AA("now.Year")})
	if err != nil {
		t.Fatalf("except (now.Year) equal to now.Year() as %v but got error %v", year, err)
	}
	if !reflect.DeepEqual(year, y) {
		t.Fatalf("except (now.Year) equal to now.Year() but got %v and %v", year, y)
	}
}

func TestDotParser(t *testing.T) {
	data := "now.Year"
	st := p.MemoryParseState(data)
	re, err := DotParser(st)
	if err != nil {
		t.Fatalf("except a Dot but error %v", err)
	}
	t.Log(re)
}
