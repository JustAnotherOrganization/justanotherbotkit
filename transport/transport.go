package transport // package github.com/justanotherorganization/justanotherbotkit/transport

import (
	"context"
	"errors"
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
		// GetUsers returns a list of all known users.
		GetUsers() ([]*User, error)
		// GetConversation gets a private conversation for the given user ID.
		// This varies depending on network (may just be the user ID) but is required for Slack.
		GetConversation(userID string) (string, error)
	}

	// Config is a transport configuration.
	Config struct {
		Token       string
		IgnoreUsers []string
	}
)

var (
	// ErrNilConfig is returned if no config is passed into a transport 'New' function.
	ErrNilConfig = errors.New("cfg cannot be nil")
	// ErrEmptyToken is returned if no token is provided in the given config.
	ErrEmptyToken = errors.New("token cannot be empty")
	// ErrNilTransport should be returned in places where a transport is required but is not set.
	ErrNilTransport = errors.New("transport cannot be nil")
	// ErrUserNotFound is returned if teh transport cannot locate a given user.
	ErrUserNotFound = errors.New("user not found")
)

// Validate a configuration.
func (c *Config) Validate() error {
	if c == nil {
		return ErrNilConfig
	}

	if c.Token == "" {
		return ErrEmptyToken
	}

	return nil
}
