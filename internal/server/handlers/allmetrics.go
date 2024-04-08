package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/render"

	"github.com/smakimka/mtrcscollector/internal/storage"
)

type GetAllMetricsHandler struct {
	s storage.Storage
}

func NewGetAllMetricsHandler(s storage.Storage) GetAllMetricsHandler {
	return GetAllMetricsHandler{s: s}
}

// GetAllMetrics godoc
// @Tags Get
// @Summary Запрос получения всех метрик
// @ID getAllMetrics
// @Accept  plain
// @Produce html
// @Success 200 {string} string ""
// @Failure 500 {string} string "Внутренняя ошибка"
// @Router / [get]
func (h GetAllMetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	gaugeMetrics, err := h.s.GetAllGaugeMetrics(ctx)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
		return
	}

	gaugeMetricsList := make([]string, len(gaugeMetrics))
	for i, mtrc := range gaugeMetrics {
		gaugeMetricsList[i] = fmt.Sprintf(`
		<tr>
			<th scope="row">%d</th>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
		</tr>`, i+1, mtrc.GetType(), mtrc.GetName(), mtrc.GetStringValue())
	}

	numHeadStart := len(gaugeMetrics)

	counterMetrics, err := h.s.GetAllCounterMetrics(ctx)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
		return
	}

	counterMetricsList := make([]string, len(counterMetrics))
	for i, mtrc := range counterMetrics {
		counterMetricsList[i] = fmt.Sprintf(`
		<tr>
			<th scope="row">%d</th>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
		</tr>`, numHeadStart+i+1, mtrc.GetType(), mtrc.GetName(), mtrc.GetStringValue())
	}

	html := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>All metrics</title>
		<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/5.0.0-alpha2/css/bootstrap.min.css">
		<link rel="icon" href="./favicon.ico" type="image/x-icon">
		<style>
			body{
				background-color: #000000;
			}
		</style>
	</head>
	<body>
		<main>
			<div class="container d-flex justify-content-center mt-5">
				<table class="table table-dark table-striped table-hover table-bordered">
					<thead class="thead-dark">
						<tr>
							<th scope="col">#</th>
							<th scope="col">Type</th>
							<th scope="col">Name</th>
							<th scope="col">Value</th>
						</tr>
					</thead>
					<tbody>
						%s
						%s
					</tbody>
				</table>  
			</div>
		</main>
		<script src="https://stackpath.bootstrapcdn.com/bootstrap/5.0.0-alpha2/js/bootstrap.min.js"></script>
	</body>
	</html>
	`, strings.Join(gaugeMetricsList, "\n"), strings.Join(counterMetricsList, "\n"))

	render.Status(r, http.StatusOK)
	render.HTML(w, r, html)
}
