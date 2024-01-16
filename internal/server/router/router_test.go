package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smakimka/mtrcscollector/internal/storage"
)

func testRequest(t *testing.T, ts *httptest.Server, url, method string) *http.Response {
	req, err := http.NewRequest(method, ts.URL+url, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	return resp
}

func TestRouter(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name   string
		method string
		url    string
		want   want
	}{
		{
			name:   "get all metrics",
			method: http.MethodGet,
			url:    "/",
			want: want{
				code:        http.StatusOK,
				contentType: "text/html; charset=utf-8",
			},
		},
		{
			name:   "update counter",
			method: http.MethodPost,
			url:    "/update/counter/test/1",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "update gauge",
			method: http.MethodPost,
			url:    "/update/gauge/test/1.1",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "update unknown",
			method: http.MethodPost,
			url:    "/update/unknown/test/1.1",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "update wrong value",
			method: http.MethodPost,
			url:    "/update/counter/test/1.1",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "get counter value",
			method: http.MethodGet,
			url:    "/value/counter/test",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "get gauge value",
			method: http.MethodGet,
			url:    "/value/gauge/test",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "get unknown value",
			method: http.MethodGet,
			url:    "/value/counter/unknown",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "get unknown value #2",
			method: http.MethodGet,
			url:    "/value/unknown/test",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	s := storage.NewMemStorage()
	ts := httptest.NewServer(GetRouter(s, &pgxpool.Pool{}))
	defer ts.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := testRequest(t, ts, test.url, test.method)
			defer resp.Body.Close()

			assert.Equal(t, test.want.code, resp.StatusCode)
			assert.Equal(t, test.want.contentType, resp.Header.Get("Content-Type"))
		})
	}
}
