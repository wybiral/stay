package stay

type Query chan chunk

func Combine(a, b Query, fn func(chunk, chunk) chunk) Query {
	ch := make(Query)
	go func() {
		for x := range a {
			y := <-b
			ch <- fn(x, y)
		}
		close(ch)
	}()
	return ch
}

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
	return Combine(a, b, func(x, y chunk) chunk { return x & y })
}

func (a Query) Or(b Query) Query {
	return Combine(a, b, func(x, y chunk) chunk { return x | y })
}

func (a Query) Xor(b Query) Query {
	return Combine(a, b, func(x, y chunk) chunk { return x ^ y })
}

func (a Query) GetIds() chan int {
	ch := make(chan int)
	go func() {
		index := 0
		for x := range a {
			for j := 0; j < 64; j++ {
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
		for j := 0; j < 64; j++ {
			if x&chunk(1<<uint(j)) != 0 {
				count += 1
			}
		}
	}
	return count
}
