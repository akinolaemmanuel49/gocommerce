package middleware

import (
	"log"
	"net/http"

	"github.com/akinolaemmanuel49/gocommerce/common/errors"
	l "github.com/akinolaemmanuel49/gocommerce/log"
	"github.com/akinolaemmanuel49/gocommerce/utils"
)

// ErrorMiddleware handles errors and logs them appropriately
func ErrorMiddleware(next http.Handler) http.Handler {
	errorLogger := l.SetupLogger("service.log", "ERROR")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				errorLogger.Printf("Panic: %v", rec)
				errors.HandleError(w, r, errors.NewInternalServerError())
			}
		}()

		er := &utils.ErrorResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}
		next.ServeHTTP(er, r)

		if er.Err != nil {
			errors.HandleError(w, r, er.Err)
		}
	})
}

// logError logs the error with contextual information
func logError(r *http.Request, logger *log.Logger, err error) {
	logger.Printf("ERROR: %s %s [User-Agent: %s]: %v", r.Method, r.URL.Path, r.UserAgent(), err)
}
