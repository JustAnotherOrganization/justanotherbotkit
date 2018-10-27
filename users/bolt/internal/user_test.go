package internal

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/justanotherorganization/justanotherbotkit/internal/test"
	"github.com/justanotherorganization/justanotherbotkit/proto"
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
