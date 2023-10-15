package main

import (
	"encoding/json"
	"errors"
	"github.com/squashd/chirpy/internal/auth"
	"net/http"
	"os"
)

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

	type response struct {
		Email       string `json:"email"`
		ID          int    `json:"id"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}

	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if apiKey != cfg.polkaApiKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusOK, nil)
		return
	}

	user, err := cfg.DB.UpgradeUser(params.Data.UserID, true)
	if errors.Is(err, os.ErrNotExist) {
		respondWithError(w, http.StatusNotFound, "Couldn't find user")
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't upgrade user")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Email:       user.Email,
		ID:          user.ID,
		IsChirpyRed: user.IsChirpyRed,
	})

}
