package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func reverseProxyServer() {
	originServerURL, err := url.Parse("http://localhost:8081")
	if err != nil {
		log.Fatal("invalid origin server URL")
	}

	reverseProxy := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// set req Host, URL and Request URI to forward a request to the origin server
		req.Host = originServerURL.Host
		req.URL.Host = originServerURL.Host
		req.URL.Scheme = originServerURL.Scheme
		req.RequestURI = ""

		//Add here code to filter out scrapers/bots

		// send a request to the origin server
		_, err := http.DefaultClient.Do(req)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(rw, err)
			return
		}
	})

	log.Fatal(http.ListenAndServe(":8080", reverseProxy))
}
