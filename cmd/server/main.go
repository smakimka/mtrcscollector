package main

import (
	"fmt"
	"net/http"
)

func main() {
	cfg := parseFlags()

	if err := run(cfg); err != nil {
		panic(err)
	}
}

func run(cfg *Config) error {
	fmt.Println("Running server on", cfg.Addr)
	return http.ListenAndServe(cfg.Addr, GetRouter())
}
