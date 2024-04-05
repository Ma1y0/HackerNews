package main

import (
	"net/http"

	"github.com/Ma1y0/HackerNews/api"
)

func main() {
	r := http.NewServeMux()

	// Routes
	r.HandleFunc("GET /", api.GetStories)

	http.ListenAndServe(":8080", r)
}
