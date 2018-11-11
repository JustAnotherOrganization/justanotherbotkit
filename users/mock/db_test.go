package mock_test

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/justanotherorganization/justanotherbotkit/internal/test"
	"github.com/justanotherorganization/justanotherbotkit/proto"
	. "github.com/justanotherorganization/justanotherbotkit/users/mock"
)

func newTestUser(tb testing.TB) *pb.BaseUser {
	uid, err := uuid.NewV4()
	test.OK(tb, err)
	id := uid.String()

	return &pb.BaseUser{
		ID:   id,
		Name: id,
	}
}

func TestCRUD(t *testing.T) {
	ctx := context.Background()
	db := New()
	// create
	u := newTestUser(t)
	_u, err := db.CreateUser(ctx, u)
	test.OK(t, err)
	test.Assert(t, u.GetID() == _u.GetID(), "ids should be equal")
	test.Assert(t, u.GetName() == _u.GetName(), "names should be equal")
	// retrieve
	_u, err = db.GetUser(ctx, u.GetID())
	test.OK(t, err)
	test.Assert(t, u.GetID() == _u.GetID(), "ids should be equal")
	test.Assert(t, u.GetName() == _u.GetName(), "names should be equal")
	// update
	u.Name = "test"
	_u, err = db.UpdateUser(ctx, u)
	test.OK(t, err)
	test.Assert(t, u.GetID() == _u.GetID(), "ids should be equal")
	test.Assert(t, u.GetName() == _u.GetName(), "names should be equal")
	// delete
	err = db.DeleteUser(ctx, u.GetID())
	test.OK(t, err)
	_u, err = db.GetUser(ctx, u.GetID())
	test.OK(t, err)
	test.Assert(t, _u == nil, "user should not exist")
}
