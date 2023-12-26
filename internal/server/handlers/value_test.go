package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smakimka/mtrcscollector/internal/logger"
	"github.com/smakimka/mtrcscollector/internal/model"
	mw "github.com/smakimka/mtrcscollector/internal/server/middleware"
	"github.com/smakimka/mtrcscollector/internal/storage"
)

func testValueRequest(t *testing.T, ts *httptest.Server, method,
	path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	if body != nil {
		req.Header.Add("Content-type", "application/json")
	}

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func getTestValueRouter() chi.Router {
	logger.SetLevel(logger.Debug)
	s := storage.NewMemStorage()

	s.UpdateGaugeMetric(model.GaugeMetric{Name: "test", Value: 1.5})

	getMetricValueHandler := GetMetricValueHandler{s}
	valueHandler := ValueHandler{s}

	r := chi.NewRouter()

	r.Post("/value", valueHandler.ServeHTTP)
	r.Route("/value/{metricKind}", func(r chi.Router) {
		r.Use(mw.MetricKind)
		r.Get("/{metricName}", getMetricValueHandler.ServeHTTP)
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
				body:        "no such metric",
			},
		},
	}

	ts := httptest.NewServer(getTestValueRouter())
	defer ts.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, body := testValueRequest(t, ts, "GET", test.url, nil)
			defer resp.Body.Close()

			assert.Equal(t, test.want.code, resp.StatusCode)
			assert.Equal(t, test.want.body, body)
			assert.Equal(t, test.want.contentType, resp.Header.Get("Content-Type"))
		})
	}
}

func TestValueHandler(t *testing.T) {
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
			body: testMetricsData{Name: "test", Kind: "gauge"},
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
				body:        "{\"id\":\"test\",\"type\":\"gauge\",\"value\":1.5}\n",
			},
		},
		{
			name: "negative test #1",
			body: testMetricsData{Name: "test", Kind: "counter"},
			want: want{
				code:        http.StatusNotFound,
				contentType: "application/json",
				body:        "{\"ok\":false,\"detail\":\"no such metric\"}\n",
			},
		},
	}

	ts := httptest.NewServer(getTestValueRouter())
	defer ts.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reqBody, err := json.Marshal(test.body)
			require.NoError(t, err)

			resp, body := testValueRequest(t, ts, "POST", "/value", bytes.NewReader(reqBody))
			defer resp.Body.Close()

			assert.Equal(t, test.want.code, resp.StatusCode)
			assert.Equal(t, test.want.body, body)
			assert.Equal(t, test.want.contentType, resp.Header.Get("Content-Type"))
		})
	}
}
