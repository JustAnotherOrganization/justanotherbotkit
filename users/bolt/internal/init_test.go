package internal

import (
	"testing"

	"github.com/justanotherorganization/justanotherbotkit/internal/test"
	bolt "go.etcd.io/bbolt"
)

const _path = "./test.db"

func initDB(tb testing.TB) *bolt.DB {
	db, err := New(_path, nil)
	test.OK(tb, err)
	return db
}
