package internal

import (
	"encoding/json"
	"strings"

	"github.com/justanotherorganization/justanotherbotkit/proto"
	"github.com/justanotherorganization/justanotherbotkit/users"
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
)

var (
	// MainBucket ...
	MainBucket = []byte("users_main")
	// ErrEmptyID is returned if Get is called without an ID set.
	ErrEmptyID = errors.New("id cannot be empty")
)

// Get ...
func Get(id string, _db *bolt.DB) (users.User, error) {
	if strings.TrimSpace(id) == "" {
		return nil, ErrEmptyID
	}

	s, err := GetSnap([]byte(id), _db)
	if err != nil {
		return nil, err
	}

	if s == nil {
		return nil, nil
	}

	return s.User, nil
}

// GetSnap ...
func GetSnap(id []byte, _db *bolt.DB) (*Snapshot, error) {
	if id == nil {
		return nil, ErrEmptyID
	}

	errNotFound := errors.New("not found")
	var byt []byte
	if err := _db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(MainBucket)
		byt = b.Get([]byte(id))
		if len(byt) == 0 {
			return errNotFound
		}

		return nil
	}); err != nil {
		if err == errNotFound {
			return nil, nil
		}

		return nil, errors.Wrap(err, "_db.View")
	}

	type snapshot struct {
		ID      string
		User    *pb.BaseUser
		Version uint32
	}
	var s snapshot
	if err := json.Unmarshal(byt, &s); err != nil {
		return nil, err
	}

	return &Snapshot{
		ID:      s.ID,
		User:    s.User,
		Version: s.Version,
	}, nil
}
