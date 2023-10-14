package main

import (
	"fmt"
	"net/http"
)

func main() {
	const port = "8080"
	const filepathRoot = "."
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir(filepathRoot))
	mux.Handle("/app/", http.StripPrefix("/app/", fs))
	mux.Handle("/healthz", http.HandlerFunc(healthzHandler))

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
