package gisp

import (
	"testing"
)

func TestParseString(t *testing.T) {
	g := NewGisp(map[string]Toolbox{})
	gisp := *g
	data := `"I'm a string"`
	str, err := gisp.Parse(data)
	if err != nil {
		t.Fatalf("except string but error: %v", err)
	}
	if str.(string) != `I'm a string` {
		t.Fatalf("except got \"I'm a string\" but \"%v\"", str)
	}
}

func TestParseFloat(t *testing.T) {
	g := NewGisp(map[string]Toolbox{})
	gisp := *g
	data := "3.14"
	ret, err := gisp.Parse(data)
	if err != nil {
		t.Fatalf("except Float(3.14) but error: %v", err)
	}
	if ret.(Float) != Float(3.14) {
		t.Fatalf("except got Float(3.14) but %v", ret)
	}
}
