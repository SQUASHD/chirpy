package main

import (
	"encoding/json"
	"fmt"
	"github.com/squashd/chirpy/internal/auth"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	issuer, err := auth.GetJWTIssuer(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT issuer")
		return
	}

	if issuer != "chirpy-access" {
		respondWithError(w, http.StatusUnauthorized, "Not a valid access token")
		return
	}

	subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		fmt.Printf("Error validating JWT: %v\n", err)
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	userIDInt, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
		return
	}

	user, err := cfg.DB.UpdateUser(userIDInt, params.Email, hashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}
