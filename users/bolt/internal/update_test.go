package internal

import (
	"context"
	"testing"

	"github.com/justanotherorganization/justanotherbotkit/internal/test"
	"github.com/justanotherorganization/justanotherbotkit/proto"
)

func TestUpdate(t *testing.T) {
	db := initDB(t)
	defer db.Close()
	u := createUser(t, db)
	_u := u.(*pb.BaseUser)
	_u.Name = "FOO"

	e, err := NewUpdateEvent(_u, 1)
	test.OK(t, err)
	fu, err := e.Apply(context.Background(), db)
	test.OK(t, err)
	test.Assert(t, u.GetID() == fu.GetID(), "ids must be equal")
	test.Assert(t, u.GetName() == "FOO", "name did not update")
}
