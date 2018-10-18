package internal

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/justanotherorganization/justanotherbotkit/users"
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
)

type deleteEvent struct {
	_e event
}

// NewDeleteEvent ...
func NewDeleteEvent(u users.User, v uint32) (Event, error) {
	if strings.TrimSpace(u.GetID()) == "" {
		return nil, errors.New("user.ID cannot be empty")
	}

	if v == 0 {
		return nil, errors.New("v should represent the current snapshot version")
	}

	return &deleteEvent{
		_e: event{
			snap: newSnapshot(u, v),
		},
	}, nil
}

func (e *deleteEvent) Apply(ctx context.Context, _db *bolt.DB) (users.User, error) {
	if err := checkPrev([]byte(e._e.snap.User.GetID()), e._e.snap.Version, _db); err != nil {
		return e._e.snap.User, err
	}

	byt, err := json.Marshal(e._e.snap)
	if err != nil {
		return e._e.snap.User, err
	}

	if err = _db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(SnapshotBucket)
		err := b.Put([]byte(e._e.snap.ID), byt)
		if err != nil {
			return err
		}

		b = tx.Bucket(MainBucket)
		return b.Delete([]byte(e._e.snap.User.GetID()))
	}); err != nil {
		return e._e.snap.User, errors.Wrap(err, "_db.Update")
	}

	return nil, nil
}
