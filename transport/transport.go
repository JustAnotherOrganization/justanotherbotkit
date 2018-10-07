package transport // package github.com/justanotherorganization/justanotherbotkit/transport

import (
	"context"
)

type (
	// Transport provides standard io access to a network transport layer.
	Transport interface {
		// TunnelEvents translates a specified network's events to Events tunneling them into evCh.
		// The session is terminated when ctx.Done returns.
		TunnelEvents(ctx context.Context, evCh chan *Event, errCh chan error)
		// SendMessage sends a message using the default format specified in the given transport.
		SendMessage(dest, msg string) error
		// SendEvent sends a new event over the given transport.
		SendEvent(ev *Event) error
		// Channels lists all the channels we have access to via this transport.
		Channels() ([]*Channel, error)
		// GetUser returns the full user data for the provided name or ID.
		GetUser(user string) (*User, error)
	}

	// Config is a transport configuration.
	Config struct {
		Token       string
		IgnoreUsers []string
	}
)
