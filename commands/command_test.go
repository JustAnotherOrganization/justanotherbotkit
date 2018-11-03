package commands_test

import (
	"testing"

	. "github.com/justanotherorganization/justanotherbotkit/commands"
	"github.com/justanotherorganization/justanotherbotkit/internal/test"
	"github.com/justanotherorganization/justanotherbotkit/transport"
	"github.com/justanotherorganization/justanotherbotkit/transport/proto"
)

func TestSingleCommand(t *testing.T) {
	tt := test.TableTest{
		Cases: []*test.TableTestCase{
			// Disabled commands should return nil.
			// TODO: test that nothing was sent over the transport in response.
			&test.TableTestCase{
				Val: &Command{Disabled: true},
				Exp: nil,
			},
			// Perms are set but no database is present (perms will be un-checked).
			// TODO: confirm that the command sent someting over the transport in response.
			&test.TableTestCase{
				Val: &Command{Perms: []string{"foo"}},
				Exp: nil,
			},
		},
		F: func(v interface{}) interface{} {
			cmd := v.(*Command)
			return cmd.Execute(&transport.Event{
				BaseEvent: &pb.BaseEvent{
					Body: "test",
				},
			})
		},
	}

	tt.Run(t)
}
