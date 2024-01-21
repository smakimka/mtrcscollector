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
	"github.com/smakimka/mtrcscollector/internal/storage"
)

func testUpdatesRequest(t *testing.T, ts *httptest.Server, method,
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

func getTestUpdatesRouter() chi.Router {
	logger.SetLevel(logger.Debug)
	s := storage.NewMemStorage()

	updatesHandler := UpdatesHandler{s}

	r := chi.NewRouter()

	r.Post("/updates/", updatesHandler.ServeHTTP)

	return r
}

type testUpdateMetricsData struct {
	Name  string  `json:"ID"`
	Kind  string  `json:"MType"`
	Value float64 `json:"Value,omitempty"`
	Delta int64   `json:"Delta,omitempty"`
}

type testMetricsData []testUpdateMetricsData

func TestUpdatesHandler(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		body testMetricsData
		want want
	}{
		{
			name: "two different metrics",
			body: testMetricsData{
				testUpdateMetricsData{
					Name:  "testCounter",
					Kind:  "counter",
					Delta: 1,
				},
				testUpdateMetricsData{
					Name:  "testGauge",
					Kind:  "gauge",
					Value: 1,
				},
			},
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
			},
		},
		{
			name: "two counter metrics",
			body: testMetricsData{
				testUpdateMetricsData{
					Name:  "testCounter",
					Kind:  "counter",
					Delta: 1,
				},
				testUpdateMetricsData{
					Name:  "testCounter",
					Kind:  "counter",
					Delta: 1,
				},
			},
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
			},
		},
		{
			name: "wrong metric",
			body: testMetricsData{
				testUpdateMetricsData{
					Name:  "wrongType",
					Kind:  "test",
					Delta: 1,
					Value: 1,
				},
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "application/json",
			},
		},
	}

	ts := httptest.NewServer(getTestUpdatesRouter())
	defer ts.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := json.Marshal(test.body)
			require.NoError(t, err)

			fmt.Println(string(data))

			resp, _ := testUpdatesRequest(t, ts, "POST", "/updates/", bytes.NewReader(data))
			defer resp.Body.Close()

			assert.Equal(t, test.want.code, resp.StatusCode)
			assert.Equal(t, test.want.contentType, resp.Header.Get("Content-Type"))
		})
	}
}
