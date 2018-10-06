package slack // package github.com/justanotherorganization/justanotherbotkit/transpot/slack

import (
	"context"
	"errors"
	"fmt"

	"github.com/justanotherorganization/justanotherbotkit/transport"
	"github.com/justanotherorganization/justanotherbotkit/transport/proto"
	"github.com/nlopes/slack"
)

// Slack provides io access to the Slack network.
type Slack struct {
	client *slack.Client
	rtm    *slack.RTM
}

// Static type checking.
var _ transport.Transport = &Slack{}

// New returns a new instance of Slack.
func New(token string) (*Slack, error) {
	_slack := &Slack{
		client: slack.New(token),
	}

	_slack.rtm = _slack.client.NewRTM()

	return _slack, nil
}

// TunnelEvents translates Slack events to transport.Events tunneling them into evCh.
// The session is terminated when ctx.Done returns.
func (s *Slack) TunnelEvents(ctx context.Context, evCh chan *transport.Event, errCh chan error) {
	go func() {
		s.rtm.ManageConnection()
	}()

	for finished := false; !finished; {
		select {
		case <-ctx.Done():
			finished = true
		case msg := <-s.rtm.IncomingEvents:
			// FIXME: this weird but slack has a habit of crashing...
			func() {
				defer func() {
					if err := recover(); err != nil {
						errCh <- fmt.Errorf("%v", err)
					}
				}()

				switch event := msg.Data.(type) {
				case *slack.MessageEvent:
					// FIXME: ignore messages from self.

					evCh <- &transport.Event{
						BaseEvent: &pb.BaseEvent{
							Origin: &pb.BaseEvent_Origin{
								Sender: &pb.BaseEvent_Origin_Sender{
									ID:   event.User,
									Name: s.rtm.GetInfo().GetUserByID(event.User).Name,
								},
							},
							Body: event.Msg.Text,
						},
						Transport: s,
					}
				case *slack.MessageTooLongEvent:
					// TODO: properly handle this...
				case *slack.ReconnectUrlEvent:
					// TODO: maybe reconnect sometimes?
				case *slack.DisconnectedEvent:
					// TODO: handle reconnecting...
				default:
				}
			}()
		}
	}

	if err := s.rtm.Disconnect(); err != nil {
		errCh <- err
	}
}

// SendMessage sends a message using the default format.
func (s *Slack) SendMessage(dest, msg string) error {
	s.client.SendMessage(
		dest,
		slack.MsgOptionText(msg, false),
		// FIXME: this should be part of the transport configuration.
		slack.MsgOptionAsUser(true), // False to send messages as slackbot.
	)
	return nil
}

// SendEvent sends a new event to Slack.
func (s *Slack) SendEvent(ev *transport.Event) error {
	// TODO:
	return errors.New("not yet implemented")
}

// Channels lists all the channels we have access to.
func (s *Slack) Channels() ([]*transport.Channel, error) {
	chs, err := s.client.GetChannels(false)
	if err != nil {
		return nil, err
	}

	ret := make([]*transport.Channel, 0, len(chs))
	for _, c := range chs {
		ret = append(ret, &transport.Channel{
			BaseChannel: &pb.BaseChannel{
				ID:        c.ID,
				Name:      c.Name,
				MemberIDs: c.Members,
				Archived:  c.IsArchived,
				// TODO:
				// Topic:
				// Purpose:
			},
		})
	}

	return ret, nil
}

// GetUser returns the full user data for the provided name or ID.
func (s *Slack) GetUser(user string) (*transport.User, error) {
	_user, err := s.client.GetUserInfo(user)
	if err != nil {
		return nil, err
	}

	return &transport.User{
		BaseUser: &pb.BaseUser{
			ID:   _user.ID,
			Name: _user.Name,
		},
		Transport: s,
	}, nil
}
