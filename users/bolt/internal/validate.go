package internal

import (
	"strings"

	"github.com/justanotherorganization/justanotherbotkit/users"
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
)

func validateUser(u users.User) error {
	if strings.TrimSpace(u.GetID()) == "" {
		return errors.New("user.ID cannot be empty")
	}

	if strings.TrimSpace(u.GetName()) == "" {
		return errors.New("user.Name cannot be empty")
	}

	return nil
}

func checkPrev(id []byte, v uint32, _db *bolt.DB) error {
	prev, err := GetSnap(id, _db)
	if err != nil {
		return err
	}

	if pv := prev.Version; v != pv {
		return errors.Errorf("version %d does not match that of database %d", v, pv)
	}

	return nil
}
