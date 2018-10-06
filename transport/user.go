package transport // package github.com/justanotherorganization/justanotherbotkit/transport

import "github.com/justanotherorganization/justanotherbotkit/transport/proto"

type (
	// User wraps a pb.BaseUser up with it's accompanied transport.
	User struct {
		*pb.BaseUser
		Transport
	}
)
