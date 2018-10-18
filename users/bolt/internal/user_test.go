package internal

import (
	"testing"

	"github.com/justanotherorganization/justanotherbotkit/proto"
	uuid "github.com/satori/go.uuid"
)

func newTestUser(tb testing.TB) *pb.BaseUser {
	uid := uuid.NewV4()
	id := uid.String()

	return &pb.BaseUser{
		ID:   id,
		Name: id,
	}
}
