package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/smakimka/mtrcscollector/internal/logger"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

type gzipReader struct {
	io.ReadCloser
	zr *gzip.Reader
}

func (r gzipReader) Read(p []byte) (n int, err error) {
	return r.zr.Read(p)
}

func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		acceptEncoding := r.Header.Get("Accept-Encoding")
		contentEncoding := r.Header.Get("Content-Encoding")
		if contentEncoding == "" && acceptEncoding == "" {
			next.ServeHTTP(w, r)
			return
		}

		if strings.Contains(contentEncoding, "gzip") {

			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				logger.Log.Error().Msg(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				io.WriteString(w, err.Error())
				return
			}
			zr := gzipReader{r.Body, gz}

			r.Body = zr
			defer zr.Close()
		}

		if !strings.Contains(acceptEncoding, "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			logger.Log.Error().Msg(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")

		next.ServeHTTP(&gzipWriter{ResponseWriter: w, Writer: gz}, r)

	})
}
