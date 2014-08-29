package gisp

import (
	"testing"
)

func TestOptionSet(t *testing.T) {
	var slot = DefOption(INT)
	slot.Set(Int(10))
	t.Logf("create a new var %v slot as %v", slot.Type(), slot.Get())
}

func TestOptionSetNil(t *testing.T) {
	var slot = DefOption(INT)
	slot.Set(nil)
	t.Logf("create a new var %v slot as %v", slot.Type(), slot.Get())
}

func TestOptionGetNil(t *testing.T) {
	var slot = DefOption(INT)
	val := slot.Get()
	if val != nil {
		t.Fatalf("except nil but %v", val)
	}
}

func TestOptionGetNilType(t *testing.T) {
	var slot = DefOption(INT)
	typ := slot.Type()
	if typ != INT {
		t.Fatalf("except INT type but %v", typ)
	}
}

func TestStrictSet(t *testing.T) {
	var slot = DefStrict(INT)
	slot.Set(Int(10))
	t.Logf("create a new var %v slot as %v", slot.Type(), slot.Get())
}

func TestStrictSetNil(t *testing.T) {
	defer func() {
		if re := recover(); re == nil {
			t.Fatal("excpet a panic when set nil to a strict value not pointer")
		}
	}()
	var slot = DefStrict(INT)
	slot.Set(nil)
}

func TestStrictGetNil(t *testing.T) {
	var slot = DefStrict(INT)
	val := slot.Get()
	if val == nil {
		t.Fatal("except zero value when init but nil")
	}
}

func TestStrictGetNilType(t *testing.T) {
	var slot = DefStrict(INT)
	typ := slot.Type()
	if typ != INT {
		t.Fatalf("except INT type but %v", typ)
	}
}
