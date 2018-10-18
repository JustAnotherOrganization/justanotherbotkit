package users // import "github.com/justanotherorganization/justanotherbotkit/users"

type (
	// User represents a bot user. This is satisfied by transport.User
	User interface {
		GetID() string
		GetName() string
		GetPermissions() []string
	}
)
