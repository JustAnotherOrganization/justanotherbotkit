package internal

import (
	"fmt"
	"testing"

	"github.com/justanotherorganization/justanotherbotkit/internal/test"
)

func TestGet(t *testing.T) {
	db := initDB(t)
	defer db.Close()
	u := createUser(t, db)
	_u, err := Get(u.GetID(), db)
	test.OK(t, err)
	test.Assert(t, u.GetID() == _u.GetID(), "ids must be equal")
	test.Assert(t, u.GetName() == _u.GetName(), "names must be equal")
}

func TestGetSnap(t *testing.T) {
	db := initDB(t)
	defer db.Close()
	u := createUser(t, db)
	s, err := GetSnap([]byte(u.GetID()), db)
	test.OK(t, err)
	test.Assert(t, s.Version == 1, "version %d should be 1", s.Version)
	test.Assert(t, s.ID == fmt.Sprintf("%s%d", u.GetID(), s.Version), "snapshot ID is not valid")
	test.Assert(t, u.GetID() == s.User.GetID(), "ids must be equal")
	test.Assert(t, u.GetName() == s.User.GetName(), "names must be equal")
}
