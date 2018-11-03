package commands

import (
	"strings"

	"github.com/justanotherorganization/justanotherbotkit/transport"
)

func (c *Command) match(ev *transport.Event) *Command {
	if ev.GetBody() == "" {
		return c
	}

	fields := strings.Fields(ev.GetBody())
	if len(fields) == 0 {
		return c
	}

	var cmd *Command
	for _, _c := range c.children {
		if isMatch(_c, fields[0]) {
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

func isMatch(c *Command, s string) bool {
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
