package gisp

import (
	"reflect"
	"testing"
)

func TestZip(t *testing.T) {
	xlist := L("a", "b", "c")
	ylist := L(3.14, 1.414, 0)
	zlist := Zip(xlist, ylist)
	data := L(L("a", 3.14), L("b", 1.414), L("c", 0))
	if !reflect.DeepEqual(zlist, data) {
		t.Fatalf("excpet zip(%v,  %v) got %v but %v",
			xlist, ylist, data, zlist)
	}
}

func TestZipNil(t *testing.T) {
	xlist := L("a", "b", "c", "d")
	ylist := L(3.14, 1.414, 0)
	zlist := Zip(xlist, ylist)
	data := L(L("a", 3.14), L("b", 1.414), L("c", 0), L("d", nil))
	if !reflect.DeepEqual(zlist, data) {
		t.Fatalf("excpet zip(%v,  %v) got %v but %v",
			xlist, ylist, data, zlist)
	}
	zlist = Zip(ylist, xlist)
	data = L(L(3.14, "a"), L(1.414, "b"), L(0, "c"), L(nil, "d"))
	if !reflect.DeepEqual(zlist, data) {
		t.Fatalf("excpet zip(%v,  %v) got %v but %v",
			ylist, xlist, data, zlist)
	}
}
