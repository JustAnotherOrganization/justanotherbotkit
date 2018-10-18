package internal

import (
	"context"

	"github.com/justanotherorganization/justanotherbotkit/users"
	bolt "go.etcd.io/bbolt"
)

type createEvent struct {
	_e event
}

// NewCreateEvent ...
func NewCreateEvent(u users.User) (Event, error) {
	if err := validateUser(u); err != nil {
		return nil, err
	}

	return &createEvent{
		_e: event{
			snap: newSnapshot(u, 1),
		},
	}, nil
}

// Apply ...
func (e *createEvent) Apply(ctx context.Context, _db *bolt.DB) (users.User, error) {
	return e._e._save(ctx, _db)
}
