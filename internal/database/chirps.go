package database

import "github.com/squashd/chirpy/internal/models"

func (db *DB) CreateChirp(body string, userId int) (models.Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return models.Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := models.Chirp{
		AuthorID: userId,
		ID:       id,
		Body:     body,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return models.Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]models.Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]models.Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) GetChirpsByAuthorId(id int) ([]models.Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]models.Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		if chirp.AuthorID == id {
			chirps = append(chirps, chirp)
		}
	}

	return chirps, nil
}

func (db *DB) GetChirp(id int) (models.Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return models.Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return models.Chirp{}, ErrNotExist
	}

	return chirp, nil
}

func (db DB) DeleteChirp(id int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	delete(dbStructure.Chirps, id)

	return nil
}
