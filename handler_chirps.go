package main

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/squashd/chirpy/internal/auth"
	"github.com/squashd/chirpy/internal/models"
	"net/http"
	"sort"
	"strconv"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	authorIDStr := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort")

	var chirpSlice []models.Chirp
	var err error

	if authorIDStr != "" {
		authorId, err := strconv.Atoi(authorIDStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID")
			return
		}
		chirpSlice, err = cfg.DB.GetChirpsByAuthorId(authorId)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error retrieving chirps by author")
			return
		}
	} else {
		chirpSlice, err = cfg.DB.GetChirps()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error retrieving chirps")
			return
		}
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "DB error")
		return
	}

	if sortOrder == "asc" {
		sort.Slice(chirpSlice, func(i, j int) bool {
			return chirpSlice[i].ID < chirpSlice[j].ID
		})
	} else if sortOrder == "desc" {
		sort.Slice(chirpSlice, func(i, j int) bool {
			return chirpSlice[i].ID > chirpSlice[j].ID
		})
	}

	respondWithJSON(w, http.StatusOK, chirpSlice)
}

func (cfg *apiConfig) handlerChirpsPost(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	authorId, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := cfg.DB.CreateChirp(cleaned, authorId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, models.Chirp{
		AuthorID: chirp.AuthorID,
		ID:       chirp.ID,
		Body:     chirp.Body,
	})
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	cleaned := censorWords(body)
	return cleaned, nil
}
func (cfg *apiConfig) handlerChirpsGetById(w http.ResponseWriter, r *http.Request) {
	chirpId := chi.URLParam(r, "chirpId")
	id, err := strconv.Atoi(chirpId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}
	chirp, err := cfg.DB.GetChirp(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}
	respondWithJSON(w, http.StatusOK, chirp)

}

func (cfg *apiConfig) handlerChirpsDeleteById(w http.ResponseWriter, r *http.Request) {
	chirpId := chi.URLParam(r, "chirpId")
	id, err := strconv.Atoi(chirpId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}
	chirp, err := cfg.DB.GetChirp(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	authorId, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
		return
	}

	if chirp.AuthorID != authorId {
		respondWithError(w, http.StatusForbidden, "You can't delete someone else's chirp")
		return
	}

	err = cfg.DB.DeleteChirp(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, nil)

}
