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
	"bufio"
	"compress/flate"
	"encoding/binary"
	"os"
	"github.com/wybiral/bitvec"
)

type Database struct {
	ids     map[string]int
	keys    []string
	columns map[string]*Column
}

func (db *Database) Save(filename string) {
	rows := len(db.keys)
	columns := len(db.columns)
	f, _ := os.Create(filename)
	c, _ := flate.NewWriter(f, 1)
	w := bufio.NewWriter(c)
	binary.Write(w, binary.LittleEndian, uint32(rows))
	binary.Write(w, binary.LittleEndian, uint32(columns))
	for _, key := range db.keys {
		w.WriteString(key)
		w.WriteString("\x00")
	}
	for name, column := range db.columns {
		w.WriteString(name)
		w.WriteString("\x00")
		binary.Write(w, binary.LittleEndian, uint32(column.count))
		for id := range bitvec.Ids(column.Iterate()) {
			binary.Write(w, binary.LittleEndian, uint32(id))
		}
	}
	w.Flush()
	c.Close()
	f.Close()
}

func (db *Database) Load(filename string) {
	var rows, columns uint32
	f, _ := os.Open(filename)
	c := flate.NewReader(f)
	r := bufio.NewReader(c)
	binary.Read(r, binary.LittleEndian, &rows)
	binary.Read(r, binary.LittleEndian, &columns)
	ids := make(map[string]int)
	db.ids = ids
	keys := make([]string, rows)
	db.keys = keys
	for i := 0; i < int(rows); i++ {
		b, _ := r.ReadBytes('\x00')
		key := string(b[:len(b)-1])
		keys[i] = key
		ids[key] = i
	}
	db.columns = make(map[string]*Column)
	for i := 0; i < int(columns); i++ {
		var id, count uint32
		name, _ := r.ReadBytes('\x00')
		binary.Read(r, binary.LittleEndian, &count)
		column := db.getColumn(string(name[:len(name)-1]))
		for j := 0; j < int(count); j++ {
			binary.Read(r, binary.LittleEndian, &id)
			column.Set(int(id), true)
		}
	}
	c.Close()
	f.Close()
}

/*
Create and return a new empty database.
*/
func NewDatabase() *Database {
	ids := make(map[string]int)
	keys := make([]string, 0)
	columns := make(map[string]*Column)
	db := &Database{ids: ids,
		keys: keys, columns: columns}
	return db
}

func (db *Database) Len() int {
	return len(db.keys)
}

/*
Get the integer id of a key or create one if the key doesn't exist.
*/
func (db *Database) getId(key string) int {
	id, ok := db.ids[key]
	if !ok {
		id = len(db.keys)
		db.ids[key] = id
		db.keys = append(db.keys, key)
	}
	return id
}

/*
Get *Column object by name or create one if it doesn't exist.
*/
func (db *Database) getColumn(column string) *Column {
	col, ok := db.columns[column]
	if !ok {
		col = NewColumn()
		db.columns[column] = col
	}
	return col
}

/*
Add "column" to "key".
*/
func (db *Database) Add(key string, column string) {
	id := db.getId(key)
	col := db.getColumn(column)
	col.Set(id, true)
}

/*
Remove "column" from "key".
*/
func (db *Database) Remove(key string, column string) {
	id := db.getId(key)
	col := db.getColumn(column)
	col.Set(id, false)
}

type emptyScan struct{}

func (s *emptyScan) Next() (bitvec.Word, bool) {
	return 0, false
}

/*
Return a Query over an entire column.
*/
func (db *Database) Query(column string) bitvec.Iterator {
	col, ok := db.columns[column]
	if ok {
		return col.Iterate()
	}
	return &emptyScan{}
}

/*
Return resulting keys from a scan.
*/
func (db *Database) Keys(s bitvec.Iterator) chan string {
	ch := make(chan string)
	go func() {
		for id := range bitvec.Ids(s) {
			ch <- db.keys[id]
		}
		close(ch)
	}()
	return ch
}

/*
Return all columns as a mapping of [name]->[# of occurrences]
*/
func (db *Database) AllColumns() map[string]int {
	out := make(map[string]int)
	for column, col := range db.columns {
		out[column] = col.count
	}
	return out
}

