package handlers

import (
	"log"
	"net/http"

	"github.com/smakimka/mtrcscollector/cmd/server/models"
	"github.com/smakimka/mtrcscollector/internal/mtrcs"
	"github.com/smakimka/mtrcscollector/internal/storage"
)

type MetricsUpdateHandler struct {
	Logger  *log.Logger
	Storage storage.Storage
}

func (h MetricsUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var data models.MetricsUpdateData

	w.Header().Add("Content-Type", "text/plain")

	code, err := data.ParsePath(r.URL.Path)
	if err != nil {
		w.WriteHeader(code)
		w.Write([]byte(err.Error()))
		h.Logger.Printf("%s    %d", r.RequestURI, code)
		return
	}

	var m mtrcs.Metric
	switch data.Kind {
	case mtrcs.Gauge:
		m = mtrcs.GaugeMetric{
			Name:  data.Name,
			Value: data.Value,
		}
	case mtrcs.Counter:
		m = mtrcs.CounterMetric{
			Name:  data.Name,
			Value: int64(data.Value),
		}
	}

	err = h.Storage.UpdateMetric(m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		h.Logger.Printf("%s    %d", r.RequestURI, http.StatusInternalServerError)
		return
	}

	h.Logger.Printf("%s    %d", r.RequestURI, http.StatusOK)
}
