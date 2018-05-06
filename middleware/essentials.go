// Package middleware will hold the middleware objects used by the frud framework.
// Future middleware will be the authentication middleware
package middleware

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)
func ConvertFlatMap(form url.Values) map[string]interface{} {
    final := make(map[string]interface{}, len(form))
    for formKey, formValue := range form {
	final[formKey] = formValue
    }
    return final
}

// Converter will make sure the body always has a JSON object
func Converter(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "GET" {
			contentType := r.Header.Get("Content-Type")
			if contentType != "application/json" {
				r.ParseForm()
				converted := ConvertFlatMap(r.Form)
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
