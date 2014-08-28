package gisp

import (
	p "github.com/Dwarfartisan/goparsec"
	"testing"
)

func TestIntParser0(t *testing.T) {
	data := "12"
	st := p.MemoryParseState(data)
	o, err := IntParser(st)
	if err != nil {
		t.Fatalf("except a Int but error %v", err)
	}
	if i, ok := o.(Int); ok {
		if i != Int(12) {
			t.Fatalf("except a Int 12 but %v", i)
		}
	} else {
		t.Fatalf("except Int but %v", o)
	}
}

func TestIntParser1(t *testing.T) {
	data := "i234"
	st := p.MemoryParseState(data)
	o, err := IntParser(st)
	if err == nil {
		t.Fatalf("except a Int parse error but got %v", o)
	}
}

func TestIntParser2(t *testing.T) {
	data := ".234"
	st := p.MemoryParseState(data)
	o, err := IntParser(st)
	if err == nil {
		t.Fatalf("except a Float parse error but got %v", o)
	}
}

func TestIntParser3(t *testing.T) {
	data := "3.14"
	st := p.MemoryParseState(data)
	o, err := p.Bind_(IntParser, p.Eof)(st)
	if err == nil {
		t.Fatalf("except a Float parse error but got %v", o)
	}
}
