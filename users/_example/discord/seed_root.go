package main

import (
	"context"

	"justanother.org/justanotherbotkit/users/repo"
)

func seedRoot(ctx context.Context, r repo.Repo) error {
	_, err := r.CreateUser(ctx, repo.User{
		NetworkID:   "",
		Permissions: []string{"root"},
	})
	return err
}
