package gisp

import (
	"testing"
)

func TestGinqSelect(t *testing.T) {
	data := QL(
		L(0, 1, 2, 3, 4, 5),
		L(1, 2, 3, 4, 5, 6),
		L(2, 3, 4, 5, 6, 7),
		L(3, 4, 5, 6, 7, 8))
	g := NewGispWith(
		map[string]Toolbox{"axiom": Axiom, "props": Propositions, "utils": Utils},
		map[string]Toolbox{"time": Time})
	g.DefAs("data", data)
	ginq, err := g.Parse(`
(ginq (select [1]))
`)
	if err != nil {
		t.Fatalf("except got a ginq query but error: %v ", err)
	}
	re, err := g.Eval(L(ginq, data))
	if err != nil {
		t.Fatalf("except got columns from data but error: %v", err)
	}

	t.Logf("ginq select got %v", re)
}

func TestGinqSelectFields(t *testing.T) {
	data := QL(
		L(0, 1, 2, 3, 4, 5),
		L(1, 2, 3, 4, 5, 6),
		L(2, 3, 4, 5, 6, 7),
		L(3, 4, 5, 6, 7, 8))
	g := NewGispWith(
		map[string]Toolbox{"axiom": Axiom, "props": Propositions, "utils": Utils},
		map[string]Toolbox{"time": Time})
	g.DefAs("data", data)
	ginq, err := g.Parse(`
(ginq (select (fs [1] [2] [4])))
`)
	if err != nil {
		t.Fatalf("except got a ginq query but error %v ", err)
	}
	re, err := g.Eval(L(ginq, data))
	if err != nil {
		t.Fatalf("except got columns from data but error %v", err)
	}

	t.Logf("ginq select got %v", re)
}

func TestGinqWhereSelect(t *testing.T) {
	data := QL(
		L(0, 1, 2, 3, 4, 5),
		L(1, 2, 3, 4, 5, 6),
		L(2, 3, 4, 5, 6, 7),
		L(3, 4, 5, 6, 7, 8))
	g := NewGispWith(
		map[string]Toolbox{"axiom": Axiom, "props": Propositions, "utils": Utils},
		map[string]Toolbox{"time": Time})
	g.DefAs("data", data)
	ginq, err := g.Parse(`
	(ginq
		(where (lambda (r) (< 1 r[0])))
		(select (fs [1] [2] [4]))
	)
	`)
	if err != nil {
		t.Fatalf("except got a ginq query but error %v ", err)
	}
	re, err := g.Eval(L(ginq, data))
	if err != nil {
		t.Fatalf("except got columns from data but error %v", err)
	}

	t.Logf("ginq select got %v", re)
}

func TestGinqSum(t *testing.T) {
	data := QL(
		L(0, 1, 2, 3, 4, 5),
		L(1, 2, 3, 4, 5, 6),
		L(0, 1, 2, 3, 4, 5),
		L(1, 2, 3, 4, 5, 6),
		L(2, 3, 4, 5, 6, 7),
		L(1, 2, 3, 4, 5, 6),
		L(2, 3, 4, 5, 6, 7),
		L(3, 4, 5, 6, 7, 8))
	g := NewGispWith(
		map[string]Toolbox{"axiom": Axiom, "props": Propositions, "utils": Utils},
		map[string]Toolbox{"time": Time})
	g.DefAs("data", data)
	ginq, err := g.Parse(`
	(ginq
		(select [5]) sum
	)
	`)
	if err != nil {
		t.Fatalf("except got a ginq query but error %v ", err)
	}
	re, err := g.Eval(L(ginq, data))
	if err != nil {
		t.Fatalf("except got columns from data but error: %v", err)
	}

	t.Logf("ginq select got %v", re)
}

func TestGinqGroupBy(t *testing.T) {
	data := QL(
		L(0, 1, 2, 3, 4, 5),
		L(1, 2, 3, 4, 5, 6),
		L(0, 1, 2, 3, 4, 5),
		L(1, 2, 3, 4, 5, 6),
		L(2, 3, 4, 5, 6, 7),
		L(1, 2, 3, 4, 5, 6),
		L(2, 3, 4, 5, 6, 7),
		L(3, 4, 5, 6, 7, 8))
	g := NewGispWith(
		map[string]Toolbox{"axiom": Axiom, "props": Propositions, "utils": Utils},
		map[string]Toolbox{"time": Time})
	g.DefAs("data", data)
	ginq, err := g.Parse(`
	(ginq
		(groupby [0] (ginq (select [5]) sum))
	)
	`)
	if err != nil {
		t.Fatalf("except got a ginq query but error %v ", err)
	}
	re, err := g.Eval(L(ginq, data))
	if err != nil {
		t.Fatalf("except got columns from data but error: %v", err)
	}

	t.Logf("ginq select got %v", re)
}

func TestGinqGroupBySumSelect(t *testing.T) {
	data := QL(
		L(0, 1, 2, 3, 4, 5),
		L(1, 2, 3, 4, 5, 6),
		L(0, 1, 2, 3, 4, 5),
		L(1, 2, 3, 4, 5, 6),
		L(2, 3, 4, 5, 6, 7),
		L(1, 2, 3, 4, 5, 6),
		L(2, 3, 4, 5, 6, 7),
		L(3, 4, 5, 6, 7, 8))
	g := NewGispWith(
		map[string]Toolbox{"axiom": Axiom, "props": Propositions, "utils": Utils},
		map[string]Toolbox{"time": Time})
	g.DefAs("data", data)
	ginq, err := g.Parse(`
	(ginq
		(groupby [0] (sums [5]))
	)
	`)
	if err != nil {
		t.Fatalf("except got a ginq query but error %v ", err)
	}
	re, err := g.Eval(L(ginq, data))
	if err != nil {
		t.Fatalf("except got columns from data but error: %v", err)
	}

	t.Logf("ginq select got %v", re)
}

func TestGinqGroupBySumSelectWhere(t *testing.T) {
	data := QL(
		L(0, 1, 2, 3, 4, 5),
		L(1, 2, 3, 4, 5, 6),
		L(0, 1, 2, 3, 4, 5),
		L(1, 2, 3, 4, 5, 6),
		L(2, 3, 4, 5, 6, 7),
		L(1, 2, 3, 4, 5, 6),
		L(2, 3, 4, 5, 6, 7),
		L(3, 4, 5, 6, 7, 8))
	g := NewGispWith(
		map[string]Toolbox{"axiom": Axiom, "props": Propositions, "utils": Utils},
		map[string]Toolbox{"time": Time})
	g.DefAs("data", data)
	ginq, err := g.Parse(`
	(ginq
		(groupby [0] (sums [5]))
		(where (lambda (x) (> 10 x[1])))
	)
	`)
	if err != nil {
		t.Fatalf("except got a ginq query but error %v ", err)
	}
	re, err := g.Eval(L(ginq, data))
	if err != nil {
		t.Fatalf("except got columns from data but error: %v", err)
	}

	t.Logf("ginq select got %v", re)
}
