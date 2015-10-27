package server

import (
	"io"
	"fmt"
	"net/http"
	"io/ioutil"
	"github.com/wybiral/stay/db"
)

type stayContext struct {
	idx      *db.Index
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
	handler func(*stayContext, []byte, chan string)
}

func (h stayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	out := make(chan string)
	body, _ := ioutil.ReadAll(r.Body)
	go h.handler(h.ctx, body, out)
	results := <-out
	text := results
	io.WriteString(w, text)
}

func Start(addr string) {
	fmt.Printf("Starting server at %s...\n", addr)
	ctx := &stayContext{db.NewIndex(), make(chan command)}
	ctx.listen()
	http.Handle("/add", stayHandler{ctx, handleAdd})
	http.Handle("/remove", stayHandler{ctx, handleRemove})
	http.Handle("/get", stayHandler{ctx, handleGet})
	http.Handle("/query", stayHandler{ctx, handleQuery})
	http.Handle("/count", stayHandler{ctx, handleCount})
	http.Handle("/columns", stayHandler{ctx, handleColumns})
	http.ListenAndServe(addr, nil)
}
