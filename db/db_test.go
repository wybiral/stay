package db

import (
	"testing"
)

func TestColumn(t *testing.T) {
	col := NewColumn()
	col.Set(100000, true)
	if col.Get(100000) != true {
		t.Errorf("col.Set(100000, true) failed")
	}
	col.Set(100000, false)
	if col.Get(100000) != false {
		t.Errorf("col.Set(100000, false) failed")
	}
}
