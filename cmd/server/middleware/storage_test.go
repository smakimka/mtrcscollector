package middleware

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/smakimka/mtrcscollector/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testStorageRequest(t *testing.T, ts *httptest.Server, method,
	path string) *http.Response {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	return resp
}

func getTestStoratgeRouter() chi.Router {
	logger := log.New(os.Stdout, "", 3)

	s := &storage.MemStorage{Logger: logger}
	err := s.Init()
	if err != nil {
		log.Fatal(err)
	}
	storageMW := WithMemStorage{S: s}

	r := chi.NewRouter()
	r.Use(storageMW.WithMemStorage)

	r.Post("/update/{metricKind}/{metricName}/{metricValue}", OkTestHTTP)

	return r
}

func TestStorage(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "positive test",
			url:  "/update/gauge/test/1",
			want: 200,
		},
	}

	ts := httptest.NewServer(getTestStoratgeRouter())
	defer ts.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := testStorageRequest(t, ts, "POST", test.url)
			defer resp.Body.Close()

			assert.Equal(t, test.want, resp.StatusCode)
		})
	}
}
