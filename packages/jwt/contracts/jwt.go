package contracts

import (
	"github.com/goravel/framework/contracts/http"
	"time"
)

type Jwt interface {
	// Guard attempts to get the guard against the local cache.
	Guard(name string) Jwt
	// Parse the given token.
	Parse(ctx http.Context, token string) (*Payload, error)
	// User returns the current authenticated user.
	User(ctx http.Context, user any, expand map[string]interface{}) error
	// Login logs a user into the application.
	Login(ctx http.Context, user any, expand map[string]interface{}) (token string, err error)
	// LoginUsingID logs the given user ID into the application.
	LoginUsingID(ctx http.Context, id any, expand map[string]interface{}) (token string, err error)
	// Refresh the token for the current user.
	Refresh(ctx http.Context) (token string, err error)
	// Logout logs the user out of the application.
	Logout(ctx http.Context) error
}

type Payload struct {
	Guard    string
	Key      string
	ExpireAt time.Time
	IssuedAt time.Time
	Expand   map[string]interface{}
}
