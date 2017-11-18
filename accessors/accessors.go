package accessors

type (
	// Accessor provides access to a network.
	Accessor interface {
		// TunnelEvents tunnels events from a connected network into a provided
		// channel (this should be ran within it's own goroutine), a separate
		// error channel should also be provided (and checked).
		TunnelEvents(chan MessageEvent, chan error, chan error)
		// SendMessage allows a process to send a message through the Accessor
		// to a connected network.
		SendMessage(msg string, dest string) error
	}
)
