package db

type Column struct {
	count int
	pages map[int]*Page
}

func NewColumn() *Column {
	pages := make(map[int]*Page)
	return &Column{pages: pages}
}

func (col *Column) GetPage(pageId int) *Page {
	page, ok := col.pages[pageId]
	if !ok {
		page = NewPage()
		col.pages[pageId] = page
	}
	return page
}

func (col *Column) Set(id int, value bool) {
	chunkId := id / BITS
	pageId := chunkId / PAGE_SIZE
	page := col.GetPage(pageId)
	if page.Set(id-(pageId*BITS*PAGE_SIZE), value) {
		if value {
			col.count += 1
		} else {
			col.count -= 1
		}
	}
}

func (col *Column) Get(id int) bool {
	chunkId := id / BITS
	pageId := chunkId / PAGE_SIZE
	page, ok := col.pages[pageId]
	if ok {
		return page.Get(id - (pageId * BITS * PAGE_SIZE))
	}
	return false
}

func (col *Column) Query(length int) Query {
	ch := make(chan chunk)
	npages := (length / PAGE_SIZE) + 1
	go func() {
		for i := 0; i < npages; i++ {
			page, ok := col.pages[i]
			if ok {
				for j := 0; j < PAGE_SIZE; j++ {
					ch <- page.chunks[j]
				}
			} else {
				for j := 0; j < PAGE_SIZE; j++ {
					ch <- chunk(0)
				}
			}
		}
		close(ch)
	}()
	return ch
}
