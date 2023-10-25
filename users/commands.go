package users

import (
	"go.uber.org/zap"
	"justanother.org/justanotherbotkit/commands"
	"justanother.org/justanotherbotkit/transport"
)

// CreateUser
// GetUser
// AlterUserPerms
// DeleteUser
// ListUsers

var (
	createUserPerms = []string{`create-user`}
)

func RegisterUserCommands(root *commands.Command, logger *zap.Logger) {
	if root.UserDB() == nil {
		panic("root command must have a user database set before calling RegisterUserCommands")
	}

	root.AddCommand(createUser(logger))
}

func createUser(logger *zap.Logger) *commands.Command {
	return &commands.Command{
		Use:  "create-user",
		Long: "Create a new privileged user",
		ExecFunc: func(cmd *commands.Command, ev transport.Event) error {
			// create-user grim

			logger.Debug("create-user called", zap.String("ev.Body", ev.Body))

			//_, err := cmd.UserDB().CreateUser(cmd.Context(), repo.User{
			//	// FIXME: figure out how to look up user data
			//	//Name:
			//	//NetworkID:
			//})
			//return err

			return nil
		},
		Perms:  createUserPerms,
		Hidden: true,
	}
}
