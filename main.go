package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

type options struct {
	port       int
	foldername string
}

func main() {
	var opts options

	flag.IntVar(&opts.port, "port", 8080, "Port to bind the web server on")
	flag.StringVar(&opts.foldername, "foldername", ".", "Path of the folder to expose")

	flag.Parse()

	if err := startServer(opts); err != nil {
		log.Fatal(err)
	}
}

func startServer(opts options) error {
	handler := makeHandler(opts.foldername)
	addr := fmt.Sprintf(":%d", opts.port)

	log.Printf("serving %q on %q", opts.foldername, addr)
	return http.ListenAndServe(addr, handler)
}
