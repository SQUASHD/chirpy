package models

import "time"

type Tokens struct {
	Tokens map[string]time.Time `json:"revoked_tokens"`
}
