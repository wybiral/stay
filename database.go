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
	"github.com/RoaringBitmap/roaring"
)

type Database struct {
	ids     map[string]uint32
	keys    []string
	columns map[string]*roaring.Bitmap
}

// Create and return a new empty database.
func NewDatabase() *Database {
	return &Database{
		ids: make(map[string]uint32),
		keys: make([]string, 0),
		columns: make(map[string]*roaring.Bitmap),
	}
}

func (db *Database) Len() int {
	return len(db.keys)
}

// Get the integer id of a key or create one if the key doesn't exist.
func (db *Database) getId(key string) uint32 {
	id, ok := db.ids[key]
	if !ok {
		id = uint32(len(db.keys))
		db.ids[key] = id
		db.keys = append(db.keys, key)
	}
	return id
}

// Get *Column object by name or create one if it doesn't exist.
func (db *Database) getColumn(column string) *roaring.Bitmap {
	col, ok := db.columns[column]
	if !ok {
		col = roaring.New()
		db.columns[column] = col
	}
	return col
}

// Add "column" to "key".
func (db *Database) Add(key string, column string) {
	id := db.getId(key)
	col := db.getColumn(column)
	col.Add(id)
}

// Remove "column" from "key".
func (db *Database) Remove(key string, column string) {
	id := db.getId(key)
	col := db.getColumn(column)
	col.Remove(id)
}

// Return a Query over an entire column.
func (db *Database) Query(column string) *Query {
	col, ok := db.columns[column]
	if !ok {
		return &Query{b: roaring.New(), n: len(db.keys)}
	}
	return &Query{b: col, n: len(db.keys)}
}

// Return resulting keys from a query.
func (db *Database) GetKeys(q *Query) chan string {
	ch := make(chan string)
	go func() {
		itr := q.b.Iterator()
		for itr.HasNext() {
			ch <- db.keys[itr.Next()]
		}
		close(ch)
	}()
	return ch
}

// Return array of columns.
func (db *Database) Columns() []string {
	out := make([]string, 0)
	for column, _ := range db.columns {
		out = append(out, column)
	}
	return out
}
