package db

import "github.com/justanotherorganization/justanotherbotkit/users"

type (
	// DB provides access to the underlying database.
	DB interface {
		WriteUser(user users.User) (err error)
		ReadUser(userID interface{}) (user users.User, err error)
	}
)
