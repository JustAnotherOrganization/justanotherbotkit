package discord // package github.com/justanotherorganization/justanotherbotkit/transport/discord

import (
	"context"
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/justanotherorganization/justanotherbotkit/transport"
	"github.com/justanotherorganization/justanotherbotkit/transport/proto"
)

type (
	// Discord provides io access to the Discord network.
	Discord struct {
		session     *discordgo.Session
		ignoreUsers []string
	}
)

// Static type checking.
var _ transport.Transport = &Discord{}

// New returns a new instance of Discord.
func New(token string, ignoreUsers ...string) (*Discord, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &Discord{
		session:     session,
		ignoreUsers: ignoreUsers,
	}, nil
}

// TunnelEvents translates discord events to transport.Events tunneling them into evCh.
// The session is terminated when ctx.Done returns.
func (d *Discord) TunnelEvents(ctx context.Context, evCh chan *transport.Event, errCh chan error) {
	d.session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore messages from self.
		if m.Author.ID == s.State.User.ID {
			return
		}

		for _, user := range d.ignoreUsers {
			if m.Author.ID == user ||
				m.Author.Username == user {
				return
			}
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

// SendEvent sends a new event to Slack.
func (d *Discord) SendEvent(ev *transport.Event) error {
	// TODO:
	return errors.New("not yet implemented")
}

// Channels lists all the channels we have access to.
func (d *Discord) Channels() ([]*transport.Channel, error) {
	// TODO:
	return nil, errors.New("not  yet implemented")
}

// GetUser returns the full user data for the provided name or ID.
func (d *Discord) GetUser(user string) (*transport.User, error) {
	// TODO:
	return nil, errors.New("not  yet implemented")
}
