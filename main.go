package main

import (
	"fmt"
	"flag"
	"github.com/wybiral/stay/server"
)

func main() {
	host := flag.String("host", "localhost", "Server host")
	port := flag.Int("port", 8080, "Server port")
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", *host, *port)
	server.Start(addr)
}
