package db

type Query chan chunk

func (a Query) Not() Query {
	ch := make(Query)
	go func() {
		for x := range a {
			ch <- ^x
		}
		close(ch)
	}()
	return ch
}

func (a Query) And(b Query) Query {
	ch := make(Query)
	go func() {
		for x := range a {
			y := <-b
			ch <- x & y
		}
		close(ch)
	}()
	return ch
}

func (a Query) Or(b Query) Query {
	ch := make(Query)
	go func() {
		for x := range a {
			y := <-b
			ch <- x | y
		}
		close(ch)
	}()
	return ch
}

func (a Query) Xor(b Query) Query {
	ch := make(Query)
	go func() {
		for x := range a {
			y := <-b
			ch <- x ^ y
		}
		close(ch)
	}()
	return ch
}

func (a Query) GetIds() chan int {
	ch := make(chan int)
	go func() {
		index := 0
		for x := range a {
			for j := 0; j < BITS; j++ {
				if x&chunk(1<<uint(j)) != 0 {
					ch <- index
				}
				index += 1
			}
		}
		close(ch)
	}()
	return ch
}

func (a Query) Count() int {
	count := 0
	for x := range a {
		for x > 0 {
			count++
			x &= (x - 1)
		}
	}
	return count
}
