package main

import (
	"fmt"
	"net/http"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	cfg.fileserverHits.Add(1)
}

func (cfg *apiConfig) printMetrics(w http.responseWriter, req *http.Request) http.Handler {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	message := fmt.Sprintf("Hits: %d", cfg.fileserverHits)
	w.Write([]byte(message))
}

func main() {
	serve_mutex := http.NewServeMux()
	serve_mutex.Handle("/app/",
		middlewareMetricsInc(
			http.StripPrefix("/app", http.FileServer(http.Dir(".")))
		)
	)

	serve_mutex.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	serve_mutex.HandleFunc("/metrics", printMetrics)

	server := http.Server{
	Addr: ":8080",
	Handler: serve_mutex,
	}

	server.ListenAndServe()
}
