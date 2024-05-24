package middleware

import (
	"net/http"
	"time"

	"github.com/smakimka/mtrcscollector/internal/logger"
)

type (
	responseDataStruct struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseDataStruct
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		url := r.RequestURI
		method := r.Method

		responseData := responseDataStruct{}
		lw := loggingResponseWriter{w, &responseData}
		next.ServeHTTP(&lw, r)

		duration := time.Since(start)

		logger.Log.Info().
			Str("url", url).
			Str("method", method).
			Int("status", responseData.status).
			Dur("duration", duration).
			Int("size", responseData.size).
			Send()
	})

}
