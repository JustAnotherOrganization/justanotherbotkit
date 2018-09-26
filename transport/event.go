package transport

import "github.com/justanotherorganization/justanotherbotkit/transport/proto"

type (
	// Event wraps a pb.BaseEvent up with it's accompanied transport.
	Event struct {
		*pb.BaseEvent
		Transport
	}
)
