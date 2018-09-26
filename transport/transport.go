package transport

import (
	"context"
)

type (
	// Transport provides standard io access to a network transport layer.
	Transport interface {
		TunnelEvents(ctx context.Context, evCh chan *Event, errCh chan error)
		// FIXME: this isn't good enough, we need to be able to set options, which
		// can vary by network...
		SendMessage(dest, msg string) error
	}
)
