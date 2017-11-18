package accessors

import (
	"errors"
	"fmt"
	"sync"

	sl "github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
)

type slack struct {
	token string
	log   *logrus.Entry
	api   *sl.Client
	rtm   *sl.RTM
	wg    *sync.WaitGroup
}

// NewSlack returns a new Accessor for the Slack network.
func NewSlack(token string, log *logrus.Entry, wg *sync.WaitGroup) (Accessor, error) {
	if token == "" || len(token) <= 0 {
		return nil, errors.New("slack token cannot be empty")
	}

	if log == nil {
		log = logrus.NewEntry(logrus.New())
	}

	if wg == nil {
		wg = &sync.WaitGroup{}
	}

	return &slack{
		log:   log,
		token: token,
		wg:    wg,
	}, nil
}

func (s *slack) TunnelEvents(eventCh chan MessageEvent, errCh, stopCh chan error) {
	s.api = sl.New(s.token)
	s.rtm = s.api.NewRTM()

	go func() {
		s.wg.Add(1)
		s.rtm.ManageConnection()
		s.wg.Done()
	}()

out:
	for {
		select {
		case <-stopCh:
			break out
		case msg := <-s.rtm.IncomingEvents:
			// TODO: revisit the necessity behind this convoluted method.
			func() {
				defer func() {
					if err := recover(); err != nil {
						errCh <- fmt.Errorf("%v", err)
					}
				}()

				switch event := msg.Data.(type) {
				case *sl.MessageEvent:
					eventCh <- MessageEvent{
						Origin: event.Channel,
						Sender: &MessageEvent_Sender{
							Name: s.rtm.GetInfo().GetUserByID(event.User).Name,
							Id:   event.User,
						},

						Body: event.Msg.Text,
					}
				case *sl.MessageTooLongEvent:
					// TODO: properly handle this...
					s.log.Errorf("%v", event)
				case *sl.ReconnectUrlEvent:
					// TODO: maybe reconnect sometimes?
					s.log.Debugf("%v", event)
				case *sl.DisconnectedEvent:
					// TODO: handle reconnecting...
					s.log.Debugf("%v", event)
				default:
					s.log.Debugf("%v", event)
				}
			}()
		}
	}

	if err := s.rtm.Disconnect(); err != nil {
		errCh <- err
	}
}

func (s *slack) SendMessage(msg, dest string) error {
	s.rtm.SendMessage(s.rtm.NewOutgoingMessage(msg, dest))
	return nil
}
