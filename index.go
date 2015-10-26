/*
The Index struct holds all of the column vectors by name and does the job of
allocating new rows and returning their id.
*/

package stay

type Index struct {
	ids     map[string]int
	keys    []string
	columns map[string]*Bitvec
}

func NewIndex() *Index {
	return &Index{ids: make(map[string]int), keys: make([]string, 0), columns: make(map[string]*Bitvec)}
}

func (index *Index) newRow(key string) int {
	id := len(index.keys)
	index.keys = append(index.keys, key)
	index.ids[key] = id
	for _, column := range index.columns {
		column.Append(false)
	}
	return id
}

func (index *Index) newColumn(key string) *Bitvec {
	length := len(index.keys)
	chunks := make([]chunk, (length/64)+1)
	bitvec := &Bitvec{length: length, chunks: chunks}
	index.columns[key] = bitvec
	return bitvec
}

func (index *Index) Get(key string, column string) bool {
	id, ok := index.ids[key]
	if ok {
		return index.columns[column].Get(id)
	} else {
		return false
	}
}

func (index *Index) Set(key string, column string, value bool) {
	id, ok := index.ids[key]
	if !ok {
		id = index.newRow(key)
	}
	col, ok := index.columns[column]
	if !ok {
		col = index.newColumn(column)
	}
	col.Set(id, value)
}

func (index *Index) Query(column string) Query {
	return index.columns[column].Query()
}

func (index *Index) GetKeys(query Query) chan string {
	ch := make(chan string)
	go func() {
		for id := range query.GetIds() {
			ch <- index.keys[id]
		}
		close(ch)
	}()
	return ch
}
