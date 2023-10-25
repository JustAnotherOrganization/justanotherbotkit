package interactions

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	"justanother.org/justanotherbotkit/transport/impl/discord"
)

type (
	Config struct {
		Log            *zap.Logger
		Transport      discord.Transport
		RemoveCommands bool
	}

	Controller struct {
		Config          Config
		mux             sync.RWMutex
		commandHandlers map[string]command
	}

	command struct {
		id      string
		handler func(session *discordgo.Session, ic *discordgo.InteractionCreate)
	}
)

func New(cfg Config) *Controller {
	c := &Controller{
		Config:          cfg,
		commandHandlers: make(map[string]command),
	}

	c.Config.Transport.Session.AddHandler(c.Handler)

	return c
}

func (c *Controller) Register(f func() (*discordgo.ApplicationCommand, func(session *discordgo.Session, ic *discordgo.InteractionCreate))) error {
	cmd, handler := f()
	name := cmd.Name
	cmd, err := c.Config.Transport.Session.ApplicationCommandCreate(c.Config.Transport.Session.State.User.ID, "", cmd)
	if err != nil {
		return fmt.Errorf("s.ApplicationCommandCreate(%s), %w", name, err)
	}

	c.mux.Lock()
	c.commandHandlers[name] = command{
		id:      cmd.ID,
		handler: handler,
	}
	c.mux.Unlock()

	return nil
}

func (c *Controller) Close() error {
	if c.Config.RemoveCommands {
		c.mux.Lock()
		defer c.mux.Unlock()

		for name, cmd := range c.commandHandlers {
			err := c.Config.Transport.Session.ApplicationCommandDelete(c.Config.Transport.Session.State.User.ID, "", cmd.id)
			if err != nil {
				return fmt.Errorf("error removing command, %s: %w", name, err)
			}

			c.Config.Log.
				Debug("removed command", zap.String("name", name))
		}
	}

	return nil
}

func (c *Controller) Handler(session *discordgo.Session, ic *discordgo.InteractionCreate) {
	c.mux.RLock()
	name := ic.ApplicationCommandData().Name
	if cmd, ok := c.commandHandlers[name]; ok {
		c.Config.Log.
			Debug("calling command handler for "+name, zap.Any("interaction create", ic))

		cmd.handler(session, ic)
	}
	c.mux.RUnlock()
}
