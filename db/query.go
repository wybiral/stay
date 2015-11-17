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
