package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var Nonce = newHash()
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:106.0) Gecko/20100101 Firefox/106.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:105.0) Gecko/20100101 Firefox/105.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64; rv:106.0) Gecko/20100101 Firefox/106.0",
	"Mozilla/5.0 (X11; Linux x86_64; rv:105.0) Gecko/20100101 Firefox/105.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
}

func reverseProxyServer() {
	originServerURL, err := url.Parse("http://localhost:8081")
	if err != nil {
		log.Fatal("invalid origin server URL")
	}

	reverseProxy := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		clientNonce := strings.Split(req.URL.String(), "nonce/")
		jsCheck := javascriptFilter(clientNonce, rw)
		if jsCheck == "denied" {
			http.Redirect(rw, req, "http://localhost:8081/denied", 403)
			closeBody(req)
		} else if jsCheck == "wait" {
			closeBody(req)
		}

		if !userAgentFilter(req) {
			closeBody(req)
			return
		}

		// set req Host, URL and Request URI to forward a request to the origin server
		req.Host = originServerURL.Host
		req.URL.Host = originServerURL.Host
		req.URL.Scheme = originServerURL.Scheme
		req.RequestURI = ""

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

func userAgentFilter(req *http.Request) bool {
	clientUserAgent := req.Header.Get("User-Agent")
	response := containsString(userAgents, clientUserAgent)
	return response
}

// Responses: ["denied", "allowed", "wait"]
func javascriptFilter(clientNonce []string, rw http.ResponseWriter) string {
	response := "denied"

	if len(clientNonce) > 1 {
		if clientNonce[1] == Nonce {
			response = "allowed"
		} else {
			response = "denied"
		}
	} else {
		rw.Header().Set("Content-Type", "text/html; charset=utf-8")
		Nonce = newHash()
		var script = `<script>window.location = "http://localhost:8080/nonce/` + Nonce + `";</script>`
		_, err := fmt.Fprint(rw, script)
		if err != nil {
			log.Println("JavascriptFilter() FAILED.")
		}
		response = "wait"
	}
	return response
}

func closeBody(req *http.Request) {
	err := req.Body.Close()
	if err != nil {
		return
	}
}

func containsString(s []string, str string) bool {
	for _, el := range s {
		if el == str {
			return true
		}
	}
	return false
}

func newHash() string {
	timestamp := time.Now().String()
	h := sha256.New()
	h.Write([]byte(timestamp))
	return hex.EncodeToString(h.Sum(nil))
}
