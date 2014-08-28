package gisp

import (
	p "github.com/Dwarfartisan/goparsec"
	"testing"
)

func TestFloatParser0(t *testing.T) {
	data := "0.012"
	st := p.MemoryParseState(data)
	o, err := FloatParser(st)
	if err != nil {
		t.Fatalf("except a Float but error %v", err)
	}
	if f, ok := o.(Float); ok {
		if f != Float(0.012) {
			t.Fatalf("except a Float 0.012 but %v", f)
		}
	} else {
		t.Fatalf("except Float but %v", o)
	}
}

func TestFloatParser1(t *testing.T) {
	data := "3.1415926"
	st := p.MemoryParseState(data)
	o, err := FloatParser(st)
	if err != nil {
		t.Fatalf("except a Float but error %v", err)
	}
	if f, ok := o.(Float); ok {
		if f != Float(3.1415926) {
			t.Fatalf("except a Float 3.1415926 but %v", f)
		}
	} else {
		t.Fatalf("except Float but %v", o)
	}
}

func TestFloatParser2(t *testing.T) {
	data := "234.0"
	st := p.MemoryParseState(data)
	o, err := FloatParser(st)
	if err != nil {
		t.Fatalf("except a Float but error %v", err)
	}
	if f, ok := o.(Float); ok {
		if f != Float(234) {
			t.Fatalf("except a Float 234.0 but %v", f)
		}
	} else {
		t.Fatalf("except Float but %v", o)
	}
}

func TestFloatParser3(t *testing.T) {
	data := ".5"
	st := p.MemoryParseState(data)
	o, err := FloatParser(st)
	if err != nil {
		t.Fatalf("except a Float but error %v", err)
	}
	if f, ok := o.(Float); ok {
		if f != Float(0.5) {
			t.Fatalf("except a Float 0.5 but %v", f)
		}
	} else {
		t.Fatalf("except Float but %v", o)
	}
}

func TestFloatParser4(t *testing.T) {
	data := "f234.0"
	st := p.MemoryParseState(data)
	o, err := FloatParser(st)
	if err == nil {
		t.Fatalf("except a Float parse error but got %v", o)
	}
}

func TestFloatParser5(t *testing.T) {
	data := "234"
	st := p.MemoryParseState(data)
	o, err := FloatParser(st)
	if err == nil {
		t.Fatalf("except a Float parse error but got %v", o)
	}
}
