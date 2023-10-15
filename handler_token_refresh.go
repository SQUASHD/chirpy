package main

import (
	"github.com/squashd/chirpy/internal/auth"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	issuers, err := auth.GetJWTIssuer(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}
	if issuers != "chirpy-refresh" {
		respondWithError(w, http.StatusUnauthorized, "Not a valid refresh token")
		return
	}
	revoked, err := cfg.DB.IsTokenRevoked(token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't check if token was revoked")
		return
	}
	if revoked {
		respondWithError(w, http.StatusUnauthorized, "Token has been revoked")
		return
	}

	subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	subjectInt, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
		return
	}

	accessToken, err := auth.MakeJWT(subjectInt, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})

}
