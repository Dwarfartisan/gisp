package gisp

import (
	"testing"

	p "github.com/Dwarfartisan/goparsec"
)

func TestingBoolParse0(t *testing.T) {
	data := "true"
	st := p.MemoryParseState(data)
	o, err := BoolParser(st)
	if err != nil {
		t.Fatalf("except bool but error %v", err)
	}
	if b, ok := o.(Bool); ok {
		if !b {
			t.Fatalf("except bool true but %v", b)
		}
	} else {
		t.Fatalf("excpet bool but %v", o)
	}
}

func TestingBoolParse1(t *testing.T) {
	data := "false"
	st := p.MemoryParseState(data)
	o, err := BoolParser(st)
	if err != nil {
		t.Fatalf("except bool but error %v", err)
	}
	if b, ok := o.(bool); ok {
		if !b {
			t.Fatalf("except bool true but %v", b)
		}
	} else {
		t.Fatalf("excpet bool but %v", o)
	}
}
