package database

import (
	"encoding/json"
	"errors"
	"github.com/squashd/chirpy/internal/models"
	"os"
	"sync"
	"time"
)

var ErrNotExist = errors.New("resource does not exist")

type DB struct {
	path  string
	debug bool
	mu    *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]models.Chirp `json:"chirps"`
	Users  map[int]models.User  `json:"users"`
	Tokens map[string]time.Time `json:"revoked_tokens"`
}

func NewDB(path string, debug bool) (*DB, error) {
	db := &DB{
		path:  path,
		debug: debug,
		mu:    &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps: map[int]models.Chirp{},
		Users:  map[int]models.User{},
		Tokens: map[string]time.Time{},
	}
	return db.writeDB(dbStructure)
}

func (db *DB) ensureDB() error {
	if db.debug {
		return db.createDB()
	}
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbStructure := DBStructure{}
	dat, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}
	err = json.Unmarshal(dat, &dbStructure)
	if err != nil {
		return dbStructure, err
	}

	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, dat, 0600)
	if err != nil {
		return err
	}
	return nil
}
