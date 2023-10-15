package models

type User struct {
	ID             int    `json:"id"`
	Email          string `json:"email"`
	HashedPassword string `json:"password"`
	IsChirpyRed    bool   `json:"is_chirpy_red"`
}
