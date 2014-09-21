package gisp

import (
	"testing"
	tm "time"
)

func TestLessTime(t *testing.T) {
	gisp := NewGispWith(
		map[string]Toolbox{"axiom": Axiom, "props": Propositions},
		map[string]Toolbox{"time": Time})
	gisp.DefAs("start", tm.Now())
	_, err := gisp.Parse("(var stop (time.now))")
	if err != nil {
		t.Fatalf("excpet def var now but error: %v", err)
	}
	ret, err := gisp.Parse("(< start stop)")
	start, _ := gisp.Parse("start")
	stop, _ := gisp.Parse("stop")
	if err != nil {
		t.Fatalf("excpet %v < %v but error: %v", start, stop, err)
	}
	if !ret.(bool) {
		t.Fatalf("excpet %v < %v", start, stop)
	}
	t.Logf("got %v < %v", start, stop)
}

func TestLessNumber(t *testing.T) {
	gisp := NewGispWith(
		map[string]Toolbox{"axiom": Axiom, "props": Propositions},
		map[string]Toolbox{"time": Time})
	gisp.DefAs("x", 15)
	gisp.Parse("(var y 3.14)")
	ret, err := gisp.Parse("(< x y)")
	x, _ := gisp.Parse("x")
	y, _ := gisp.Parse("y")
	if err != nil {
		t.Fatalf("excpet not %v < %v but error: %v", x, y, err)
	}
	if ret.(bool) {
		t.Fatalf("excpet not %v < %v", x, y)
	}
	t.Logf("got not %v < %v", x, y)
}

func TestLessString(t *testing.T) {
	gisp := NewGispWith(
		map[string]Toolbox{"axiom": Axiom, "props": Propositions},
		map[string]Toolbox{"time": Time})
	err := gisp.DefAs("a", "to be or not to be")
	if err != nil {
		t.Fatalf("excpet define a but error: %v", err)
	}
	_, err = gisp.Parse(`(var b "to be or not too be")`)
	if err != nil {
		t.Fatalf("excpet define b but error: %v", err)
	}
	ret, err := gisp.Parse("(< a b)")
	a, _ := gisp.Parse("a")
	b, _ := gisp.Parse("b")
	if err != nil {
		t.Fatalf("excpet \"%v\" < \"%v\" but error: %v", a, b, err)
	}
	if !ret.(bool) {
		t.Fatalf("excpet %v < %v", a, b)
	}
	t.Logf("got %v < %v", a, b)
}

func TestLessList(t *testing.T) {
	gisp := NewGispWith(
		map[string]Toolbox{"axiom": Axiom, "props": Propositions},
		map[string]Toolbox{"time": Time})
	gisp.DefAs("x", L(33, 29, "abc"))
	_, err := gisp.Parse("(var y '(33 45 \"def\"))")
	if err != nil {
		t.Fatalf("excpet define b but error: %v", err)
	}
	x, _ := gisp.Parse("x")
	y, _ := gisp.Parse("y")
	ret, err := gisp.Parse("(< x y)")
	if err != nil {
		t.Fatalf("excpet not %v < %v but error: %v", x, y, err)
	}
	if !ret.(bool) {
		t.Fatalf("excpet not %v < %v", x, y)
	}
	t.Logf("got not %v < %v", x, y)
}

func TestLessOrEqualList(t *testing.T) {
	gisp := NewGispWith(
		map[string]Toolbox{"axiom": Axiom, "props": Propositions},
		map[string]Toolbox{"time": Time})
	gisp.DefAs("x", L(33, 29, "abc"))
	_, err := gisp.Parse(`(var y '(33 45 "def"))`)
	if err != nil {
		t.Fatalf("excpet define b but error: %v", err)
	}
	gisp.Parse(`(var z '(33 45 "def"))`)
	x, _ := gisp.Parse("x")
	y, _ := gisp.Parse("y")
	ret, err := gisp.Parse("(<= x y)")
	if err != nil {
		t.Fatalf("excpet not %v <= %v but error: %v", x, y, err)
	}
	if !ret.(bool) {
		t.Fatalf("excpet not %v <= %v", x, y)
	}
	z, _ := gisp.Parse("z")
	ret, err = gisp.Parse("(<= x z)")
	if err != nil {
		t.Fatalf("excpet not %v <= %v but error: %v", x, z, err)
	}
	if !ret.(bool) {
		t.Fatalf("excpet not %v <= %v", x, z)
	}
	ret, err = gisp.Parse("(<= y z)")
	if err != nil {
		t.Fatalf("excpet not %v <= %v but error: %v", y, z, err)
	}
	if !ret.(bool) {
		t.Fatalf("excpet not %v <= %v", y, z)
	}
}

func TestLessOptionList(t *testing.T) {
	gisp := NewGispWith(
		map[string]Toolbox{"axiom": Axiom, "props": Propositions},
		map[string]Toolbox{"time": Time})
	gisp.DefAs("x", L(33, nil, "abc"))
	_, err := gisp.Parse(`(var y '(33 45 "def"))`)
	if err != nil {
		t.Fatalf("excpet define b but error: %v", err)
	}
	gisp.Parse(`(var z '(33 45 "def"))`)
	x, _ := gisp.Parse("x")
	y, _ := gisp.Parse("y")
	ret, err := gisp.Parse("(<? x y)")
	if err != nil {
		t.Fatalf("excpet not %v <? %v but error: %v", x, y, err)
	}
	if !ret.(bool) {
		t.Fatalf("excpet not %v <? %v", x, y)
	}
	z, _ := gisp.Parse("z")
	ret, err = gisp.Parse("(<? x z)")
	if err != nil {
		t.Fatalf("excpet not %v <? %v but error: %v", x, z, err)
	}
	if !ret.(bool) {
		t.Fatalf("excpet not %v <= %v", x, z)
	}
	ret, err = gisp.Parse("(<? y z)")
	if err != nil {
		t.Fatalf("excpet not %v <? %v but error: %v", y, z, err)
	}
	if ret.(bool) {
		t.Fatalf("excpet not %v <? %v", y, z)
	}
}
