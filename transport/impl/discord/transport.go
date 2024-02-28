package discord

import (
	"context"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"justanother.org/justanotherbotkit/transport"
	"justanother.org/justanotherbotkit/transport/pkg/option"
)

type Transport struct {
	Session *discordgo.Session
	Config  transport.Config
}

func New(cfg transport.Config) (Transport, error) {
	t := Transport{
		Config: cfg,
	}

	var err error
	t.Session, err = discordgo.New("Bot " + cfg.Token)
	return t, err
}

func (t Transport) MessageEventHandler(h func(ctx context.Context, ev transport.Event) error, errHandler func(err error)) {
	t.Session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		for _, user := range t.Config.IgnoreUsers {
			if m.Author.ID == user ||
				m.Author.Username == user {
				return
			}
		}

		// FIXME: get the context somewhere passed into the transport.
		err := h(context.Background(), transport.Event{
			Origin: transport.EventOrigin{
				ID: m.ChannelID,
				Sender: transport.EventOriginSender{
					ID:   m.Author.ID,
					Name: m.Author.Username,
				},
			},
			Body:      m.Content,
			Transport: t,
		})
		if err != nil {
			errHandler(err)
		}
	})
}

func (t Transport) SendMessage(_ context.Context, dest string, options ...transport.MsgOption) error {
	var (
		msg            string
		requestOptions []discordgo.RequestOption
	)
	for _, opt := range options {
		switch o := opt.(type) {
		case option.Text:
			msg = o.Value
		case discordgo.RequestOption:
			requestOptions = append(requestOptions, o)
		default:
			return fmt.Errorf("unsupported option: %T", o)
		}
	}

	if msg == "" {
		return errors.New("expected Text option to be included")
	}

	_, err := t.Session.ChannelMessageSend(dest, msg)
	return err
}
