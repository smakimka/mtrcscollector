package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testMetricRequest(t *testing.T, ts *httptest.Server, method,
	path string) *http.Response {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	return resp
}

func OkTestHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func getTestMetricRouter() chi.Router {
	r := chi.NewRouter()

	r.Route("/update/{metricKind}", func(r chi.Router) {
		r.Use(MetricKind)
		r.Post("/{metricName}/{metricValue}", OkTestHTTP)
	})

	return r
}

func TestMetricKind(t *testing.T) {
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
		{
			name: "negative test #1",
			url:  "/update/non-existent/test/1",
			want: 400,
		},
	}

	ts := httptest.NewServer(getTestMetricRouter())
	defer ts.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := testMetricRequest(t, ts, "POST", test.url)
			defer resp.Body.Close()

			assert.Equal(t, test.want, resp.StatusCode)
		})
	}
}
