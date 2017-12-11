package backend

type (
	// DB provides access to the underlying database.
	DB interface {
		WriteUser(userID, name string) error
		CheckUser(userID string) (hasPermission bool, err error)
		GetPerms(userID string) (permissions []string, err error)
		SetPerms(userID string, permissions ...string) error
	}
)
