package api

import (
	"github.com/labstack/echo/v4"

	auth "github.com/kl09/auth-go"
)

type Router struct {
	credService auth.CredentialService
}

func NewRouter(credService auth.CredentialService) *Router {
	return &Router{
		credService: credService,
	}
}

func (r *Router) Handler() *echo.Echo {
	e := echo.New()
	e.HTTPErrorHandler = customHTTPErrorHandler

	e.GET("/v1/users-by-token/:token", r.userByToken)
	e.POST("/v1/register", r.registerUser)
	e.POST("/v1/auth", r.auth)

	return e
}
