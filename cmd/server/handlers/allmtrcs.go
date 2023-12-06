package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/smakimka/mtrcscollector/cmd/server/middleware"
	"github.com/smakimka/mtrcscollector/internal/storage"
)

func GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	s := r.Context().Value(middleware.StorageKey).(storage.Storage)
	metrics, err := s.GetAllMetrics()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
	}

	mtrcsList := make([]string, len(metrics))
	for i, mtrc := range metrics {
		mtrcsList[i] = fmt.Sprintf(`
		<tr>
			<th scope="row">%d</th>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
		</tr>`, i+1, mtrc.GetType(), mtrc.GetName(), mtrc.GetStringValue())
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
					</tbody>
				</table>  
			</div>
		</main>
		<script src="https://stackpath.bootstrapcdn.com/bootstrap/5.0.0-alpha2/js/bootstrap.min.js"></script>
	</body>
	</html>
	`, strings.Join(mtrcsList, "\n"))

	render.Status(r, http.StatusOK)
	render.HTML(w, r, html)
}
