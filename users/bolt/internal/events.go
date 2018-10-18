package internal

import (
	"context"
	"encoding/json"

	"github.com/justanotherorganization/justanotherbotkit/users"
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
)

// ErrNotImplemented is returned if Apply is attempted on an internal event type.
var ErrNotImplemented = errors.New("not implemented")

type (
	// Event represents a database event.
	Event interface {
		Apply(ctx context.Context, _db *bolt.DB) (users.User, error)
	}

	event struct {
		snap *Snapshot
	}
)

// Apply is a dummy apply function present in all events.
// It should be overwritten.
func (e *event) Apply(ctx context.Context, _db *bolt.DB) (users.User, error) {
	return nil, ErrNotImplemented
}

func (e *event) _save(ctx context.Context, _db *bolt.DB) (users.User, error) {
	byt, err := json.Marshal(e.snap)
	if err != nil {
		return e.snap.User, err
	}

	if err = _db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(SnapshotBucket)
		err := b.Put([]byte(e.snap.ID), byt)
		if err != nil {
			return err
		}

		b = tx.Bucket(MainBucket)
		return b.Put([]byte(e.snap.User.GetID()), byt)
	}); err != nil {
		return e.snap.User, errors.Wrap(err, "_db.Update")
	}

	return Get(e.snap.User.GetID(), _db)
}
