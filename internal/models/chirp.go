package models

type Chirp struct {
	AuthorID int    `json:"author_id"`
	ID       int    `json:"id"`
	Body     string `json:"body"`
}
