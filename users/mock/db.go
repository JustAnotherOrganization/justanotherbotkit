package mock // import "github.com/justanotherorganization/justanotherbotkit/users/mock"

import (
	"context"
	"sync"

	"github.com/justanotherorganization/justanotherbotkit/users"
)

// DB is a mock database.
type DB struct {
	sync.RWMutex
	m map[string]users.User
}

// Static type checking
var _ users.DB = &DB{}

// New returns a new in-memory user database.
func New() (db *DB) {
	db = new(DB)
	db.m = make(map[string]users.User)
	return db
}

// CreateUser creates a new user and saves it in memory.
func (db *DB) CreateUser(_ context.Context, u users.User) (users.User, error) {
	db.Lock()
	defer db.Unlock()
	db.m[u.GetID()] = u
	return u, nil
}

// GetUser gets a user from the database using the ID.
func (db *DB) GetUser(_ context.Context, id string) (users.User, error) {
	db.RLock()
	defer db.RUnlock()
	u := db.m[id]
	return u, nil
}

// UpdateUser udpates a user record in memory.
func (db *DB) UpdateUser(_ context.Context, u users.User) (users.User, error) {
	db.Lock()
	defer db.Unlock()
	db.m[u.GetID()] = u
	return u, nil
}

// DeleteUser deletes a user record from memory.
func (db *DB) DeleteUser(_ context.Context, id string) error {
	db.Lock()
	defer db.Unlock()
	_, ok := db.m[id]
	if ok {
		delete(db.m, id)
	}

	return nil
}
