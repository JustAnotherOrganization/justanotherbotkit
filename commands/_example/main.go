package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/justanotherorganization/justanotherbotkit/commands"
	"github.com/justanotherorganization/justanotherbotkit/transport"
	"github.com/justanotherorganization/justanotherbotkit/transport/discord"
)

func setupTransport(token string) (transport.Transport, error) {
	return discord.New(&transport.Config{
		Token: token,
	})
}

func setupCommands() *commands.Command {
	// Root command will select appropriate child command based on the message's first word
	root := &commands.Command{}
	// Child commands
	root.AddCommand(&commands.Command{
		Use:      "greet",                          // Command name
		Aliases:  []string{"hello"},                // Command aliases
		Long:     "A command which greets someone", // Command help text
		ExecFunc: greetCommand,                     // Command implementation
	})
	return root
}

func greetCommand(ev *transport.Event) error {
	param := ev.Body
	// Check for missing parameters
	if param == "" {
		ev.SendMessage(ev.Origin.ID, fmt.Sprintf("@%s: not enough arguments", ev.Origin.Sender.Name))
		return nil
	}
	// Reply to command with a greeting
	ev.SendMessage(ev.Origin.ID, fmt.Sprintf("@%s: Hello %s!", ev.Origin.Sender.Name, param))
	return nil
}

func main() {
	// Parse parameters
	var token string
	flag.StringVar(&token, "token", "", "Discord bot token")
	flag.Parse()
	if token == "" {
		token = os.Getenv("DISCORDTOKEN")
	}
	if token == "" {
		fmt.Println("Discord token must be set")
		os.Exit(1)
	}

	// Setup transport and root command
	t, err := setupTransport(token)
	if err != nil {
		log.Fatalf("could not setup transport: %v", err)
	}
	c := setupCommands()

	// Start events processor
	ctx, cancel := context.WithCancel(context.Background())
	events := make(chan *transport.Event)
	errors := make(chan error)
	go t.TunnelEvents(ctx, events, errors)

	// Start signal handler
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	log.Println("Bot is now running. Press Ctrl+C to exit.")

	// Handle incoming events
loop:
	for {
		select {
		// Exit on SIGINT
		case <-signals:
			log.Println("Exiting...")
			cancel()
			break loop
		// Execute commands in incoming events
		case ev := <-events:
			err := c.Execute(ev)
			if err != nil {
				log.Printf("could not execute command: %v", err)
			}
		// Print event processing errors
		case err := <-errors:
			if err != nil {
				log.Printf("could not receive event: %v", err)
			}
		}
	}
}
