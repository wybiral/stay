package main

import (
	"fmt"
	"github.com/wybiral/stay/db"
)

func main() {
	db := db.NewDatabase()

	db.Add("user:0", "likes:a")
	db.Add("user:0", "likes:b")
	db.Add("user:0", "likes:c")

	db.Add("user:1", "likes:a")
	db.Add("user:1", "likes:b")
	db.Add("user:1", "likes:d")

	db.Add("user:2", "likes:a")
	db.Add("user:2", "likes:d")

	db.Add("user:3", "likes:d")

	// q1 = a & b
	q1 := db.Query("likes:a").And(db.Query("likes:b"))

	// q2 = not(a) & d
	q2 := db.Query("likes:a").Not().And(db.Query("likes:d"))

	// query = q1 | q2
	query := q1.Or(q2)

	for key := range db.Keys(query) {
		fmt.Println(key)
	}
}
