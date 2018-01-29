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

type Query bitvec.Iterator

func Not(x Query) Query {
	return bitvec.Not(x)
}

func And(x, y Query) Query {
	return bitvec.And(x, y)
}

func Or(x, y Query) Query {
	return bitvec.Or(x, y)
}

func Xor(x, y Query) Query {
	return bitvec.Xor(x, y)
}
