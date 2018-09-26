package discord

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/justanotherorganization/justanotherbotkit/transport"
	"github.com/justanotherorganization/justanotherbotkit/transport/proto"
)

type (
	// Discord provides io access to the Discord network.
	Discord struct {
		session *discordgo.Session
	}
)

// Static type checking.
var _ transport.Transport = &Discord{}

// New returns a new instance of Discord.
func New(token string) (*Discord, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &Discord{
		session: session,
	}, nil
}

// TunnelEvents translates discord events to transport.Events and passes them back up the stack.
func (d *Discord) TunnelEvents(ctx context.Context, evCh chan *transport.Event, errCh chan error) {
	d.session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore messages from self.
		if m.Author.ID == s.State.User.ID {
			return
		}

		evCh <- &transport.Event{
			BaseEvent: &pb.BaseEvent{
				Origin: &pb.BaseEvent_Origin{
					Sender: &pb.BaseEvent_Origin_Sender{
						ID:   m.Author.ID,
						Name: m.Author.Username,
					},
					ID: m.ChannelID,
				},
				Body: m.Content,
			},
			Transport: d,
		}
	})

	if err := d.session.Open(); err != nil {
		errCh <- err
		return
	}

	// Block until context is finished.
	<-ctx.Done()

	if err := d.session.Close(); err != nil {
		errCh <- err
	}
}

// SendMessage sends a message over the discord network.
func (d *Discord) SendMessage(dest, msg string) error {
	_, err := d.session.ChannelMessageSend(dest, msg)
	return err
}
