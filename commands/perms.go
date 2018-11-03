package commands

import (
	"context"

	"github.com/justanotherorganization/justanotherbotkit/transport"
	"github.com/pkg/errors"
)

func checkPerms(c *Command, ev *transport.Event) (bool, error) {
	if c.UserDB != nil && len(c.Perms) > 0 {
		return hasPerms(c, ev.GetOrigin().GetSender().GetID())
	}

	return true, nil
}

func hasPerms(c *Command, id string) (bool, error) {
	u, err := c.UserDB.GetUser(context.Background(), id)
	if err != nil {
		return false, errors.Wrap(err, "UserDB.GetUser")
	}

	for _, p := range u.GetPermissions() {
		// Root users can do all the things!!!
		if p == "root" {
			return true, nil
		}

		for _, _p := range c.Perms {
			if p == _p {
				return true, nil
			}
		}
	}

	return false, nil
}
