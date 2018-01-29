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
	"github.com/wybiral/bitvec"
)

type Database struct {
	ids     map[string]int
	keys    []string
	columns map[string]*Column
}

// Create and return a new empty database.
func NewDatabase() *Database {
	return &Database{
		ids: make(map[string]int),
		keys: make([]string, 0),
		columns: make(map[string]*Column),
	}
}

func (db *Database) Len() int {
	return len(db.keys)
}

// Get the integer id of a key or create one if the key doesn't exist.
func (db *Database) getId(key string) int {
	id, ok := db.ids[key]
	if !ok {
		id = len(db.keys)
		db.ids[key] = id
		db.keys = append(db.keys, key)
	}
	return id
}

// Get *Column object by name or create one if it doesn't exist.
func (db *Database) getColumn(column string) *Column {
	col, ok := db.columns[column]
	if !ok {
		col = NewColumn()
		db.columns[column] = col
	}
	return col
}

// Add "column" to "key".
func (db *Database) Add(key string, column string) {
	id := db.getId(key)
	col := db.getColumn(column)
	col.Set(id, true)
}

// Remove "column" from "key".
func (db *Database) Remove(key string, column string) {
	id := db.getId(key)
	col := db.getColumn(column)
	col.Set(id, false)
}

// Return a Query over an entire column.
func (db *Database) Query(column string) Query {
	query := bitvec.ZeroIterator(len(db.keys))
	col, ok := db.columns[column]
	if ok {
		query = bitvec.Or(query, col.Query())
	}
	return query
}

// Return resulting keys from a query.
func (db *Database) GetKeys(s Query) chan string {
	ch := make(chan string)
	go func() {
		for id := range bitvec.Indices(s) {
			ch <- db.keys[id]
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
