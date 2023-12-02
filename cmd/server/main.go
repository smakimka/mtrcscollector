package main

import (
	"log"
	"net/http"
	"os"

	"github.com/smakimka/mtrcscollector/cmd/server/handlers"
	"github.com/smakimka/mtrcscollector/cmd/server/middleware"
	"github.com/smakimka/mtrcscollector/internal/storage"
)

func main() {
	logger := log.New(os.Stdout, "", 5)

	s := &storage.MemStorage{Logger: logger}
	err := s.Init()
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()

	mux.Handle(`/update/`,
		middleware.MethodPOST(
			http.StripPrefix("/update/", handlers.MetricsUpdateHandler{
				Storage: s,
				Logger:  logger,
			})))

	log.Fatal(http.ListenAndServe(`localhost:8080`, mux))
}
