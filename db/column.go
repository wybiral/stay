package db

type Column struct {
	count  int
	vector *Bitvec
}

func NewColumn() *Column {
	return &Column{count: 0, vector: NewBitvec()}
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

func (col *Column) Scan() Scan {
	return col.vector.Scan()
}
