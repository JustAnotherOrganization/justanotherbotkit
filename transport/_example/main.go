package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/justanotherorganization/justanotherbotkit/transport"
	discord "github.com/justanotherorganization/justanotherbotkit/transport/discord"
	slack "github.com/justanotherorganization/justanotherbotkit/transport/slack"
)

var (
	discordToken string
	slackToken   string
)

func init() {
	flag.StringVar(&discordToken, "dt", "", "Discord token")
	flag.StringVar(&slackToken, "st", "", "Slack token")
	flag.Parse()

	if discordToken == "" {
		discordToken = os.Getenv("DISCORDTOKEN")
	}

	if slackToken == "" {
		slackToken = os.Getenv("SLACKTOKEN")
	}

	if discordToken == "" && slackToken == "" {
		fmt.Println("slack or discord token must be set")
		os.Exit(1)
	}
}

func main() {
	evCh := make(chan *transport.Event)
	errCh := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())

	if discordToken != "" {
		d, err := discord.New(discordToken)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		go d.TunnelEvents(ctx, evCh, errCh)
	}

	if slackToken != "" {
		s, err := slack.New(slackToken)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		go s.TunnelEvents(ctx, evCh, errCh)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	for finished := false; !finished; {
		select {
		case <-sc:
			finished = true
		case ev := <-evCh:
			fmt.Fprintf(os.Stdout, "%v\n", ev)
		case err := <-errCh:
			fmt.Fprintln(os.Stderr, err.Error())
			// non-fatal, for now
		}
	}

	cancel()
}
