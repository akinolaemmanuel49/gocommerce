package utils

import (
	"encoding/json"
	"net/http"

	l "github.com/akinolaemmanuel49/gocommerce/log"
)

func WriteJSON(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}) {
	logger := l.SetupLogger("service.log", "INFO")
	w.Header().Set("Content-Type", "application/json")
	if rw, ok := w.(*ErrorResponseWriter); ok && !rw.Written {
		rw.WriteHeader(statusCode)
	} else {
		w.WriteHeader(statusCode)
	}

	logger.Printf("%s %s", r.Method, r.URL.Path)
	json.NewEncoder(w).Encode(data)
}

// ErrorResponseWriter tracks errors during request handling
type ErrorResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	Written    bool
	Err        error
}

// WriteHeader intercepts the status code
func (er *ErrorResponseWriter) WriteHeader(code int) {
	if !er.Written {
		er.StatusCode = code
		er.ResponseWriter.WriteHeader(code)
		er.Written = true
	}

}

// Write intercepts the response body and detects errors
func (er *ErrorResponseWriter) Write(body []byte) (int, error) {
	n, err := er.ResponseWriter.Write(body)
	if err != nil {
		er.Err = err
	}
	return n, err
}
