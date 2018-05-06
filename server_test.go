package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

var files = []struct {
	c        string
	database string
}{
	{"mongo.json", "mongo"},
	//{"neo.json", "neo4j"},
}

var endpoints = []struct {
	method         string
	endpoint       string
	data           string
	expectedStatus int
}{
	{"POST", "people", `{"name":"bob"}`, 201},
	{"POST", "people", `{"name":"tim"}`, 201},
	{"POST", "people", `{"name":"tim"}`, 400},
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
	for _, conf := range files {
		srv := StartServer("testResources/" + conf.c)
		client := &http.Client{}
		for i, e := range endpoints {
			url := "http://localhost:8080/" + e.endpoint
			var jsonStr = []byte(e.data)
			req, err := http.NewRequest(e.method, url, bytes.NewBuffer(jsonStr))
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			println(string(bodyBytes))
			if err != nil {
				t.Error(err.Error())
			}
			defer resp.Body.Close()
			_, _ = ioutil.ReadAll(resp.Body)
			if resp.StatusCode != e.expectedStatus {
				t.Errorf("%d: %s, In %s, expected %d, but got %d", i, conf.c, e.method, e.expectedStatus, resp.StatusCode)
			}
		}
		srv.Shutdown(nil)
	}
}
