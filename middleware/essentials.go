package middleware

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/krhancoc/gotech/util"
)

// Converter will make sure the body always has a JSON object
func Converter(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "GET" {
			contentType := r.Header.Get("Content-Type")
			if contentType != "application/json" {

				r.ParseForm()
				converted := util.ConvertFlatMap(r.Form)
				jsonString, err := json.Marshal(converted)
				if err != nil {
					panic(err)
				}
				r.Body = ioutil.NopCloser(strings.NewReader(string(jsonString)))
			}
		}
		next.ServeHTTP(w, r)
	})
}
