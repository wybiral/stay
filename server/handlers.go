package server

import (
	"encoding/json"
)

func errorMsg(msg string) string {
	obj := make(map[string]string)
	obj["error"] = msg
	bytes, _ := json.Marshal(obj)
	return string(bytes)
}

func handleAdd(ctx *stayContext, body []byte, out chan string) {
	var value map[string][]string
	err := json.Unmarshal(body, &value)
	if err != nil {
		out <- errorMsg("Malformed request body")
	} else {
		ctx.commands <- &cmdUpdate{value, true, out}
	}
}

func handleRemove(ctx *stayContext, body []byte, out chan string) {
	var value map[string][]string
	err := json.Unmarshal(body, &value)
	if err != nil {
		out <- errorMsg("Malformed request body")
	} else {
		ctx.commands <- &cmdUpdate{value, false, out}
	}
}

func handleGet(ctx *stayContext, body []byte, out chan string) {
	var value []string
	err := json.Unmarshal(body, &value)
	if err != nil {
		out <- errorMsg("Malformed request body")
	} else {
		ctx.commands <- &cmdGet{value, out}
	}
}

func handleQuery(ctx *stayContext, body []byte, out chan string) {
	var value interface{}
	err := json.Unmarshal(body, &value)
	if err != nil {
		out <- errorMsg("Malformed request body")
	} else {
		ctx.commands <- &cmdQuery{value, out}
	}
}

func handleCount(ctx *stayContext, body []byte, out chan string) {
	var value interface{}
	err := json.Unmarshal(body, &value)
	if err != nil {
		out <- errorMsg("Malformed request body")
	} else {
		ctx.commands <- &cmdCount{value, out}
	}
}

func handleColumns(ctx *stayContext, body []byte, out chan string) {
	columns := ctx.idx.Columns()
	bytes, _ := json.Marshal(columns)
	out <- string(bytes)
}

func handleTest(ctx *stayContext, body []byte, out chan string) {
	out <- "Hello world!"
}
