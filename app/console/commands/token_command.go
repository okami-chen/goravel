package commands

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/http"
)

type TokenCommand struct {
}

// Signature The name and signature of the console command.
func (receiver *TokenCommand) Signature() string {
	return "command:name"
}

// Description The console command description.
func (receiver *TokenCommand) Description() string {
	return "Command description"
}

// Extend The console command extend.
func (receiver *TokenCommand) Extend() command.Extend {
	return command.Extend{}
}

// Handle Execute the console command.
func (receiver *TokenCommand) Handle(ctx console.Context) error {
	token, _ := facades.Auth().LoginUsingID(http.NewContext(), 1)
	facades.Log().Info(token)
	return nil
}
