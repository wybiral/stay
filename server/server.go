package server

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"github.com/wybiral/stay/db"
)

func errorMsg(msg string) string {
	obj := make(map[string]string)
	obj["error"] = msg
	bytes, _ := json.Marshal(obj)
	return string(bytes)
}

type Context struct {
	db *db.Database
}

type stayHandler struct {
	ctx     *Context
	handler func(*Context, []byte, chan string)
}

func (h stayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	out := make(chan string)
	body, _ := ioutil.ReadAll(r.Body)
	go h.handler(h.ctx, body, out)
	results := <-out
	text := results
	io.WriteString(w, text)
}

func handleAdd(ctx *Context, body []byte, out chan string) {
	var value map[string][]string
	err := json.Unmarshal(body, &value)
	if err != nil {
		out <- errorMsg("Malformed request body")
	} else {
		for key, columns := range value {
			for _, column := range columns {
				ctx.db.Add(key, column)
			}
		}
		out <- "{}"
	}
}

func handleRemove(ctx *Context, body []byte, out chan string) {
	var value map[string][]string
	err := json.Unmarshal(body, &value)
	if err != nil {
		out <- errorMsg("Malformed request body")
	} else {
		for key, columns := range value {
			for _, column := range columns {
				ctx.db.Remove(key, column)
			}
		}
		out <- "{}"
	}
}
func buildQuery(ctx *Context, x interface{}) db.Query {
	var query db.Query
	switch v := x.(type) {
	case string:
		query = ctx.db.Query(v)
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

func handleQuery(ctx *Context, body []byte, out chan string) {
	var value interface{}
	err := json.Unmarshal(body, &value)
	if err != nil {
		out <- errorMsg("Malformed request body")
	} else {
		query := buildQuery(ctx, value)
		keys := ctx.db.Keys(query)
		obj := make([]string, 0)
		for key := range keys {
			obj = append(obj, key)
		}
		bytes, _ := json.Marshal(obj)
		out <- string(bytes)
	}
}

func handleCount(ctx *Context, body []byte, out chan string) {
	var value interface{}
	err := json.Unmarshal(body, &value)
	if err != nil {
		out <- errorMsg("Malformed request body")
	} else {
		query := buildQuery(ctx, value)
		out <- fmt.Sprintf(`{"count":%d}`, query.Count())
	}
}

func handleColumns(ctx *Context, body []byte, out chan string) {
	columns := ctx.db.AllColumns()
	bytes, _ := json.Marshal(columns)
	out <- string(bytes)
}

func Start(addr string) {
	fmt.Printf("Starting server at %s...\n", addr)
	ctx := &Context{db.NewDatabase()}
	http.Handle("/add", stayHandler{ctx, handleAdd})
	http.Handle("/remove", stayHandler{ctx, handleRemove})
	http.Handle("/query", stayHandler{ctx, handleQuery})
	http.Handle("/count", stayHandler{ctx, handleCount})
	http.Handle("/columns", stayHandler{ctx, handleColumns})
	http.ListenAndServe(addr, nil)
}
