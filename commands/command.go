package commands

import (
	"context"
	"fmt"
	"strings"

	"justanother.org/justanotherbotkit/transport"
	"justanother.org/justanotherbotkit/transport/pkg/option"
	"justanother.org/justanotherbotkit/users/repo"
)

type (
	Command struct {
		Use     string
		Aliases []string
		Short   string
		Long    string

		ExecFunc func(cmd *Command, event transport.Event) error

		Perms    []string
		Disabled bool
		Hidden   bool

		parent   *Command
		children []*Command

		ctx    context.Context
		userDB repo.Repo
	}
)

func (c *Command) Execute(ctx context.Context, ev transport.Event) error {
	if c.Disabled {
		return nil
	}

	if ev.Body == "" {
		return nil
	}

	c.ctx = ctx
	defer func() {
		c.ctx = nil
	}()

	if c.UserDB() != nil && len(c.Perms) > 0 {
		ok, err := _hasPerms(c, ev.Origin.Sender.ID)
		if err != nil {
			return err
		}

		if !ok {
			// FIXME: consider returning an error instead of ignoring the command.
			return nil
		}
	}

	fields := strings.Fields(ev.Body)

	if strings.Compare("help", fields[0]) == 0 {
		ev.Body = strings.Join(fields[1:], " ")
		return c.help(ev)
	}

	_c := c.match(&ev)
	if _c == nil {
		return nil
	}

	if _c.UserDB() != nil && len(_c.Perms) > 0 {
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

	return _c.ExecFunc(_c, ev)
}

func (c *Command) AddCommand(cmd *Command) {
	cmd.parent = c
	c.children = append(c.children, cmd)
}

func (c *Command) Root() *Command {
	r := c
	for {
		if r.parent == nil {
			return r
		}

		r = r.parent
	}
}

func (c *Command) WithUserDB(db repo.Repo) *Command {
	c.Root().userDB = db
	return c
}

func (c *Command) UserDB() repo.Repo {
	return c.Root().userDB
}

func (c *Command) Context() context.Context {
	return c.ctx
}

func (c *Command) help(ev transport.Event) error {
	// FIXME: this should return proper help from the root command
	// IE: it should show children commands.
	if ev.Body == "" {
		var b strings.Builder
		if c.Use != "" {
			b.WriteString(c.name() + "\n")
		}
		if c.Long != "" {
			b.WriteString(c.Long + "\n")
		} else if c.Short != "" {
			b.WriteString(c.Short + "\n")
		}

		c.helpForChildren(&b)
		return ev.SendMessage(
			c.Context(),
			ev.Origin.ID,
			option.Text{
				Value: b.String(),
			},
		)
	}

	_c := c.match(&ev)
	if _c == nil {
		return nil
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s\n%s\n", _c.name(), _c.Long))
	_c.helpForChildren(&b)

	return ev.SendMessage(
		c.Context(),
		ev.Origin.ID,
		option.Text{
			Value: b.String(),
		},
	)
}

func (c *Command) helpForChildren(b *strings.Builder) {
	if len(c.children) > 0 {
		if b.Len() > 0 {
			b.WriteString("\n")
		}

		b.WriteString("sub commands:\n")

		for _, _c := range c.children {
			desc := _c.Short
			if desc == "" {
				desc = _c.Long
			}

			b.WriteString(fmt.Sprintf("\t%s:\t%s\n", _c.Use, desc))
		}
	}
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
	cmd, ev.Body = c.matchChildren(fields...)

	if cmd == nil {
		cmd = c
	}

	return cmd
}

func (c *Command) matchChildren(
	fields ...string,
) (cmd *Command, body string) {
	for _, _c := range c.children {
		if _isCommand(_c, fields[0]) {
			if len(_c.children) > 0 && len(fields) > 1 {
				return _c.matchChildren(fields[1:]...)
			}

			body = strings.Join(fields[1:], " ")
			cmd = _c
			break
		}
	}

	return cmd, body
}

func (c *Command) name() string {
	var b strings.Builder

	if c.parent != nil {
		p := c.parent
		for {
			if p == nil || p.Use == "" {
				break
			}

			b.WriteString(p.Use + " ")

			p = p.parent
		}
	}

	b.WriteString(c.Use)

	return b.String()
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

func _hasPerms(c *Command, networkID string) (bool, error) {
	u, err := c.UserDB().GetUserByNetworkID(c.Context(), networkID)
	if err != nil {
		return false, fmt.Errorf("UserDB.GetUserByNetworkID, %w", err)
	}

	for _, p := range u.Permissions {
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
