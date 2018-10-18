package transport // package github.com/justanotherorganization/justanotherbotkit/transport

import "github.com/justanotherorganization/justanotherbotkit/transport/internal/proto"

type (
	// Channel wraps a pb.BaseChannel up with it's accomanied transport.
	Channel struct {
		*pb.BaseChannel
		Transport
	}
)
