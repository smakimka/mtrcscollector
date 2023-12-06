package middleware

import (
	"context"
	"net/http"

	"github.com/smakimka/mtrcscollector/internal/storage"
)

type WithMemStorage struct {
	S storage.Storage
}

func (s WithMemStorage) WithMemStorage(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), StorageKey, s.S)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
