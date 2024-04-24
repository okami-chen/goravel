package facades

import (
	"log"

	"goravel/packages/clash"
	"goravel/packages/clash/contracts"
)

func Clash() contracts.Clash {
	instance, err := clash.App.Make(clash.Binding)
	if err != nil {
		log.Println(err)
		return nil
	}

	return instance.(contracts.Clash)
}
