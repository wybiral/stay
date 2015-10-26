package main

import (
	"encoding/json"
	"fmt"
	"github.com/wybiral/stay"
	"net/http"
)

type command interface {
	Execute(*stayContext)
}

type stayContext struct {
	idx      *stay.Index
	commands chan command
}

func (ctx *stayContext) listen() {
	go func() {
		for cmd := range ctx.commands {
			cmd.Execute(ctx)
		}
	}()
}

type stayHandler struct {
	ctx     *stayContext
	handler func(*stayContext, http.ResponseWriter, *http.Request)
}

func (h stayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	h.handler(h.ctx, w, r)
}

type cmdUpdate struct {
	body    map[string]map[string]bool
	results chan string
}

func (c cmdUpdate) Execute(ctx *stayContext) {
	for key, columns := range c.body {
		for column, value := range columns {
			ctx.idx.Set(key, column, value)
		}
	}
	c.results <- `{}`
}

func handleUpdate(ctx *stayContext, w http.ResponseWriter, r *http.Request) {
	var body map[string]map[string]bool
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		fmt.Println(err)
		return
	}
	results := make(chan string)
	ctx.commands <- cmdUpdate{body, results}
	_ = <-results
}

type cmdGet struct {
	body    []string
	results chan string
}

func (c cmdGet) Execute(ctx *stayContext) {
	results := make(map[string][]string)
	for _, key := range c.body {
		results[key] = ctx.idx.GetColumns(key)
	}
	bytes, _ := json.Marshal(results)
	c.results <- string(bytes)
}

func handleGet(ctx *stayContext, w http.ResponseWriter, r *http.Request) {
	var body []string
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		fmt.Println(err)
		return
	}
	results := make(chan string)
	ctx.commands <- cmdGet{body, results}
	out := <-results
	fmt.Println(out)
}

type cmdQuery struct {
	query   interface{}
	results chan string
}

func buildQuery(ctx *stayContext, x interface{}) stay.Query {
	var query stay.Query
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

func (c cmdQuery) Execute(ctx *stayContext) {
	query := buildQuery(ctx, c.query)
	results := make([]string, 0)
	for key := range ctx.idx.GetKeys(query) {
		results = append(results, key)
	}
	bytes, _ := json.Marshal(results)
	c.results <- string(bytes)
}

func handleQuery(ctx *stayContext, w http.ResponseWriter, r *http.Request) {
	var body interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		fmt.Println(err)
		return
	}
	results := make(chan string)
	ctx.commands <- cmdQuery{body, results}
	out := <-results
	fmt.Fprint(w, out)
}

type cmdCount struct {
	query   interface{}
	results chan string
}

func (c cmdCount) Execute(ctx *stayContext) {
	query := buildQuery(ctx, c.query)
	results := make(map[string]int)
	results["count"] = query.Count()
	bytes, _ := json.Marshal(results)
	c.results <- string(bytes)
}

func handleCount(ctx *stayContext, w http.ResponseWriter, r *http.Request) {
	var body interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		fmt.Println(err)
		return
	}
	results := make(chan string)
	ctx.commands <- cmdCount{body, results}
	out := <-results
	fmt.Fprint(w, out)
}

func main() {
	fmt.Println("Starting server...")
	ctx := &stayContext{stay.NewIndex(), make(chan command)}
	ctx.listen()
	http.Handle("/update", stayHandler{ctx, handleUpdate})
	http.Handle("/get", stayHandler{ctx, handleGet})
	http.Handle("/query", stayHandler{ctx, handleQuery})
	http.Handle("/count", stayHandler{ctx, handleCount})
	http.ListenAndServe(":8080", nil)
}
