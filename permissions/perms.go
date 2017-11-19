package permissions

import (
	"github.com/justanotherorganization/justanotherbotkit/permissions/backend"
	"github.com/pkg/errors"
)

type (
	// Manager manages user permissions in postgres.
	Manager struct {
		db backend.DB
	}
)

// NewManager returns a newly created Permissions manager.
func NewManager(db backend.DB) (*Manager, error) {
	if db == nil {
		return nil, errors.New("db cannot be  nil")
	}

	return &Manager{
		db: db,
	}, nil
}

// NewUser creates a new user in the database.
func (pm *Manager) NewUser(id, name string, perms ...string) (*User, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}

	// Write the user to the database.
	if err := pm.db.WriteUser(id, name); err != nil {
		return nil, err
	}

	u := &User{
		ID: id,
		pm: pm,
	}

	if err := u.AddPerms(perms...); err != nil {
		return nil, err
	}

	return u, nil
}

// GetUser checks if the user exists in the database and returns a user if so.
func (pm *Manager) GetUser(id string) (*User, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}

	exists, err := pm.db.CheckUser(id)
	if err != nil {
		return nil, err
	}

	if exists {
		return &User{
			ID: id,
			pm: pm,
		}, nil
	}

	return nil, nil
}
