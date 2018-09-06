package transport

import "github.com/justanotherorganization/justanotherbotkit/permissions"

// Channel is an interface for channels, or channel in the case of there being
//  no channels and only one place to talk.
type Channel interface {
	GetUsers() []*permissions.User
	SendQuickMessage(string) error
	SendMessage(Message) error
}
