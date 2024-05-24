package middleware

import (
	"bytes"
	"crypto/rsa"
	"io"
	"net/http"
	"strings"

	"github.com/smakimka/mtrcscollector/internal/logger"
)

type DecryptMiddleware struct {
	PrivateKey *rsa.PrivateKey
}

func NewDecryptMiddleware(privateKey *rsa.PrivateKey) *DecryptMiddleware {
	return &DecryptMiddleware{PrivateKey: privateKey}
}

func (m *DecryptMiddleware) Decrypt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		encryption := r.Header.Get("Encryption")
		if encryption == "" || !strings.Contains(encryption, "crypto-key") {
			next.ServeHTTP(w, r)
			return
		}

		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			logger.Log.Error().Msg(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, err.Error())
			return
		}

		decryptedBody, err := rsa.DecryptPKCS1v15(nil, m.PrivateKey, body)
		if err != nil {
			logger.Log.Error().Msg(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, err.Error())
			return
		}

		dr := io.NopCloser(bytes.NewReader(decryptedBody))
		r.Body = dr
		defer dr.Close()

		next.ServeHTTP(w, r)

	})
}
