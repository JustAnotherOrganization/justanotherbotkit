package internal

import (
	"fmt"

	"github.com/justanotherorganization/justanotherbotkit/users"
)

// SnapshotBucket ...
var SnapshotBucket = []byte("users_snapshots")

type (
	// Snapshot represents a snapshot of a value in a database.
	Snapshot struct {
		ID      string
		User    users.User
		Version uint32
	}
)

func newSnapshot(u users.User, v uint32) *Snapshot {
	return &Snapshot{
		ID:      fmt.Sprintf("%s%d", u.GetID(), v),
		User:    u,
		Version: v,
	}
}
