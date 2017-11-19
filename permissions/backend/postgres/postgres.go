package postgres

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx"
	"github.com/justanotherorganization/justanotherbotkit/permissions/backend"
)

const (
	sqlInsertUser = `INSERT into "%s" (id, name, perms, groups)
	VALUES ($1, $2, $3, $4);`
	sqlSelectUser  = `SELECT id from "%s" where id = $1;`
	sqlUpdatePerms = `UPDATE "%s"
	SET perms = $1
	WHERE id = $2;`
	sqlSelectPerms = `SELECT perms from "%s" where id = $1;`
)

type (
	pg struct {
		pool  *pgx.ConnPool
		table string
	}
)

// New returns a new DB.
func New(conf *pgx.ConnPoolConfig, table string) (backend.DB, error) {
	if conf == nil {
		return nil, errors.New("conf cannot be nil")
	}

	if table == "" {
		return nil, errors.New("table must be set")
	}

	pool, err := pgx.NewConnPool(*conf)
	if err != nil {
		return nil, err
	}

	return &pg{
		pool:  pool,
		table: table,
	}, nil
}

// WriteUser writes a new ID to postgres.
func (p *pg) WriteUser(id, name string) error {
	tx, err := p.pool.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			// TODO: properly handle this error.
			tx.Rollback()
		}
	}()

	if _, err := tx.Exec(fmt.Sprintf(sqlInsertUser, p.table), id, name, []string{}, []int32{}); err != nil {
		return err
	}

	return tx.Commit()
}

// CheckUser checks if a user exists in postgres.
func (p *pg) CheckUser(id string) (bool, error) {
	result := p.pool.QueryRow(fmt.Sprintf(sqlSelectUser, p.table), id)
	var _id string
	if err := result.Scan(&_id); err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// SetPerms replaces the current slice of permissions with a new one.
// This is destructive and should be used to add and remove perms.
func (p *pg) SetPerms(id string, perm ...string) error {
	tx, err := p.pool.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			// TODO: properly handle this error.
			tx.Rollback()
		}
	}()

	if _, err := tx.Exec(fmt.Sprintf(sqlUpdatePerms, p.table), perm, id); err != nil {
		return err
	}

	return tx.Commit()
}

// GetPerms get the perms for a user from postgres.
func (p *pg) GetPerms(id string) ([]string, error) {
	result := p.pool.QueryRow(fmt.Sprintf(sqlSelectPerms, p.table), id)
	var s []string
	if err := result.Scan(&s); err != nil {
		return nil, err
	}

	return s, nil
}
