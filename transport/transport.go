package transport

//go:generate go install github.com/golang/mock/mockgen@latest
//go:generate mockgen -build_flags=--mod=mod -destination=./impl/mock/mock_transport.go -package=mock justanother.org/justanotherbotkit/transport Transport

import (
	"context"
)

type (
	Transport interface {
		SendMessage(ctx context.Context, dest string, options ...MsgOption) error
		MessageEventHandler(h func(ctx context.Context, ev Event) error, errHandler func(err error))
	}

	Config struct {
		Token       string
		IgnoreUsers []string
	}
)
