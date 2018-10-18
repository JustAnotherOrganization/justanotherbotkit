package internal

import (
	"context"
	"testing"

	"github.com/justanotherorganization/justanotherbotkit/internal/test"
	"github.com/justanotherorganization/justanotherbotkit/users"
	bolt "go.etcd.io/bbolt"
)

func createUser(tb testing.TB, db *bolt.DB) users.User {
	u := newTestUser(tb)
	e, err := NewCreateEvent(u)
	test.OK(tb, err)
	_u, err := e.Apply(context.Background(), db)
	test.OK(tb, err)
	test.Assert(tb, u.GetID() == _u.GetID(), "ids must be equal")
	return _u
}

func TestCreate(t *testing.T) {
	db := initDB(t)
	defer db.Close()
	_ = createUser(t, db)
}
