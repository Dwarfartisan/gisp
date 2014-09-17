package gisp

import (
	"reflect"
	"testing"
	tm "time"
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

func TestParseExt(t *testing.T) {
	g := NewGispWith(
		map[string]Toolbox{"axiom": Axiom, "props": Propositions},
		map[string]Toolbox{"time": Time})
	gisp := *g
	ret, err := gisp.Parse("(time.now)")
	if err != nil {
		t.Fatalf("except got time.Now() but error: %v", err)
	}
	if now, ok := ret.(tm.Time); ok {
		t.Logf("got now time is %v", now)
	} else {
		t.Fatalf("except got now time but %v", now)
	}
}

func TestParseCallToolkitFunction(t *testing.T) {
	g := NewGispWith(
		map[string]Toolbox{"axiom": Axiom, "props": Propositions},
		map[string]Toolbox{"time": Time})
	gisp := *g
	ret, err := gisp.Parse(`(time.parseDuration "24h")`)
	if err != nil {
		t.Fatalf("except got time.Duration 24 hours but error: %v", err)
	}
	dur, err := tm.ParseDuration("24h")
	if err != nil {
		t.Fatalf("except got time.Duration 24 hours but error: %v", err)
	}
	if !reflect.DeepEqual(dur, ret) {
		t.Fatalf("except got time.Duration 24 hours but got: %v", ret)
	}
	t.Logf("parse duration \"24h\" got %v\n", ret)
}
