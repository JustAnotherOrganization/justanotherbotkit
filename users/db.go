package users // import "github.com/justanotherorganization/justanotherbotkit/users"

import "context"

type (
	// DB provides access to an underlying database.
	DB interface {
		// CreateUser creates a new user and saves it to the database.
		CreateUser(ctx context.Context, u User) (User, error)
		// GetUser gets a user from the database using the ID.
		GetUser(ctx context.Context, id string) (User, error)
		// UpdateUser updates a user record in the database.
		UpdateUser(ctx context.Context, u User) (User, error)
		// DeleteUser deletes a user record from the database.
		DeleteUser(ctx context.Context, id string) (User, error)
	}
)
