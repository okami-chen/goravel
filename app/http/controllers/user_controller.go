package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"goravel/packages/jwt/facades"
)

type UserController struct {
	//Dependent services
}

func NewUserController() *UserController {
	return &UserController{
		//Inject services
	}
}

func (r *UserController) Index(ctx http.Context) http.Response {
	facades.Jwt().LoginUsingID(ctx, 1, make(map[string]interface{}))
	return ctx.Response().Success().Json(http.Json{
		"success": true,
		"code":    0,
		"data":    nil,
	})
}
