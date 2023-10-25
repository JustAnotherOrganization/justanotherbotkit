package repo

//go:generate go install github.com/golang/mock/mockgen@latest
//go:generate mockgen -build_flags=--mod=mod --destination=./mock/mock_repo.go -package=mock justanother.org/justanotherbotkit/users/repo Repo

import (
	"context"
)

type PageOrder string

const (
	ASC  PageOrder = "ASC"
	DESC PageOrder = "DESC"
)

type Repo interface {
	CreateUser(ctx context.Context, u User) (User, error)
	GetUser(ctx context.Context, id uint) (User, error)
	GetUserByNetworkID(ctx context.Context, networkID string) (User, error)
	ListUsers(ctx context.Context, cursor, limit uint, order PageOrder) ([]User, error)
	UpdateUser(ctx context.Context, u User) (User, error)
	DeleteUser(ctx context.Context, id uint) error
}
