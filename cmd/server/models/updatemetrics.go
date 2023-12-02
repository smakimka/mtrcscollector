package models

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/smakimka/mtrcscollector/internal/mtrcs"
)

type MetricsUpdateData struct {
	Kind  int
	Name  string
	Value float64
}

func (d *MetricsUpdateData) ParsePath(path string) (int, error) {
	rawParams := strings.Split(path, "/")
	if len(rawParams) != 3 {
		return http.StatusNotFound, fmt.Errorf("expected 3 params, got %d", len(rawParams))
	}

	if rawParams[0] != "gauge" && rawParams[0] != "counter" {
		return http.StatusBadRequest, fmt.Errorf("expected metric type to be \"gauge\" or \"counter\", got \"%s\" instead", rawParams[0])
	}

	switch rawParams[0] {
	case "gauge":
		d.Kind = mtrcs.Gauge
	case "counter":
		d.Kind = mtrcs.Counter
	}
	d.Name = rawParams[1]

	var conv_err error
	switch d.Kind {
	case mtrcs.Counter:
		val, err := strconv.ParseInt(rawParams[2], 10, 64)
		conv_err = err
		d.Value = float64(val)
	case mtrcs.Gauge:
		val, err := strconv.ParseFloat(rawParams[2], 64)
		conv_err = err
		d.Value = val
	}

	if conv_err != nil {
		return http.StatusBadRequest, conv_err
	}

	return 200, nil
}
