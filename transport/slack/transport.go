package slack // package github.com/justanotherorganization/justanotherbotkit/transpot/slack

import (
	"context"
	"errors"
	"fmt"

	bkPb "github.com/justanotherorganization/justanotherbotkit/proto"
	"github.com/justanotherorganization/justanotherbotkit/transport"
	"github.com/justanotherorganization/justanotherbotkit/transport/proto"
	"github.com/nlopes/slack"
)

// Slack provides io access to the Slack network.
type Slack struct {
	client *slack.Client
	rtm    *slack.RTM
	cfg    *transport.Config
}

// Static type checking.
var _ transport.Transport = &Slack{}

// New returns a new instance of Slack.
// cfg.IgnoreUsers must be set with the bot name or ID otherwise it will potentially read
// it's own messages.
func New(cfg *transport.Config) (*Slack, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	_slack := &Slack{
		client: slack.New(cfg.Token),
		cfg:    cfg,
	}

	// Rebuild the ignore user list with IDs only.
	finalIgnore := make([]string, len(cfg.IgnoreUsers), len(cfg.IgnoreUsers))
	for i, u := range cfg.IgnoreUsers {
		user, err := _slack.GetUser(u)
		if err != nil {
			return nil, err
		}

		finalIgnore[i] = user.GetID()
	}
	cfg.IgnoreUsers = finalIgnore

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
				// TODO: if this BS is necessary we need an option to dump a stacktrace...
				defer func() {
					if err := recover(); err != nil {
						errCh <- fmt.Errorf("%v", err)
					}
				}()

				switch event := msg.Data.(type) {
				case *slack.MessageEvent:
					for _, user := range s.cfg.IgnoreUsers {
						if event.User == user {
							return
						}
					}

					u, err := s.client.GetUserInfo(event.User)
					if err != nil {
						errCh <- errors.New("slack user not known, this should not happen")
						return
					}

					evCh <- &transport.Event{
						BaseEvent: &pb.BaseEvent{
							Origin: &pb.BaseEvent_Origin{
								Sender: &pb.BaseEvent_Origin_Sender{
									ID:   event.User,
									Name: u.Name,
								},
								ID: event.Channel,
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
	if err == nil {
		return s.fromSlackUser(_user), nil
	}

	users, _err := s.GetUsers()
	if _err != nil {
		// Return the original error, knowing why we failed to lookup the one user
		// is probably related to why we failed to lookup the full list anyway.
		// FIXME: log this instead and return the proper error?
		return nil, err
	}

	for _, u := range users {
		if u.GetID() == user ||
			u.GetName() == user {
			return u, nil
		}
	}

	return nil, transport.ErrUserNotFound
}

// GetUsers returns a list of all known users.
func (s *Slack) GetUsers() ([]*transport.User, error) {
	_users, err := s.client.GetUsers()
	if err != nil {
		return nil, err
	}

	users := make([]*transport.User, len(_users), len(_users))
	for i, u := range _users {
		users[i] = s.fromSlackUser(&u)
	}

	return users, nil
}

// GetConversation gets a private conversation for the given user ID.
func (s *Slack) GetConversation(userID string) (string, error) {
	_, _, id, err := s.client.OpenIMChannel(userID)
	return id, err
}

func (s *Slack) fromSlackUser(user *slack.User) *transport.User {
	return &transport.User{
		BaseUser: &bkPb.BaseUser{
			ID:   user.ID,
			Name: user.Name,
		},
		Transport: s,
	}
}
