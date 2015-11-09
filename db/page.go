package db

type Page struct {
	chunks []chunk
}

func NewPage() *Page {
	return &Page{make([]chunk, PAGE_SIZE)}
}

func (page *Page) Set(i int, x bool) bool {
	n := i / BITS
	mask := chunk(1 << (uint(i) % BITS))
	c := page.chunks[n]
	if x {
		page.chunks[n] = c | mask
		return c&mask == 0
	} else {
		page.chunks[n] = c & ^mask
		return c&mask != 0
	}
}

func (page *Page) Get(i int) bool {
	n := i / BITS
	mask := chunk(1 << (uint(i) % BITS))
	return (page.chunks[n] & mask) > 0
}
