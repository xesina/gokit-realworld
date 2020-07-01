package http

import (
	transport "github.com/go-kit/kit/transport/http"
	"github.com/xesina/go-kit-realworld-example-app/user"
	"net/http"
)

func (h UserHandler) registerHandlerFunc() http.HandlerFunc {
	return wrapHandler(transport.NewServer(
		user.RegisterEndpoint(h.service),
		h.decodeRegisterRequest,
		h.encodeUserResponse,
		h.serverOptions...,
	))
}

func (h UserHandler) loginHandlerFunc() http.HandlerFunc {
	return wrapHandler(transport.NewServer(
		user.LoginEndpoint(h.service),
		h.decodeLoginRequest,
		h.encodeUserResponse,
		h.serverOptions...,
	))
}

func (h UserHandler) getHandlerFunc() http.HandlerFunc {
	return wrapHandler(transport.NewServer(
		user.GetEndpoint(h.service),
		h.decodeGetRequest,
		h.encodeUserResponse,
		h.serverOptions...,
	))
}

func (h UserHandler) profileHandlerFunc() http.HandlerFunc {
	return wrapHandler(transport.NewServer(
		user.GetProfileEndpoint(h.service),
		h.decodeProfileRequest,
		h.encodeProfileResponse,
		h.serverOptions...,
	))
}
