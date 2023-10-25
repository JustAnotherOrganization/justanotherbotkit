package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"justanother.org/justanotherbotkit/commands"
	"justanother.org/justanotherbotkit/transport"
	"justanother.org/justanotherbotkit/transport/impl/discord"
	"justanother.org/justanotherbotkit/users"
	"justanother.org/justanotherbotkit/users/repo/postgres"
)

func main() {
	// Suggest attaching signal handling to context.
	if err := start(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func start(ctx context.Context) error {
	ts, err := discord.New(transport.Config{
		Token: os.Getenv(`EXAMPLE_DISCORD_TOKEN`),
	})
	if err != nil {
		return err
	}

	pool, err := pgxpool.New(ctx, postgres.DSN())
	if err != nil {
		return err
	}

	root := new(commands.Command).
		WithUserDB(postgres.New(pool))

	ts.MessageEventHandler(root.Execute, func(err error) {
		log.Println(err.Error())
	})

	_log, err := zap.NewDevelopment()
	if err != nil {
		return err
	}

	users.RegisterUserCommands(root, _log)

	return ts.Start(ctx)
}
