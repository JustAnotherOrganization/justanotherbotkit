package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"justanother.org/justanotherbotkit/users/repo"
	"justanother.org/justanotherbotkit/users/repo/postgres"
)

type fixture struct {
	pool *pgxpool.Pool
	repo postgres.Repo
}

func getFixture(t *testing.T) fixture {
	t.Helper()

	var (
		f   fixture
		err error
	)

	f.pool, err = pgxpool.New(context.Background(), postgres.DSN())
	require.NoError(t, err)

	f.repo = postgres.New(f.pool)
	return f
}

func TestRepo_CreateUser(t *testing.T) {
	f := getFixture(t)

	user := repo.User{
		NetworkID:   uuid.NewString(),
		Name:        "user-name",
		Permissions: []string{"view-members"},
	}

	created, err := f.repo.CreateUser(context.Background(), user)
	require.NoError(t, err)

	fmt.Println(created.ID)
	found, err := f.repo.GetUser(context.Background(), created.ID)
	require.NoError(t, err)

	assert.Equal(t, "user-name", found.Name)
	assert.Equal(t, []string{"view-members"}, found.Permissions)

	_, err = f.repo.CreateUser(context.Background(), found)
	require.Error(t, err)
}

func TestRepo_UpdateUser(t *testing.T) {
	f := getFixture(t)

	user := repo.User{
		NetworkID:   uuid.NewString(),
		Name:        "user-name",
		Permissions: []string{"view-members"},
	}

	user, err := f.repo.CreateUser(context.Background(), user)
	require.NoError(t, err)

	user.Name = "user-name-2"
	user.Permissions = append(user.Permissions, "delete-members")

	_, err = f.repo.UpdateUser(context.Background(), user)
	require.NoError(t, err)

	user, err = f.repo.GetUser(context.Background(), user.ID)
	require.NoError(t, err)

	assert.Equal(t, "user-name-2", user.Name)
	assert.Equal(t, []string{"view-members", "delete-members"}, user.Permissions)
}

func TestRepo_GetUser(t *testing.T) {
	t.Run("user exists", func(tt *testing.T) {
		tt.Skip("if this test existed it would match the first part of the create test")
	})

	t.Run("user doesn't exist", func(tt *testing.T) {
		f := getFixture(tt)
		_, err := f.repo.GetUser(context.Background(), 9000)
		require.Error(tt, err)
	})
}

func TestRepo_ListUsers(t *testing.T) {
	f := getFixture(t)

	// FIXME: this is gross
	_, err := f.pool.Exec(context.Background(), "TRUNCATE jab.users;")
	require.NoError(t, err)

	user := repo.User{
		NetworkID:   uuid.NewString(),
		Name:        "user-name-1",
		Permissions: []string{"view-members"},
	}

	_, err = f.repo.CreateUser(context.Background(), user)
	require.NoError(t, err)

	user.NetworkID = uuid.NewString()
	user.Name = "user-name-2"
	_, err = f.repo.CreateUser(context.Background(), user)
	require.NoError(t, err)

	user.NetworkID = uuid.NewString()
	user.Name = "user-name-3"
	_, err = f.repo.CreateUser(context.Background(), user)
	require.NoError(t, err)

	t.Run("descending (default)", func(tt *testing.T) {
		users, err := f.repo.ListUsers(context.Background(), 0, 1, "")
		require.NoError(tt, err)
		assert.Len(tt, users, 1)
		assert.Equal(tt, "user-name-3", users[0].Name)

		users, err = f.repo.ListUsers(context.Background(), users[0].ID, 1, "")
		require.NoError(tt, err)
		assert.Len(tt, users, 1)
		assert.Equal(tt, "user-name-2", users[0].Name)

		users, err = f.repo.ListUsers(context.Background(), users[0].ID, 1, "")
		require.NoError(tt, err)
		assert.Len(tt, users, 1)
		assert.Equal(tt, "user-name-1", users[0].Name)
	})

	t.Run("ascending", func(tt *testing.T) {
		users, err := f.repo.ListUsers(context.Background(), 0, 1, repo.ASC)
		require.NoError(tt, err)
		assert.Len(tt, users, 1)
		assert.Equal(tt, "user-name-1", users[0].Name)

		users, err = f.repo.ListUsers(context.Background(), users[0].ID, 1, repo.ASC)
		require.NoError(tt, err)
		assert.Len(tt, users, 1)
		assert.Equal(tt, "user-name-2", users[0].Name)

		users, err = f.repo.ListUsers(context.Background(), users[0].ID, 1, repo.ASC)
		require.NoError(tt, err)
		assert.Len(tt, users, 1)
		assert.Equal(tt, "user-name-3", users[0].Name)
	})
}

func TestRepo_DeleteUser(t *testing.T) {
	f := getFixture(t)

	user := repo.User{
		NetworkID:   uuid.NewString(),
		Name:        "user-name",
		Permissions: []string{"view-members"},
	}

	user, err := f.repo.CreateUser(context.Background(), user)
	require.NoError(t, err)

	err = f.repo.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)

	_, err = f.repo.GetUser(context.Background(), user.ID)
	require.Error(t, err)

	err = f.repo.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)
}
