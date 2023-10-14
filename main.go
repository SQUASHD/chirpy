package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
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

	r := chi.NewRouter()
	fs := apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))

	r.Handle("/app/*", http.StripPrefix("/app", fs))
	r.Handle("/app", http.StripPrefix("/app", fs))

	r.Get("/healthz", healthzHandler)
	r.Get("/metrics", apiCfg.metricsHandler)
	r.Handle("/reset", http.HandlerFunc(apiCfg.resetMetricsHandler))
	corsMux := middlewareCors(r)
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
