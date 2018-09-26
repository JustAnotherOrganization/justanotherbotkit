package slack

import (
	"context"
	"fmt"

	"github.com/justanotherorganization/justanotherbotkit/transport"
	"github.com/justanotherorganization/justanotherbotkit/transport/proto"
	"github.com/nlopes/slack"
)

type (
	// Slack provides io access to the Slack network.
	Slack struct {
		client *slack.Client
		rtm    *slack.RTM
	}
)

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

// TunnelEvents translates slack events to transport.Events and passes them back up the stack.
func (s *Slack) TunnelEvents(ctx context.Context, evCh chan *transport.Event, errCh chan error) {
	go func() {
		s.rtm.ManageConnection()
	}()

	for finished := false; !finished; {
		select {
		case <-ctx.Done():
			finished = true
		case msg := <-s.rtm.IncomingEvents:
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

// SendMessage sends a message over the slack network.
func (s *Slack) SendMessage(dest, msg string) error {
	s.client.SendMessage(
		dest,
		slack.MsgOptionText(msg, false),
		slack.MsgOptionAsUser(true), // False to send messages as slackbot.
	)
	return nil
}
