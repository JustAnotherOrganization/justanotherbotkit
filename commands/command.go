package commands

import (
	"fmt"
	"strings"

	"github.com/justanotherorganization/justanotherbotkit/transport"
	"github.com/justanotherorganization/justanotherbotkit/users"
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
	if ev == nil ||
		ev.GetBody() == "" ||
		c.Disabled {
		return nil
	}

	ok, err := checkPerms(c, ev)
	if err != nil {
		return err
	}

	if !ok {
		return nil
	}

	fields := strings.Fields(ev.GetBody())

	if strings.Compare("help", fields[0]) == 0 {
		ev.Body = strings.Join(fields[1:], " ")
		return c.help(ev)
	}

	_c := c.match(ev)
	if _c == nil ||
		_c.Disabled {
		return nil
	}

	ok, err = checkPerms(_c, ev)
	if err != nil {
		return err
	}

	if !ok {
		return nil
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
		ev.GetOrigin().GetID(),
		fmt.Sprintf("%s:\n%s\n", _c.Use, _c.Long),
	)
}
