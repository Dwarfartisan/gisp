package gisp

import (
	"reflect"
	"testing"

	p "github.com/Dwarfartisan/goparsec"
)

func TestAtomParse0(t *testing.T) {
	data := "x"
	state := p.MemoryParseState(data)
	a, err := AtomParser(state)
	if err == nil {
		test := Atom{"x", Type{ANY, false}}
		if !reflect.DeepEqual(test, a) {
			t.Fatalf("except Atom{\"x\", ANY} but %v", a)
		}
	} else {
		t.Fatalf("except Atom{\"x\", ANY} but %v", err)
	}
}

func TestAtomParse1(t *testing.T) {
	data := "x::atom"
	state := p.MemoryParseState(data)
	a, err := AtomParser(state)
	if err == nil {
		test := Atom{"x", Type{ATOM, false}}
		if !reflect.DeepEqual(test, a) {
			t.Fatalf("except Atom{\"x\", ATOM} but %v", a)
		}
	} else {
		t.Fatalf("except Atom{\"x\", ATOM} but %v", err)
	}
}

func TestAtomParse2(t *testing.T) {
	data := "x::any"
	state := p.MemoryParseState(data)
	a, err := AtomParser(state)
	if err == nil {
		test := Atom{"x", Type{ANY, false}}
		if !reflect.DeepEqual(test, a) {
			t.Fatalf("except Atom{\"x\", ANY} but %v", a)
		}
	} else {
		t.Fatalf("except Atom{\"x\", ANY} but %v", err)
	}
}

func TestAtomParse3(t *testing.T) {
	data := "x::int"
	state := p.MemoryParseState(data)
	a, err := AtomParser(state)
	if err == nil {
		test := Atom{"x", Type{INT, false}}
		if !reflect.DeepEqual(test, a) {
			t.Fatalf("except Atom{\"x\", INT} but %v", a)
		}
	} else {
		t.Fatalf("except Atom{\"x\", INT} but %v", err)
	}
}
