package main

import (
	"fmt"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const port = "8080"
	const filepathRoot = "."

	apiCfg := &apiConfig{
		fileserverHits: 0,
	}

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir(filepathRoot))
	mux.Handle("/app/", http.StripPrefix("/app/", apiCfg.middlewareMetricsInc(fs)))
	mux.Handle("/healthz", http.HandlerFunc(healthzHandler))
	mux.Handle("/metrics", http.HandlerFunc(apiCfg.metricsHandler))
	mux.Handle("/reset", http.HandlerFunc(apiCfg.resetMetricsHandler))
	corsMux := middlewareCors(mux)
	server := &http.Server{
		Handler: corsMux,
		Addr:    "localhost:" + port,
	}
	fmt.Printf("Server listening on port %s\n", port)
	server.ListenAndServe()
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (c *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	hits := c.fileserverHits
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", hits)))
}

func (c *apiConfig) resetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	c.fileserverHits = 0
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
