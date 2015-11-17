package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"github.com/wybiral/stay/db"
)

func errorMsg(msg string) []byte {
	obj := make(map[string]string)
	obj["error"] = msg
	bytes, _ := json.Marshal(obj)
	return bytes
}

type updatable interface {
	update(ctx *Context) 
}

type Context struct {
	db *db.Database
	updates chan updatable
}

func (ctx *Context) startUpdateLoop() {
	go func() {
		for x := range ctx.updates {
			x.update(ctx)
		}
	}()
}

type stayHandler struct {
	ctx     *Context
	handler func(*Context, []byte, chan []byte)
}

func (h stayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	out := make(chan []byte)
	body, _ := ioutil.ReadAll(r.Body)
	go h.handler(h.ctx, body, out)
	w.Write(<-out)
}

type addRemoveUpdate struct {
	action string
	mapping map[string][]string
	out chan []byte
}

func (u *addRemoveUpdate) update(ctx *Context) {
	if u.action == "add" {
		for key, columns := range u.mapping {
			for _, column := range columns {
				ctx.db.Add(key, column)
			}
		}
	} else {
		for key, columns := range u.mapping {
			for _, column := range columns {
				ctx.db.Remove(key, column)
			}
		}
	}
	u.out <- []byte("{}")
}

func handleAdd(ctx *Context, body []byte, out chan []byte) {
	var mapping map[string][]string
	err := json.Unmarshal(body, &mapping)
	if err != nil {
		out <- errorMsg("Malformed request body")
	} else {
		u := &addRemoveUpdate{action: "add", mapping: mapping, out: out}
		ctx.updates <- u
	}
}

func handleRemove(ctx *Context, body []byte, out chan []byte) {
	var mapping map[string][]string
	err := json.Unmarshal(body, &mapping)
	if err != nil {
		out <- errorMsg("Malformed request body")
	} else {
		u := &addRemoveUpdate{action: "remove", mapping: mapping, out: out}
		ctx.updates <- u
	}
}
func buildQuery(ctx *Context, x interface{}) db.Scan {
	var query db.Scan
	switch v := x.(type) {
	case string:
		query = ctx.db.Query(v)
	case []interface{}:
		op := v[0].(string)
		query = buildQuery(ctx, v[1])
		if op == "not" {
			query = db.Not(query)
		} else {
			if op == "and" {
				for _, q := range v[2:] {
					query = db.And(query, buildQuery(ctx, q))
				}
			} else if op == "or" {
				for _, q := range v[2:] {
					query = db.Or(query, buildQuery(ctx, q))
				}
			} else if op == "xor" {
				for _, q := range v[2:] {
					query = db.Xor(query, buildQuery(ctx, q))
				}
			}
		}
	}
	return query
}

func handleQuery(ctx *Context, body []byte, out chan []byte) {
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
		out <- bytes
	}
}

func handleCount(ctx *Context, body []byte, out chan []byte) {
	var value interface{}
	err := json.Unmarshal(body, &value)
	if err != nil {
		out <- errorMsg("Malformed request body")
	} else {
		query := buildQuery(ctx, value)
		out <- []byte(fmt.Sprintf(`{"count":%d}`, db.Count(query)))
	}
}

func handleStats(ctx *Context, body []byte, out chan []byte) {
	results := make(map[string]interface{})
	results["rows"] = ctx.db.Len()
	results["columns"] = ctx.db.AllColumns()
	bytes, _ := json.Marshal(results)
	out <- bytes
}

func handleSave(ctx *Context, body []byte, out chan []byte) {
	ctx.db.Save("backup.txt")
	out <- []byte("{}")
}

func handleLoad(ctx *Context, body []byte, out chan []byte) {
	ctx.db.Load("backup.txt")
	out <- []byte("{}")
}

func Start(addr string) {
	fmt.Printf("Starting server at %s...\n", addr)
	ctx := &Context{db.NewDatabase(), make(chan updatable)}
	ctx.startUpdateLoop()
	http.Handle("/add", stayHandler{ctx, handleAdd})
	http.Handle("/remove", stayHandler{ctx, handleRemove})
	http.Handle("/query", stayHandler{ctx, handleQuery})
	http.Handle("/count", stayHandler{ctx, handleCount})
	http.Handle("/stats", stayHandler{ctx, handleStats})
	http.Handle("/save", stayHandler{ctx, handleSave})
	http.Handle("/load", stayHandler{ctx, handleLoad})
	http.ListenAndServe(addr, nil)
}
