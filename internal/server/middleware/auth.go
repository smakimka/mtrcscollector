package middleware

import (
	"bytes"
	"encoding/hex"
	"hash"
	"io"
	"net/http"

	"github.com/smakimka/mtrcscollector/internal/auth"
)

type HashingResponseWriter struct {
	http.ResponseWriter
	hasher hash.Hash
}

func (w *HashingResponseWriter) Write(b []byte) (int, error) {
	size, err := w.hasher.Write(b)
	if err != nil {
		return size, err
	}
	size, err = w.ResponseWriter.Write(b)
	return size, err
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !auth.Enabled() {
			next.ServeHTTP(w, r)
			return
		}

		sign := r.Header.Get("HashSHA256")
		if sign == "" {
			next.ServeHTTP(w, r)
			return
		}

		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		decodedSign, err := hex.DecodeString(sign)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ok, err := auth.Check(decodedSign, body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		hashingWriter := &HashingResponseWriter{w, auth.GetHasher()}
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		next.ServeHTTP(hashingWriter, r)
		sign = string(hashingWriter.hasher.Sum(nil))
		w.Header().Add("HashSHA256", sign)
	})
}
