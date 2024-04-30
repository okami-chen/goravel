package commands

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
)

type TokenCommand struct {
}

//Signature The name and signature of the console command.
func (receiver *TokenCommand) Signature() string {
	return "command:name"
}

//Description The console command description.
func (receiver *TokenCommand) Description() string {
	return "Command description"
}

//Extend The console command extend.
func (receiver *TokenCommand) Extend() command.Extend {
	return command.Extend{}
}

//Handle Execute the console command.
func (receiver *TokenCommand) Handle(ctx console.Context) error {
	
	return nil
}
