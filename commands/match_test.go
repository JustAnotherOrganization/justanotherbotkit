package commands

import (
	"testing"

	"github.com/justanotherorganization/justanotherbotkit/internal/test"
	"github.com/justanotherorganization/justanotherbotkit/transport"
)

func TestMatch(t *testing.T) {
	type _mt struct {
		cmd *Command
		ev  *transport.Event
	}

	tt := test.TableTest{
		Cases: []*test.TableTestCase{
			&test.TableTestCase{
				Val: &_mt{cmd: &Command{}, ev: &transport.Event{}},
				Exp: &Command{},
			},
		},
		F: func(v interface{}) interface{} {
			mt := v.(*_mt)
			return mt.cmd.match(mt.ev)
		},
	}

	tt.Run(t)
}

func TestIsMatch(t *testing.T) {
	type _mt struct {
		cmd *Command
		s   string
	}

	tt := test.TableTest{
		Cases: []*test.TableTestCase{
			&test.TableTestCase{
				Val: &_mt{cmd: &Command{Use: ""}, s: ""},
				Exp: true,
			},
			&test.TableTestCase{
				Val: &_mt{cmd: &Command{Use: "foo"}, s: ""},
				Exp: false,
			},
			&test.TableTestCase{
				Val: &_mt{cmd: &Command{Use: "foo"}, s: "bar"},
				Exp: false,
			},
			&test.TableTestCase{
				Val: &_mt{cmd: &Command{Use: "foo"}, s: "foo"},
				Exp: true,
			},
			&test.TableTestCase{
				Val: &_mt{cmd: &Command{Use: "foo", Aliases: []string{"bar"}}, s: "bar"},
				Exp: true,
			},
			&test.TableTestCase{
				Val: &_mt{cmd: &Command{Use: "foo", Aliases: []string{"bar", "foobar"}}, s: "bar"},
				Exp: true,
			},
			&test.TableTestCase{
				Val: &_mt{cmd: &Command{Use: "foo", Aliases: []string{"bar", "foobar"}}, s: "foobar"},
				Exp: true,
			},
		},
		F: func(v interface{}) interface{} {
			mt := v.(*_mt)
			return isMatch(mt.cmd, mt.s)
		},
	}

	tt.Run(t)
}
