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
	"github.com/wybiral/bitvec"
)

type Column struct {
	count  int
	vector *bitvec.Bitvec
}

func NewColumn() *Column {
	return &Column{count: 0, vector: bitvec.NewBitvec()}
}

func (col *Column) Set(id int, value bool) {
	if col.vector.Set(id, value) {
		if value {
			col.count += 1
		} else {
			col.count -= 1
		}
	}
}

func (col *Column) Get(id int) bool {
	return col.vector.Get(id)
}

func (col *Column) Iterate() bitvec.Iterator {
	return col.vector.Iterate()
}
