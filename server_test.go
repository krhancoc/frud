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
	{"POST", "people", `{"name":"bob"}`, 201},
	{"POST", "people", `{"name":"tim"}`, 201},
	{"POST", "people", `{"name":"tim"}`, 409},
	{"PUT", "people", `{"name":"tim", "nickname": "Scott"}`, 201},
	{"POST", "people", `{"name":"testEntry", "nickname": "bigTuna", "supervisor":"bob", "partner":"tim" }`, 201},
	{"GET", "people?name=testEntry", ``, 200},
	{"POST", "meeting", `{"date":"March 21st", "attending": "testEntry"}`, 201},
	{"DELETE", "people", `{"name":"testEntry"}`, 200},
	{"DELETE", "meeting", `{"date":"March 21st"}`, 200},
	{"DELETE", "people", `{"name":"bob"}`, 200},
	{"DELETE", "people", `{"name":"tim"}`, 200},
}

func TestServerEndpoints(t *testing.T) {
	srv := StartServer("testResources/neo.json")

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
