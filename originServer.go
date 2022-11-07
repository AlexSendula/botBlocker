package main

import (
	"fmt"
	"log"
	"net/http"
)

// The origin server represents the target for the bots and the proxy.
func originServer() {
	originServerHandler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		_, _ = fmt.Fprint(rw, "origin originServer response")
	})

	log.Fatal(http.ListenAndServe(":8081", originServerHandler))
}
