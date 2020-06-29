package http

import (
	transport "github.com/go-kit/kit/transport/http"
	"github.com/xesina/go-kit-realworld-example-app/user"
	"net/http"
)

func (h UserHandler) registerHandlerFunc(opts ...transport.ServerOption) http.HandlerFunc {
	return wrapHandler(transport.NewServer(
		user.RegisterEndpoint(h.service),
		h.decodeRegisterRequest,
		h.encodeUserResponse,
		opts...,
	))
}

func (h UserHandler) loginHandlerFunc(opts ...transport.ServerOption) http.HandlerFunc {
	return wrapHandler(transport.NewServer(
		user.LoginEndpoint(h.service),
		h.decodeLoginRequest,
		h.encodeUserResponse,
		opts...,
	))
}

func (h UserHandler) getHandlerFunc(opts ...transport.ServerOption) http.HandlerFunc {
	return wrapHandler(transport.NewServer(
		user.GetEndpoint(h.service),
		h.decodeGetRequest,
		h.encodeUserResponse,
		opts...,
	))
}
