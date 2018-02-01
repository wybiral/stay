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

type Query struct {
	b *roaring.Bitmap
	n int
}

func Not(x *Query) *Query {
	return &Query {
		b: roaring.Flip(x.b, 0, uint64(x.n)),
		n: x.n,
	}
}

func And(x, y *Query) *Query {
	return &Query {
		b: roaring.And(x.b, y.b),
		n: x.n,
	}
}

func Or(x, y *Query) *Query {
	return &Query {
		b: roaring.Or(x.b, y.b),
		n: x.n,
	}
}

func Xor(x, y *Query) *Query {
	return &Query {
		b: roaring.Xor(x.b, y.b),
		n: x.n,
	}
}
