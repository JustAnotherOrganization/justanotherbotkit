package transport // package github.com/justanotherorganization/justanotherbotkit/transport

import "github.com/justanotherorganization/justanotherbotkit/transport/internal/proto"

type (
	// Event wraps a pb.BaseEvent up with it's accompanied transport.
	Event struct {
		*pb.BaseEvent
		Transport
	}
)
