package commands

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
	"github.com/goravel/gin"
)

type TokenCommand struct {
}

// Signature The name and signature of the console command.
func (receiver *TokenCommand) Signature() string {
	return "jwt:token"
}

// Description The console command description.
func (receiver *TokenCommand) Description() string {
	return "Generate a JWT token"
}

// Extend The console command extend.
func (receiver *TokenCommand) Extend() command.Extend {
	return command.Extend{
		Category: "jwt",
	}
}

// Handle Execute the console command.
func (receiver *TokenCommand) Handle(ctx console.Context) error {
	token, _ := facades.Auth().LoginUsingID(gin.Background(), 1)
	facades.Log().Infof("Token: %s", token)
	return nil
}
