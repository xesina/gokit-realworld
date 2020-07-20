package http

import (
	"github.com/go-chi/chi"
	transport "github.com/go-kit/kit/transport/http"
	realworld "github.com/xesina/gokit-realworld"
	"github.com/xesina/gokit-realworld/http/middleware"
)

type Context struct {
	router         *chi.Mux
	jwt            *middleware.JWTAuth
	serverOptions  []transport.ServerOption
	userService    realworld.UserService
	articleService realworld.ArticleService
}
