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

type Scan interface {
	Next() (val word, ok bool)
}

type notScan struct {
	x Scan
}

func (s *notScan) Next() (word, bool) {
	val, ok := s.x.Next()
	return ^(val | FILL_BIT), ok
}

func Not(x Scan) Scan {
	return &notScan{x}
}

type binaryScan struct {
	x, y Scan
}

type andScan binaryScan
type orScan binaryScan
type xorScan binaryScan

func (s *andScan) Next() (word, bool) {
	x, xok := s.x.Next()
	y, yok := s.y.Next()
	return x & y, xok && yok
}

func And(x, y Scan) Scan {
	return &andScan{x, y}
}

func (s *orScan) Next() (word, bool) {
	x, xok := s.x.Next()
	y, yok := s.y.Next()
	return x | y, xok && yok
}

func Or(x, y Scan) Scan {
	return &orScan{x, y}
}

func (s *xorScan) Next() (word, bool) {
	x, xok := s.x.Next()
	y, yok := s.y.Next()
	return x ^ y, xok && yok
}

func Xor(x, y Scan) Scan {
	return &xorScan{x, y}
}

func Count(s Scan) int {
	count := 0
	for {
		x, ok := s.Next()
		if !ok {
			break
		}
		for x > 0 {
			count++
			x &= (x - 1)
		}
	}
	return count
}

func Ids(s Scan) chan int {
	ch := make(chan int)
	go func() {
		id := 0
		for {
			w, ok := s.Next()
			if !ok {
				break
			}
			for i := word(0); i < wordbits-1; i++ {
				if w&(1<<i) != 0 {
					ch <- id + int(i)
				}
			}
			id += wordbits - 1
		}
		close(ch)
	}()
	return ch
}
