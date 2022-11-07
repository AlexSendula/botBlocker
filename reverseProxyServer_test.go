package main

import (
	"net/http"
	"testing"
)

func TestReverseProxyServerResponding(t *testing.T) {
	go reverseProxyServer()

	expectedStatus := "200 OK"

	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		t.Errorf("reverseProxyServer() FAILED. Failed to get response.\nError:\n%v\n", err)
	}

	if expectedStatus == resp.Status {
		t.Logf("HTTP GET request on reverseProxyServer() PASSED. Expected response code: %s, got %s.\n", expectedStatus, resp.Status)
	} else {
		t.Errorf("reverseProxyServer() FAILED. Expected response code: %s, got %s.\nFull respone:\n%v\n", expectedStatus, resp.Status, resp)
	}
}
