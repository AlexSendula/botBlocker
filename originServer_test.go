package main

import (
	"net/http"
	"testing"
)

func TestOriginServerResponding(t *testing.T) {
	go originServer()

	resp, err := http.Get("http://localhost:8081")
	expectedStatus := "200 OK"

	if err != nil {
		t.Errorf("originServer() FAILED. Failed to get response.\nError:\n%v\n", err)
	}

	if expectedStatus == resp.Status {
		t.Logf("HTTP GET request on originServer() PASSED. Expected response code: %s, got %s.\n", expectedStatus, resp.Status)
	} else {
		t.Errorf("originServer() FAILED. Expected response code: %s, got %s.\nFull respone:\n%v\n", expectedStatus, resp.Status, resp)
	}
}
