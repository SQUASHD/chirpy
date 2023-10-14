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

	router := chi.NewRouter()
	fs := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	router.Handle("/app/*", fs)
	router.Handle("/app", fs)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/metrics", apiCfg.handlerMetrics)
	apiRouter.Get("/reset", apiCfg.handlerReset)
	router.Mount("/api", apiRouter)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		htmlContent := fmt.Sprintf(`
        <html>
        <body>
            <h1>Welcome, Chirpy Admin</h1>
            <p>Chirpy has been visited %d times!</p>
        </body>
        </html>
    `, apiCfg.fileserverHits)

		w.Write([]byte(htmlContent))
	})

	router.Mount("/admin", adminRouter)
	corsMux := middlewareCors(router)
	server := &http.Server{
		Handler: corsMux,
		Addr:    "localhost:" + port,
	}
	fmt.Printf("Server listening on port %s\n", port)
	server.ListenAndServe()
}
