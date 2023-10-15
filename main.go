package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/squashd/chirpy/internal/database"
	"net/http"
	"os"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	jwtSecret      string
	polkaApiKey    string
}

var dbg = flag.Bool("debug", false, "Enable debug mode")

func main() {
	flag.Parse()

	const port = "8080"
	const filepathRoot = "."

	db, err := database.NewDB("database.json", *dbg)
	if err != nil {
		panic(err)
	}

	err = godotenv.Load()
	if err != nil {
		panic(err)
	}

	apiCfg := &apiConfig{
		fileserverHits: 0,
		DB:             db,
		jwtSecret:      os.Getenv("JWT_SECRET"),
		polkaApiKey:    os.Getenv("POLKA_API_KEY"),
	}

	if err != nil {
		panic(err)
	}

	router := chi.NewRouter()
	fs := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	router.Handle("/app/*", fs)
	router.Handle("/app", fs)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handlerCheckHealth)
	apiRouter.Get("/metrics", apiCfg.handlerGetMetrics)
	apiRouter.Get("/reset", apiCfg.handlerReset)

	apiRouter.Get("/chirps", apiCfg.handlerChirpsGet)
	apiRouter.Post("/chirps", apiCfg.handlerChirpsPost)
	apiRouter.Get("/chirps/{chirpId}", apiCfg.handlerChirpsGetById)
	apiRouter.Delete("/chirps/{chirpId}", apiCfg.handlerChirpsDeleteById)

	apiRouter.Post("/users", apiCfg.handlerUsersCreate)
	apiRouter.Put("/users", apiCfg.handlerUsersUpdate)

	apiRouter.Post("/refresh", apiCfg.handlerRefresh)
	apiRouter.Post("/revoke", apiCfg.handlerRevoke)
	apiRouter.Post("/login", apiCfg.handlerLogin)

	apiRouter.Post("/polka/webhooks", apiCfg.handlerPolkaWebhooks)

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
