package handlers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/smakimka/mtrcscollector/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
				contentType: "text/plain",
			},
		},
		{
			name: "positive test #2",
			url:  "/update/counter/test/1",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain",
			},
		},
		{
			name: "wrong metric value",
			url:  "/update/counter/test/1.5",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
		{
			name: "not existent metric type",
			url:  "/update/not_existing/test/1",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
		{
			name: "two params",
			url:  "/update/gauge/test",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain",
			},
		},
		{
			name: "one param",
			url:  "/update/counter",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain",
			},
		},
	}

	logger := log.New(os.Stdout, "", 5)
	s := &storage.MemStorage{Logger: logger}
	err := s.Init()
	require.NoError(t, err, "error initializing memstorage")
	handler := http.StripPrefix("/update/", MetricsUpdateHandler{Logger: logger, Storage: s})

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.url, nil)

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
