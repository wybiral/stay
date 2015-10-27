package db

import (
	"testing"
)

func TestBitvec(t *testing.T) {
	bv := NewBitvec()
	bv.Append(true)
	if bv.Len() != 1 {
		t.Errorf("Append didn't change Len")
	}
	if bv.Get(0) != true {
		t.Errorf("Append didn't set correct value")
	}
	bv.Set(0, false)
	if bv.Get(0) != false {
		t.Errorf("Set didn't change value")
	}
	bv.Append(true)
	if bv.Len() != 2 {
		t.Errorf("Second Append didn't change Len")
	}
	if bv.Get(0) != false || bv.Get(1) != true {
		t.Errorf("Second Append caused incorrect values")
	}
}
