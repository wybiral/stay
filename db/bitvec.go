/*
Bitvec is the core bit vector implementation for columns. Currently it's naive
and should be replaced with a compressed representation.
*/
package db

type chunk uint64
const BITS = 64

type Bitvec struct {
	length int
	chunks []chunk
}

func NewBitvec() *Bitvec {
	return &Bitvec{length: 0, chunks: make([]chunk, 0, 4)}
}

func (bv *Bitvec) Len() int {
	return bv.length
}

func (bv *Bitvec) Copy() *Bitvec {
	chunks := make([]chunk, len(bv.chunks), cap(bv.chunks))
	copy(chunks, bv.chunks)
	return &Bitvec{length: bv.length, chunks: chunks}
}

func (bv *Bitvec) Append(x bool) {
	i := bv.length
	bv.length += 1
	if bv.length > len(bv.chunks)*BITS {
		bv.chunks = append(bv.chunks, 0)
	}
	bv.Set(i, x)
}

func (bv *Bitvec) Set(i int, x bool) {
	n := i / BITS
	mask := chunk(1 << (uint(i) % BITS))
	if x {
		bv.chunks[n] |= mask
	} else {
		bv.chunks[n] &= ^mask
	}
}

func (bv *Bitvec) Get(i int) bool {
	n := i / BITS
	mask := chunk(1 << (uint(i) % BITS))
	return (bv.chunks[n] & mask) > 0
}

func (bv *Bitvec) Query() Query {
	ch := make(Query)
	go func() {
		for _, x := range bv.chunks {
			ch <- x
		}
		close(ch)
	}()
	return ch
}
