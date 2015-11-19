package db

import (
	"testing"
	"math/rand"
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

func TestBitvecP01(t *testing.T) {
	b := NewBitvec()
	n := 1000000
	p := 0.01
	data := make([]int, 0)
	for i := 0; i < n; i++ {
		if rand.Float64() < p {
			data = append(data, i)
			b.Set(i, true)
		}
	}
	ids := Ids(b.Scan())
	for i, x := range data {
		y := <- ids
		if x != y {
			t.Errorf("Failed on the %dth value", i)
			return
		}
	}
}

func TestBitvecP50(t *testing.T) {
	b := NewBitvec()
	n := 1000000
	p := 0.50
	data := make([]int, 0)
	for i := 0; i < n; i++ {
		if rand.Float64() < p {
			data = append(data, i)
			b.Set(i, true)
		}
	}
	ids := Ids(b.Scan())
	for i, x := range data {
		y := <- ids
		if x != y {
			t.Errorf("Failed on the %dth value", i)
			return
		}
	}
}

func TestDatabase(t *testing.T) {

}
