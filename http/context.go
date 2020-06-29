package http

import (
	"github.com/go-chi/chi"
	transport "github.com/go-kit/kit/transport/http"
	realworld "github.com/xesina/go-kit-realworld-example-app"
	"github.com/xesina/go-kit-realworld-example-app/http/middleware"
)

type Context struct {
	router        *chi.Mux
	jwt           *middleware.JWTAuth
	serverOptions []transport.ServerOption
	userService   realworld.UserService
}
