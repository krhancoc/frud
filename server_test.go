package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

var endpoints = []struct {
	method         string
	endpoint       string
	data           string
	expectedStatus int
}{
	{"POST", "modelonly", `{"name":"testEntry"}`, 201},
	{"GET", "modelonly", `{"name":"testEntry"}`, 200},
	{"PUT", "modelonly", `{"name":"testEntry", "anotherField: "yeppers"}`, 200},
	{"DELETE", "modelonly", `{"name ":"testEntry"}`, 200},
}

func TestServerEndpoints(t *testing.T) {
	srv := StartServer("config.json")

	client := &http.Client{}
	for _, e := range endpoints {
		url := "http://localhost:8080/" + e.endpoint
		var jsonStr = []byte(e.data)
		req, err := http.NewRequest(e.method, url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			t.Error(err.Error())
		}
		defer resp.Body.Close()
		_, _ = ioutil.ReadAll(resp.Body)
		if resp.StatusCode != e.expectedStatus {
			t.Errorf("In %s, expected %d, but got %d", e.method, e.expectedStatus, resp.StatusCode)
		}
	}
	srv.Shutdown(nil)
}
