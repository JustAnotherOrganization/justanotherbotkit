package transport

//go:generate go install github.com/golang/mock/mockgen@latest
//go:generate mockgen -build_flags=--mod=mod -destination=./impl/mock/mock_transport.go -package=mock justanother.org/justanotherbotkit/transport Transport

import (
	"context"
)

type (
	Transport interface {
		Start(ctx context.Context) error
		SendMessage(ctx context.Context, dest string, options ...MsgOption) error
		MessageEventHandler(h func(ctx context.Context, ev Event) error, errHandler func(err error))
		// TODO: consider providing a way to register event handlers.
		// 	* for discord this is relatively easy but for slack the RTM doesn't _quite_ support this on its own
		// 		instead we'll have to implement a registration function inside the slack transport package.
		//
	}

	Config struct {
		Token       string
		IgnoreUsers []string
	}
)
