// Copyright 2015 Davy Wybiral <davy.wybiral@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package staydb

import (
	"testing"
)

func TestDatabase(t *testing.T) {
	db := NewDatabase()
	db.Add("user:0", "likes:a")
	db.Add("user:0", "likes:b")
	db.Add("user:0", "likes:c")
	db.Add("user:1", "likes:a")
	db.Add("user:1", "likes:b")
	db.Add("user:1", "likes:d")
	db.Add("user:2", "likes:a")
	db.Add("user:2", "likes:d")
	db.Add("user:3", "likes:d")
	// q1 = a & b
	q1 := And(db.Query("likes:a"), db.Query("likes:b"))
	// q2 = not(a) & d
	q2 := And(Not(db.Query("likes:a")), db.Query("likes:d"))
	// query = q1 | q2
	query := Or(q1, q2)
	keys := db.GetKeys(query)
	x := <-keys
	if x != "user:0" {
		t.Errorf("First result should be \"user:0\", got %s", x)
	}
	x = <-keys
	if x != "user:1" {
		t.Errorf("Second result should be \"user:1\", got %s", x)
	}
	x = <-keys
	if x != "user:3" {
		t.Errorf("Third result should be \"user:3\", got %s", x)
	}
}
