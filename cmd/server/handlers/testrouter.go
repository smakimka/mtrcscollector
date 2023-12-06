package handlers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	chiMW "github.com/go-chi/chi/v5/middleware"
	mw "github.com/smakimka/mtrcscollector/cmd/server/middleware"
	"github.com/smakimka/mtrcscollector/internal/storage"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method,
	path string) *http.Response {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	return resp
}

func GetTestRouter() chi.Router {
	logger := log.New(os.Stdout, "", 3)

	s := &storage.MemStorage{Logger: logger}
	err := s.Init()
	if err != nil {
		log.Fatal(err)
	}
	storageMW := mw.WithMemStorage{S: s}

	r := chi.NewRouter()
	r.Use(chiMW.Logger)
	r.Use(storageMW.WithMemStorage)

	r.Route("/update/{metricKind}", func(r chi.Router) {
		r.Use(mw.MetricKind)
		r.Post("/{metricName}/{metricValue}", UpdateMetricHandler)
	})

	return r
}
