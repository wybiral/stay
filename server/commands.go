package server

import (
	"fmt"
	"encoding/json"
	"github.com/wybiral/stay/db"
)

type command interface {
	Execute(*stayContext)
}

type cmdUpdate struct {
	body    map[string][]string
	value bool
	results chan string
}

func (c *cmdUpdate) Execute(ctx *stayContext) {
	for key, columns := range c.body {
		for _, column := range columns {
			ctx.idx.Set(key, column, c.value)
		}
	}
	c.results <- ""
}

type cmdGet struct {
	body    []string
	results chan string
}

func (c *cmdGet) Execute(ctx *stayContext) {
	results := make(map[string][]string)
	for _, key := range c.body {
		results[key] = ctx.idx.GetColumns(key)
	}
	bytes, _ := json.Marshal(results)
	c.results <- string(bytes)
}

type cmdQuery struct {
	query   interface{}
	results chan string
}

func buildQuery(ctx *stayContext, x interface{}) db.Query {
	var query db.Query
	switch v := x.(type) {
	case string:
		query = ctx.idx.Query(v)
	case []interface{}:
		op := v[0].(string)
		query = buildQuery(ctx, v[1])
		if op == "not" {
			query = query.Not()
		} else {
			if op == "and" {
				for _, q := range v[2:] {
					query = query.And(buildQuery(ctx, q))
				}
			} else if op == "or" {
				for _, q := range v[2:] {
					query = query.Or(buildQuery(ctx, q))
				}
			} else if op == "xor" {
				for _, q := range v[2:] {
					query = query.Xor(buildQuery(ctx, q))
				}
			}
		}
	}
	return query
}

func (c *cmdQuery) Execute(ctx *stayContext) {
	query := buildQuery(ctx, c.query)
	results := make([]string, 0)
	for key := range ctx.idx.GetKeys(query) {
		results = append(results, key)
	}
	bytes, _ := json.Marshal(results)
	c.results <- string(bytes)
}

type cmdCount struct {
	query   interface{}
	results chan string
}

func (c *cmdCount) Execute(ctx *stayContext) {
	query := buildQuery(ctx, c.query)
	count := query.Count()
	c.results <- fmt.Sprintf(`{"count":%d}`, count)
}
