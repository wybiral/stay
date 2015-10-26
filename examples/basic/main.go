package main

import (
	"fmt"
	"github.com/wybiral/stay"
)

func main() {
	idx := stay.NewIndex()

	idx.Set("user:0", "likes:a", true)
	idx.Set("user:0", "likes:b", true)
	idx.Set("user:0", "likes:c", true)

	idx.Set("user:1", "likes:a", true)
	idx.Set("user:1", "likes:b", true)
	idx.Set("user:1", "likes:d", true)

	idx.Set("user:2", "likes:a", true)
	idx.Set("user:2", "likes:d", true)

	idx.Set("user:3", "likes:d", true)

	// q1 = a & b
	q1 := idx.Query("likes:a").And(idx.Query("likes:b"))

	// q2 = not(a) & d
	q2 := idx.Query("likes:a").Not().And(idx.Query("likes:d"))

	// query = q1 | q2
	query := q1.Or(q2)

	for key := range idx.GetKeys(query) {
		fmt.Println(key)
	}
}
