package handlers

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	mw "github.com/smakimka/mtrcscollector/cmd/server/middleware"
	"github.com/smakimka/mtrcscollector/internal/mtrcs"
	"github.com/smakimka/mtrcscollector/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testValueRequest(t *testing.T, ts *httptest.Server, method,
	path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(body)
}

func getTestValueRouter() chi.Router {
	logger := log.New(os.Stdout, "", 3)

	s := &storage.MemStorage{Logger: logger}
	err := s.Init()
	if err != nil {
		log.Fatal(err)
	}
	s.UpdateMetric(mtrcs.GaugeMetric{Name: "test", Value: 1.5})
	storageMW := mw.WithMemStorage{S: s}

	r := chi.NewRouter()
	r.Use(storageMW.WithMemStorage)

	r.Route("/value/{metricKind}", func(r chi.Router) {
		r.Use(mw.MetricKind)
		r.Get("/{metricName}", GetMetricValueHandler)
	})

	return r
}

func TestGetMetricsHandler(t *testing.T) {
	type want struct {
		code        int
		contentType string
		body        string
	}
	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "positive test #1",
			url:  "/value/gauge/test",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
				body:        "1.5",
			},
		},
		{
			name: "negative test #1",
			url:  "/value/counter/test",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
				body:        "no such metric\n",
			},
		},
	}

	ts := httptest.NewServer(getTestValueRouter())
	defer ts.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, body := testValueRequest(t, ts, "GET", test.url)

			assert.Equal(t, test.want.code, resp.StatusCode)
			assert.Equal(t, test.want.body, body)
			assert.Equal(t, test.want.contentType, resp.Header.Get("Content-Type"))
		})
	}
}
