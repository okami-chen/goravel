package facades

import (
	"log"

	"goravel/packages/jwt"
	"goravel/packages/jwt/contracts"
)

func Jwt() contracts.Jwt {
	instance, err := jwt.App.Make(jwt.Binding)
	if err != nil {
		log.Println(err)
		return nil
	}

	return instance.(contracts.Jwt)
}
