package main

import (
	"fmt"
	"net/http"
)

func main() {
	const port = "8080"
	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)
	server := &http.Server{
		Handler: corsMux,
		Addr:    "localhost:" + port,
	}
	fmt.Printf("Server listening on port %s\n", port)
	server.ListenAndServe()
}
