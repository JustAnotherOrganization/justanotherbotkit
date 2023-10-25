package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	"justanother.org/justanotherbotkit/transport"
	"justanother.org/justanotherbotkit/transport/impl/discord"
	"justanother.org/justanotherbotkit/transport/impl/discord/interactions"
)

func main() {
	if err := start(); err != nil {
		log.Fatal(err)
	}
}

func start() error {
	_log, err := zap.NewDevelopment()
	if err != nil {
		return err
	}

	ts, err := discord.New(transport.Config{
		Token: os.Getenv(`EXAMPLE_DISCORD_TOKEN`),
	})
	if err != nil {
		return err
	}

	interactionController := interactions.New(interactions.Config{
		Log:            _log,
		Transport:      ts,
		RemoveCommands: true, // Remove our command/s when done.
	})
	defer func() {
		_ = interactionController.Close()
	}()

	if err = interactionController.Register(hello(_log)); err != nil {
		return err
	}

	return ts.Start(context.Background()) // Suggest attaching signal handling to the context.
}

func hello(log *zap.Logger) func() (*discordgo.ApplicationCommand, func(session *discordgo.Session, ic *discordgo.InteractionCreate)) {
	return func() (*discordgo.ApplicationCommand, func(session *discordgo.Session, ic *discordgo.InteractionCreate)) {
		return &discordgo.ApplicationCommand{
				Name:        "hello",
				Description: "hello <foo>",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "param",
						Description: "Hello what?",
						Type:        discordgo.ApplicationCommandOptionString,
					},
				},
			},
			func(session *discordgo.Session, ic *discordgo.InteractionCreate) {
				if len(ic.ApplicationCommandData().Options) != 1 {
					log.Error("malformed InteractionCreate event",
						zap.Any("ic", ic))
					return
				}

				paramRaw := ic.ApplicationCommandData().Options[0].Value
				param, ok := paramRaw.(string)
				if !ok {
					// Consider passing these sorts of errors back in the response.
					log.Error(fmt.Sprintf("malformed InteractionCreate Option value, expected string, got %T", paramRaw))
					return
				}

				if err := session.InteractionRespond(ic.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Hello %s!", param),
					},
				}); err != nil {
					log.Error("error sending response to interaction command", zap.Error(err))
					return
				}
			}
	}
}
