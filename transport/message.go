package transport

import (
	"time"

	"github.com/justanotherorganization/justanotherbotkit/permissions"
)

// Message ...
type Message interface {
	GetSender() *permissions.User
	GetBody() string
	QuickReply(string) error
	SendReply(Message) error
	GetChannel() Channel
	GetSentDate() time.Time
	GetRecievedDate() time.Time
}
