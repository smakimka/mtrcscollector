package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smakimka/mtrcscollector/internal/logger"
	mw "github.com/smakimka/mtrcscollector/internal/server/middleware"
	"github.com/smakimka/mtrcscollector/internal/storage"
)

func testUpdateRequest(t *testing.T, ts *httptest.Server, method,
	path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func getTestUpdateRouter() chi.Router {
	logger.SetLevel(logger.Debug)
	s := storage.NewMemStorage()

	updateMetricHandler := UpdateMetricHandler{s}
	updateHandler := UpdateHandler{s}

	r := chi.NewRouter()

	r.Post("/update", updateHandler.ServeHTTP)
	r.Route("/update/{metricKind}", func(r chi.Router) {
		r.Use(mw.MetricKind)
		r.Post("/{metricName}/{metricValue}", updateMetricHandler.ServeHTTP)
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
				code:        http.StatusBadRequest,
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
			resp, _ := testUpdateRequest(t, ts, "POST", test.url, nil)
			defer resp.Body.Close()

			assert.Equal(t, test.want.code, resp.StatusCode)
			assert.Equal(t, test.want.contentType, resp.Header.Get("Content-Type"))
		})
	}
}

type testMetricsData struct {
	Name  string  `json:"id"`
	Kind  string  `json:"type"`
	Value float64 `json:"value,omitempty"`
	Delta int64   `json:"delta,omitempty"`
}

func TestUpdateHandler(t *testing.T) {
	type want struct {
		code        int
		contentType string
		body        string
	}
	tests := []struct {
		name string
		body testMetricsData
		want want
	}{
		{
			name: "positive test #1",
			body: testMetricsData{Name: "test", Kind: "gauge", Value: 1.1},
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
				body:        "{\"id\":\"test\",\"type\":\"gauge\",\"value\":1.1}\n",
			},
		},
		{
			name: "positive test #2",
			body: testMetricsData{Name: "test", Kind: "counter", Delta: 1},
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
				body:        "{\"id\":\"test\",\"type\":\"counter\",\"value\":1}\n",
			},
		},
		{
			name: "not existent metric type",
			body: testMetricsData{Name: "test", Kind: "not_existing", Value: 1},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "application/json",
				body:        "{\"ok\":false,\"detail\":\"wrong metric kind\"}\n",
			},
		},
		{
			name: "no metric kind",
			body: testMetricsData{Name: "test", Value: 1},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "application/json",
				body:        "{\"ok\":false,\"detail\":\"missing some of required fields\"}\n",
			},
		},
		{
			name: "no metric value",
			body: testMetricsData{Name: "test", Kind: "gauge"},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "application/json",
				body:        "{\"ok\":false,\"detail\":\"missing some of required fields\"}\n",
			},
		},
		{
			name: "no metric name",
			body: testMetricsData{Kind: "gauge", Value: 1},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "application/json",
				body:        "{\"ok\":false,\"detail\":\"missing some of required fields\"}\n",
			},
		},
	}

	ts := httptest.NewServer(getTestUpdateRouter())
	defer ts.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := json.Marshal(test.body)
			require.NoError(t, err)

			fmt.Println(string(data))

			resp, body := testUpdateRequest(t, ts, "POST", "/update", bytes.NewReader(data))
			defer resp.Body.Close()

			assert.Equal(t, test.want.code, resp.StatusCode)
			assert.Equal(t, test.want.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, test.want.body, string(body))
		})
	}
}