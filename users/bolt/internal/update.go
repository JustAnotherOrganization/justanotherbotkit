package internal

import (
	"context"

	"github.com/justanotherorganization/justanotherbotkit/users"
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
)

type updateEvent struct {
	_e event
}

// NewUpdateEvent ...
func NewUpdateEvent(u users.User, v uint32) (Event, error) {
	if err := validateUser(u); err != nil {
		return nil, err
	}

	if v == 0 {
		return nil, errors.New("v should represent the previous snapshot version")
	}

	return &updateEvent{
		_e: event{
			snap: newSnapshot(u, v),
		},
	}, nil
}

func (e *updateEvent) Apply(ctx context.Context, _db *bolt.DB) (users.User, error) {
	if err := checkPrev([]byte(e._e.snap.User.GetID()), e._e.snap.Version, _db); err != nil {
		return e._e.snap.User, err
	}

	// Overwrite immutable values from u here (if any).

	return e._e._save(ctx, _db)
}
