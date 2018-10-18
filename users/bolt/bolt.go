// Package bolt provides a boltdb implementation of the justanotherbotkit/users/DB interface.
// It provides this in the form of 2 event driven snapshot buckets; this means that the snapshot
// bucket will always retain a previous version of a DB record.
// While the context arguments are currently ignored future updates will use the context to house
// metadata pertaining to a given change.
package bolt // import "github.com/justanotherorganization/justanotherbotkit/users/bolt"

import (
	"context"

	"github.com/justanotherorganization/justanotherbotkit/users"
	"github.com/justanotherorganization/justanotherbotkit/users/bolt/internal"
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
)

type (
	// DB provides a Bolt implementation of a users.DB
	DB struct {
		_db *bolt.DB
	}
)

// Static type checking
var _ users.DB = &DB{}

// New returns a new users.DB running over bolt.
func New(cfg *Config) (db *DB, err error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	db = new(DB)
	db._db, err = internal.New(cfg.File, cfg.Options)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// CreateUser creates a new user and saves it to the database.
func (db *DB) CreateUser(ctx context.Context, u users.User) (users.User, error) {
	e, err := internal.NewCreateEvent(u)
	if err != nil {
		return nil, errors.Wrap(err, "NewCreateEvent")
	}

	return e.Apply(ctx, db._db)
}

// GetUser gets a user from the database using the ID.
func (db *DB) GetUser(ctx context.Context, id string) (users.User, error) {
	return internal.Get(id, db._db)
}

// UpdateUser updates a user record in the database.
func (db *DB) UpdateUser(ctx context.Context, u users.User) (users.User, error) {
	prev, err := internal.GetSnap([]byte(u.GetID()), db._db)
	if err != nil {
		return nil, err
	}

	e, err := internal.NewUpdateEvent(u, prev.Version)
	if err != nil {
		return nil, errors.Wrap(err, "NewUpdateEvent")
	}

	return e.Apply(ctx, db._db)
}

// DeleteUser deletes a user record from the database.
func (db *DB) DeleteUser(ctx context.Context, id string) (users.User, error) {
	prev, err := internal.GetSnap([]byte(id), db._db)
	if err != nil {
		return nil, err
	}

	e, err := internal.NewDeleteEvent(prev.User, prev.Version)
	if err != nil {
		return nil, errors.Wrap(err, "NewDeleteEvent")
	}

	return e.Apply(ctx, db._db)
}
