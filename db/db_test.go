/*
Copyright 2015 Davy Wybiral <davy.wybiral@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
