package internal

import (
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
)

const _fileMode = 0666

// New creates a new bolt database with the provided config.
func New(file string, opts *bolt.Options) (*bolt.DB, error) {
	db, err := bolt.Open(file, _fileMode, opts)
	if err != nil {
		return nil, errors.Wrap(err, "bolt.Open")
	}

	tx, err := db.Begin(true)
	if err != nil {
		return nil, errors.Wrap(err, "_db.Begin")
	}
	defer tx.Rollback()

	if _, err = tx.CreateBucketIfNotExists(SnapshotBucket); err != nil {
		return nil, errors.Wrapf(err, "error creating bucket %s", string(SnapshotBucket))
	}

	if _, err = tx.CreateBucketIfNotExists(MainBucket); err != nil {
		return nil, errors.Wrapf(err, "error creating bucket %s", string(MainBucket))
	}

	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "tx.Commit")
	}

	return db, nil
}
