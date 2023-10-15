package database

import "github.com/squashd/chirpy/internal/models"

func (db *DB) CreateUser(email string, hashedPassword string) (models.User, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return models.User{}, err
	}

	id := len(dbStructure.Users) + 1
	user := models.User{
		ID:             id,
		Email:          email,
		HashedPassword: hashedPassword,
		IsChirpyRed:    false,
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (db *DB) GetUser(id int) (models.User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return models.User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return models.User{}, ErrNotExist
	}

	return user, nil
}

func (db *DB) GetUserByEmail(email string) (models.User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return models.User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return models.User{}, ErrNotExist
}

func (db *DB) UpdateUser(id int, email, hashedPassword string) (models.User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return models.User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return models.User{}, ErrNotExist
	}

	user.Email = email
	user.HashedPassword = hashedPassword
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (db *DB) UpgradeUser(id int, isChirpyRed bool) (models.User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return models.User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return models.User{}, ErrNotExist
	}

	user.IsChirpyRed = isChirpyRed
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
