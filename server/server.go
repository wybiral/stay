/*
Copyright 2015 Davy Wybiral <davy.wybiral@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"github.com/wybiral/bitvec"
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
func buildQuery(ctx *Context, x interface{}) *bitvec.Iterator {
	var query *bitvec.Iterator
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
		out <- []byte(fmt.Sprintf(`{"count":%d}`, query.Count()))
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
