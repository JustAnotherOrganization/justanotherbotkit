package backend

type (
	// DB provides access to the underlying database.
	DB interface {
		WriteUser(id, name string) error
		CheckUser(string) (bool, error)
		GetPerms(string) ([]string, error)
		SetPerms(string, ...string) error
	}
)
