package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testGzipRequest(t *testing.T, ts *httptest.Server, encoding,
	path string) *http.Response {

	body := bytes.NewBufferString("{test: 123}")
	if encoding != "" {
		bodyString := "{test: 123}"
		compressedString := bytes.NewBuffer([]byte{})
		zw := gzip.NewWriter(compressedString)

		zw.Write([]byte(bodyString))
		zw.Close()

		body = compressedString
	}

	req, err := http.NewRequest("POST", ts.URL+path, body)
	require.NoError(t, err)
	if encoding != "" {
		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Accept-Encoding", encoding)
	}

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	return resp
}

func MirrorTestHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(body)
	w.WriteHeader(http.StatusOK)
}

func getTestGzipRouter() chi.Router {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Use(Gzip)
		r.Post("/mirror", MirrorTestHTTP)
	})

	return r
}

func TestGzip(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		acceptEncoding string
		wantEncoding   string
	}{
		{
			name:           "gzip encoding",
			url:            "/mirror",
			acceptEncoding: "gzip",
			wantEncoding:   "gzip",
		},
		{
			name:           "no encoding",
			url:            "/mirror",
			acceptEncoding: "",
			wantEncoding:   "",
		},
	}

	ts := httptest.NewServer(getTestGzipRouter())
	defer ts.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := testGzipRequest(t, ts, test.acceptEncoding, test.url)
			defer resp.Body.Close()

			assert.Equal(t, test.wantEncoding, resp.Header.Get("Content-Encoding"))
		})
	}
}
