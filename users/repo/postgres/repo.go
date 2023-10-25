package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"justanother.org/justanotherbotkit/users/repo"
)

type Repo struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) Repo {
	return Repo{
		pool: pool,
	}
}

func (r Repo) CreateUser(ctx context.Context, u repo.User) (repo.User, error) {
	const sql = `INSERT INTO jab.users (network_id, name, permissions)
	VALUES ($1, $2, $3)
	RETURNING id
	;`

	err := r.pool.QueryRow(ctx, sql, u.NetworkID, u.Name, u.Permissions).Scan(&u.ID)
	return u, err
}

func (r Repo) UpdateUser(ctx context.Context, u repo.User) (repo.User, error) {
	const sql = `UPDATE jab.users
	SET name = $2, permissions = $3
	WHERE id = $1
	;`

	_, err := r.pool.Exec(ctx, sql, u.ID, u.Name, u.Permissions)
	return u, err
}

func (r Repo) DeleteUser(ctx context.Context, id uint) error {
	const sql = `DELETE FROM jab.users WHERE id = $1;`
	_, err := r.pool.Exec(ctx, sql, id)
	return err
}

func (r Repo) GetUser(ctx context.Context, id uint) (repo.User, error) {
	return r.getUser(ctx, "WHERE id = $1", id)
}

func (r Repo) GetUserByNetworkID(ctx context.Context, networkID string) (repo.User, error) {
	return r.getUser(ctx, "WHERE network_id = $1", networkID)
}

func (r Repo) getUser(ctx context.Context, whereClause string, params ...any) (repo.User, error) {
	const sql = `SELECT id, network_id, name, permissions FROM jab.users`
	_sql := sql + "\n" + whereClause

	var user repo.User
	err := r.pool.QueryRow(ctx, _sql, params...).Scan(
		&user.ID,
		&user.NetworkID,
		&user.Name,
		&user.Permissions,
	)
	return user, err
}

func (r Repo) ListUsers(ctx context.Context, cursor, limit uint, order repo.PageOrder) ([]repo.User, error) {
	if limit == 0 {
		limit = 100
	}

	if order == "" {
		order = repo.DESC
	}

	var (
		args []interface{}
		b    strings.Builder
	)
	b.WriteString("SELECT id, network_id, name, permissions FROM jab.users\n")
	if cursor > 0 {
		switch order {
		case repo.ASC:
			b.WriteString(fmt.Sprintf("WHERE id > $1 ORDER BY id %s\n", order))
		case repo.DESC:
			b.WriteString(fmt.Sprintf("WHERE id < $1 ORDER BY id %s\n", order))
		}

		args = append(args, cursor)
	} else {
		b.WriteString(fmt.Sprintf("ORDER BY id %s\n", order))
	}

	b.WriteString(fmt.Sprintf("LIMIT $%d", len(args)+1))
	args = append(args, limit)

	rows, err := r.pool.Query(ctx, b.String(), args...)
	if err != nil {
		return nil, err
	}

	var users []repo.User
	for rows.Next() {
		var user repo.User
		if err = rows.Scan(
			&user.ID,
			&user.NetworkID,
			&user.Name,
			&user.Permissions,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
