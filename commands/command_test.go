package commands_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"justanother.org/justanotherbotkit/commands"
	"justanother.org/justanotherbotkit/transport"
	"justanother.org/justanotherbotkit/transport/impl/mock"
	"justanother.org/justanotherbotkit/transport/pkg/option"
)

func TestCommand_Execute(t *testing.T) {
	tr := mock.NewMockTransport(gomock.NewController(t))

	var called bool
	root := &commands.Command{}
	root.AddCommand(&commands.Command{
		Use:     "greet",
		Aliases: []string{"hello"},
		Long:    "A command which greets someone",
		ExecFunc: func(ev transport.Event) error {
			assert.Equal(t, "", ev.Body)
			called = true
			return nil
		},
	})
	root.AddCommand(&commands.Command{
		Use:     "dnc",
		Aliases: []string{"do not call"},
		Long:    "A command that should never be called",
		ExecFunc: func(ev transport.Event) error {
			return errors.New("should not have been called")
		},
	})

	ev := transport.Event{
		Origin: transport.EventOrigin{
			ID: "foo",
			Sender: transport.EventOriginSender{
				ID:   "foobar",
				Name: "baz",
			},
		},
		Body:      "greet",
		Transport: tr,
	}

	called = false
	err := root.Execute(ev)
	require.NoError(t, err)
	assert.True(t, called)

	called = false
	ev.Body = "hello"
	err = root.Execute(ev)
	require.NoError(t, err)
	assert.True(t, called)

	ev.Body = "dnc"
	err = root.Execute(ev)
	require.Error(t, err)
}

type optionMatcher struct {
	options []transport.MsgOption
}

func (m optionMatcher) Matches(x interface{}) bool {
	return cmp.Equal(m.options, x)
}

func (m optionMatcher) String() string {
	return fmt.Sprintf("%v", m.options)
}

func TestCommand_Help(t *testing.T) {
	tr := mock.NewMockTransport(gomock.NewController(t))

	root := &commands.Command{}

	parent := &commands.Command{
		Use:   "parent",
		Short: "A parent command",
		Long:  "A parent command",
	}

	parent.AddCommand(&commands.Command{
		Use:   "child",
		Short: "A child command",
		Long:  "A child command",
	})
	root.AddCommand(parent)

	root.AddCommand(&commands.Command{
		Use:   "alt",
		Short: "An alternative command",
		Long:  "An alternative command",
	})

	t.Run("help", func(tt *testing.T) {
		tr.EXPECT().
			SendMessage(gomock.Any(), gomock.Eq("channel-id"), optionMatcher{
				options: []transport.MsgOption{
					option.Text{
						Value:  "sub commands:\n\tparent:\tA parent command\n\talt:\tAn alternative command\n",
						Escape: false,
					},
				},
			})

		err := root.Execute(transport.Event{
			Origin: transport.EventOrigin{
				ID: "channel-id",
				Sender: transport.EventOriginSender{
					ID:   "foobar",
					Name: "baz",
				},
			},
			Body:      "help",
			Transport: tr,
		})
		require.NoError(tt, err)
	})

	t.Run("help parent", func(tt *testing.T) {
		tr.EXPECT().
			SendMessage(gomock.Any(), gomock.Eq("channel-id"), optionMatcher{
				options: []transport.MsgOption{
					option.Text{
						Value:  "parent\nA parent command\n\nsub commands:\n\tchild:\tA child command\n",
						Escape: false,
					},
				},
			})

		err := root.Execute(transport.Event{
			Origin: transport.EventOrigin{
				ID: "channel-id",
				Sender: transport.EventOriginSender{
					ID:   "foobar",
					Name: "baz",
				},
			},
			Body:      "help parent",
			Transport: tr,
		})
		require.NoError(tt, err)
	})

	t.Run("help parent child", func(tt *testing.T) {
		tr.EXPECT().
			SendMessage(gomock.Any(), gomock.Eq("channel-id"), optionMatcher{
				options: []transport.MsgOption{
					option.Text{
						Value:  "parent child\nA child command\n",
						Escape: false,
					},
				},
			})

		err := root.Execute(transport.Event{
			Origin: transport.EventOrigin{
				ID: "channel-id",
				Sender: transport.EventOriginSender{
					ID:   "foobar",
					Name: "baz",
				},
			},
			Body:      "help parent child",
			Transport: tr,
		})
		require.NoError(tt, err)
	})
}
