package gisp

import (
	"testing"
)

func TestAdds0(t *testing.T) {
	var data = []interface{}{0, 1, 2, 3, 4, 5, 6}
	s, err := adds(data...)
	if err != nil {
		t.Fatalf("except error is nil but %v", err)
	}
	if s.(int) != 21 {
		t.Fatalf("except sum 0~6 is 21 but got %v", s)
	}
}

func TestAdds1(t *testing.T) {
	var data = []interface{}{0, 1, 2, 3.14, 4, 5, 6}
	s, err := adds(data...)
	if err != nil {
		t.Fatalf("except error is nil but %v", err)
	}
	if s.(float64) != 21.14 {
		t.Fatalf("except sum 0, 1, 2, 3.14, 4, 5, 6 is 21.14 but got %v", s)
	}
}
