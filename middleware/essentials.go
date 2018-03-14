package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/krhancoc/gotech/util"
)

// Bytes wrapper around byte slice
type Bytes []byte

// Close closes a byte string which does nothing...
func (Bytes) Close() error {
	return nil
}

// Read just copies over to another buffer
func (w Bytes) Read(p []byte) (int, error) {
	copy(p, w)
	return len(w), nil
}

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
				println(string(jsonString))
				r.Body = Bytes(jsonString)
			}
		}
		next.ServeHTTP(w, r)
	})
}
