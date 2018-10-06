package commands

import (
	"fmt"
	"strings"

	"github.com/justanotherorganization/justanotherbotkit/transport"
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

		pm       *_pm // FIXME: replace this with a real permissions manager
		parent   *Command
		children []*Command
	}

	_pm struct{}
)

// Execute executes the command using the provided event.
// It attempts to match the event against any children commands first
// allowing a base Command to act as a command tree.
func (c *Command) Execute(ev *transport.Event) error {
	if c.Disabled {
		return nil
	}

	// TODO: handle perms

	fields := strings.Fields(ev.Body)

	if strings.Compare("help", fields[0]) == 0 {
		ev.Body = strings.Join(fields[1:], " ")
		return c.help(ev)
	}

	_c := c.match(ev)
	if _c == nil {
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

	if cmd.pm == nil {
		cmd.pm = c.pm
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
