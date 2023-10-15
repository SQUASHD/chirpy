package main

import "strings"

func censorWords(input string) string {
	wordsToCensor := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}

	split := strings.Fields(input)

	for i, word := range split {
		if _, exists := wordsToCensor[strings.ToLower(word)]; exists {
			split[i] = "****"
		}
	}

	return strings.Join(split, " ")
}
