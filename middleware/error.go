package middleware

import (
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/gocommerce/common/errors"
	l "github.com/akinolaemmanuel49/gocommerce/log"
)

// ErrorMiddleware handles errors and logs them appropriately
func ErrorMiddleware(next http.Handler) http.Handler {
	errorLogger, err := l.SetupLogger("service.log", "ERROR")
	if err != nil {
		log.Fatalf("Error setting up error logger: %v", err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wrap the response writer
		er := &errors.ErrorResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}

		// Defer to handle panics and log errors
		defer func() {
			if rec := recover(); rec != nil {
				errorLogger.Printf("Panic: %v", rec)
				errors.HandleError(w, r, errors.NewInternalServerError(), errorLogger)
			}
		}()

		next.ServeHTTP(er, r)
	})
}
