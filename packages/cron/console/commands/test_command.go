package commands

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
)

type TestCommand struct {
}

// Signature The name and signature of the console command.
func (receiver *TestCommand) Signature() string {
	return "test"
}

// Description The console command description.
func (receiver *TestCommand) Description() string {
	return "Command description"
}

// Extend The console command extend.
func (receiver *TestCommand) Extend() command.Extend {
	return command.Extend{
		Category: "cron",
		Flags: []command.Flag{
			&command.StringFlag{
				Name:    "lang",
				Value:   "default",
				Aliases: []string{"l"},
				Usage:   "language for the greeting",
			},
		},
	}
}

// Handle Execute the console command.
func (receiver *TestCommand) Handle(ctx console.Context) error {

	return nil
}
