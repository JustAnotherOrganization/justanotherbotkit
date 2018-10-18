// Package commands provides command handling capable of ingesting transport.Events
// and responding over the provided transport wire.
// If the userDB is set permissions on a command will be honored otherwise
// if no userDB is set permissions will be ignored entirely (this could possibly
// use improvement).
package commands // import "github.com/justanotherorganization/justanotherbotkit/commands"

import (
	"context"
	"fmt"
	"strings"

	"github.com/justanotherorganization/justanotherbotkit/transport"
	"github.com/justanotherorganization/justanotherbotkit/users"
	"github.com/pkg/errors"
)

type (
	// Command is the base for all commands.
	Command struct {
		Use     string
		Aliases []string
		Short   string
		Long    string

		ExecFunc func(ev *transport.Event) error

		Perms    []string
		Disabled bool
		Hidden   bool

		UserDB   users.DB
		parent   *Command
		children []*Command
	}
)

// Execute executes the command using the provided event.
// It attempts to match the event against any children commands first
// allowing a base Command to act as a command tree.
func (c *Command) Execute(ev *transport.Event) error {
	if c.Disabled {
		return nil
	}

	if c.UserDB != nil && len(c.Perms) > 0 {
		ok, err := _hasPerms(c, ev.Origin.Sender.ID)
		if err != nil {
			return err
		}

		if !ok {
			return nil
		}
	}

	fields := strings.Fields(ev.Body)

	if strings.Compare("help", fields[0]) == 0 {
		ev.Body = strings.Join(fields[1:], " ")
		return c.help(ev)
	}

	_c := c.match(ev)
	if _c == nil {
		return nil
	}

	if _c.UserDB != nil && len(_c.Perms) > 0 {
		ok, err := _hasPerms(c, ev.Origin.Sender.ID)
		if err != nil {
			return err
		}

		if !ok {
			return nil
		}
	}

	if _c.ExecFunc == nil {
		return nil
	}

	return _c.ExecFunc(ev)
}

// AddCommand adds a new Command to the current command as a child.
func (c *Command) AddCommand(cmd *Command) {
	cmd.parent = c

	if c.UserDB != nil {
		cmd.UserDB = c.UserDB
	}

	c.children = append(c.children, cmd)
}

// Parent returns the parent command if there is one.
func (c *Command) Parent() *Command {
	return c.parent
}

// Children returns the children commands if there are any.
func (c *Command) Children() []*Command {
	return c.children
}

func (c *Command) help(ev *transport.Event) error {
	// FIXME: this should return proper help from the root command
	// IE: it should show children commands.
	_c := c.match(ev)
	if _c == nil {
		return nil
	}

	return ev.SendMessage(
		ev.Origin.ID,
		fmt.Sprintf("%s:\n%s\n", _c.Use, _c.Long),
	)
}

func (c *Command) match(ev *transport.Event) *Command {
	if ev.Body == "" {
		return c
	}

	fields := strings.Fields(ev.Body)
	if len(fields) == 0 {
		return c
	}

	var cmd *Command
	for _, _c := range c.children {
		if _isCommand(_c, fields[0]) {
			ev.Body = strings.Join(fields[1:], " ")
			cmd = _c
			break
		}
	}

	if cmd == nil {
		cmd = c
	}

	return cmd
}

func _isCommand(c *Command, s string) bool {
	if strings.Compare(c.Use, s) == 0 {
		return true
	}

	for _, a := range c.Aliases {
		if strings.Compare(a, s) == 0 {
			return true
		}
	}

	return false
}

func _hasPerms(c *Command, id string) (bool, error) {
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
