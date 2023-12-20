package handlers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smakimka/mtrcscollector/internal/storage"
)

func testAllMtrcsRequest(t *testing.T, ts *httptest.Server, method string) *http.Response {
	req, err := http.NewRequest(method, ts.URL, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	return resp
}

func getTestAllMtrcsRouter() chi.Router {
	l := log.New(os.Stdout, "", 3)

	s := storage.NewMemStorage().
		WithLogger(l)

	getAllMetricsHandler := GetAllMetricsHandler{s, l}
	r := chi.NewRouter()
	r.Get("/", getAllMetricsHandler.ServeHTTP)

	return r
}

func TestAllMtrcsHandler(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				code:        http.StatusOK,
				contentType: "text/html; charset=utf-8",
			},
		},
	}

	ts := httptest.NewServer(getTestAllMtrcsRouter())
	defer ts.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := testAllMtrcsRequest(t, ts, "GET")
			defer resp.Body.Close()

			assert.Equal(t, test.want.code, resp.StatusCode)
			assert.Equal(t, test.want.contentType, resp.Header.Get("Content-Type"))
		})
	}
}
