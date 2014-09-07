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

type Box map[string]interface{}

func (b Box) Get(name string) interface{} {
	return b[name]
}

func TestDotMap(t *testing.T) {
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
	g, _ := NewGisp(map[string]Toolbox{
		"axioms": Axiom,
		"props":  Propositions,
	})
	g.Defvar("box", VarSlot(ANYOPTION))
	g.Setvar("box", box)
	c, err := g.Parse(`(box.Get "c")`)
	if err != nil {
		t.Fatalf("excpet got b but error %v", err)
	}
	if !reflect.DeepEqual(Quote{c}, box["c"]) {
		t.Fatalf("except %v but got %v", box["c"], c)
	}

}
