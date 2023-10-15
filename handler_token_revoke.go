package main

import (
	"github.com/squashd/chirpy/internal/auth"
	"net/http"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}
	alreadyRevoked, err := cfg.DB.IsTokenRevoked(token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't check if token was revoked")
		return
	}
	if alreadyRevoked {
		respondWithError(w, http.StatusUnauthorized, "Token already revoked")
		return
	}
	err = cfg.DB.AddToken(token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't add token to revoked list")
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}
