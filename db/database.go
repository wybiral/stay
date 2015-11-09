package db

type chunk uint64

const (
	BITS      = 64               // Bits in a chunk
	PAGE_SIZE = 64               // Chunks in a page
	PAGE_BITS = BITS * PAGE_SIZE // Bits in a page
)

type Database struct {
	ids     map[string]int
	keys    []string
	columns map[string]*Column
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

/*
Check if "key" contains "column"
*/
func (db *Database) Get(key string, column string) bool {
	id, ok := db.ids[key]
	if !ok {
		return false
	}
	col, ok := db.columns[column]
	if !ok {
		return false
	}
	return col.Get(id)
}

/*
Return a Query over an entire column.
*/
func (db *Database) Query(column string) Query {
	length := (len(db.keys) / BITS) + 1
	col, ok := db.columns[column]
	if ok {
		return col.Query(length)
	}
	return emptyQuery(length)
}

/*
Return resulting keys from a query.
*/
func (db *Database) Keys(query Query) chan string {
	ch := make(chan string)
	go func() {
		id := 0
		for x := range query {
			mask := chunk(1)
			for i := 0; i < BITS; i++ {
				if x&mask != 0 {
					ch <- db.keys[id]
				}
				mask <<= 1
				id += 1
			}
		}
		close(ch)
	}()
	return ch
}

func emptyQuery(length int) Query {
	ch := make(chan chunk)
	go func() {
		for i := 0; i < length; i++ {
			ch <- chunk(0)
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
