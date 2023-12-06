package handlers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	mw "github.com/smakimka/mtrcscollector/cmd/server/middleware"
	"github.com/smakimka/mtrcscollector/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testUpdateRequest(t *testing.T, ts *httptest.Server, method,
	path string) *http.Response {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	return resp
}

func getTestUpdateRouter() chi.Router {
	logger := log.New(os.Stdout, "", 3)

	s := &storage.MemStorage{Logger: logger}
	err := s.Init()
	if err != nil {
		log.Fatal(err)
	}
	storageMW := mw.WithMemStorage{S: s}

	r := chi.NewRouter()
	r.Use(storageMW.WithMemStorage)

	r.Route("/update/{metricKind}", func(r chi.Router) {
		r.Use(mw.MetricKind)
		r.Post("/{metricName}/{metricValue}", UpdateMetricHandler)
	})

	return r
}

func TestMetricsUpdateHandler(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "positive test #1",
			url:  "/update/gauge/test/1.0",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "positive test #2",
			url:  "/update/counter/test/1",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "wrong metric value",
			url:  "/update/counter/test/1.5",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "not existent metric type",
			url:  "/update/not_existing/test/1",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "two params",
			url:  "/update/gauge/test",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "one param",
			url:  "/update/counter",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	ts := httptest.NewServer(getTestUpdateRouter())
	defer ts.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := testUpdateRequest(t, ts, "POST", test.url)

			assert.Equal(t, test.want.code, resp.StatusCode)
			assert.Equal(t, test.want.contentType, resp.Header.Get("Content-Type"))
		})
	}
}
