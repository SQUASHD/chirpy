package database

import "time"

func (db *DB) AddToken(token string) error {

	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	dbStructure.Tokens[token] = time.Now()

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) IsTokenRevoked(token string) (bool, error) {

	dbStructure, err := db.loadDB()
	if err != nil {
		return false, err
	}

	_, ok := dbStructure.Tokens[token]
	if !ok {
		return false, nil
	}

	return true, nil
}
