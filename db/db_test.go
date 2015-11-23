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

func TestDatabase(t *testing.T) {
	db := NewDatabase()
	db.Add("nina", "species:human")
	db.Add("nina", "sex:female")
	db.Add("elaine", "species:cat")
	db.Add("elaine", "sex:female")
	db.Add("davy", "species:human")
	db.Add("davy", "sex:male")
	db.Add("percy", "species:cat")
	db.Add("percy", "sex:male")
	human := db.Query("species:human")
	cat := db.Query("species:cat")
	female := db.Query("sex:female")
	male := db.Query("sex:male")
	ch := db.Keys(human.And(male).Or(cat.And(female)))
	x := <-ch
	if x != "elaine" {
		t.Errorf("First result should be \"elaine\", got %s", x)
	}
	x = <-ch
	if x != "davy" {
		t.Errorf("Second result should be \"davy\", got %s", x)
	}
}
