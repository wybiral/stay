/*
The Index struct holds all of the column vectors by name and does the job of
allocating new rows and returning their id.
*/

package stay

type Index struct {
    nrows int
    columns map[string]*Bitvec
}

func NewIndex() *Index {
    return &Index{nrows: 0, columns: make(map[string]*Bitvec)}
}

func (index *Index) Row() uint {
    id := uint(index.nrows)
    index.nrows += 1
    for _, column := range index.columns {
        column.Append(false)
    }
    return id
}

func (index *Index) AddColumn(key string) {
    chunks := make([]chunk, (index.nrows / 64) + 1)
    bitvec := &Bitvec{length: index.nrows, chunks: chunks}
    index.columns[key] = bitvec
}

func (index *Index) Get(id uint, column string) bool {
    return index.columns[column].Get(id)
}

func (index *Index) Set(id uint, column string, value bool) {
    index.columns[column].Set(id, value)
}

func (index *Index) Query(column string) Query {
    return index.columns[column].Query()
}
