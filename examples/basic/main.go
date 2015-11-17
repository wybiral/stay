package main

import (
	"fmt"
	"github.com/wybiral/stay/db"
)

func main() {
	d := db.NewDatabase()

	d.Add("user:0", "likes:a")
	d.Add("user:0", "likes:b")
	d.Add("user:0", "likes:c")

	d.Add("user:1", "likes:a")
	d.Add("user:1", "likes:b")
	d.Add("user:1", "likes:d")

	d.Add("user:2", "likes:a")
	d.Add("user:2", "likes:d")

	d.Add("user:3", "likes:d")

	// q1 = a & b
	q1 := db.And(d.Query("likes:a"), d.Query("likes:b"))

	// q2 = not(a) & d
	q2 := db.And(db.Not(d.Query("likes:a")), d.Query("likes:d"))

	// query = q1 | q2
	query := db.Or(q1, q2)

	for key := range d.Keys(query) {
		fmt.Println(key)
	}
}
