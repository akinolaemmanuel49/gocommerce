package errors

import "net/http"

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
