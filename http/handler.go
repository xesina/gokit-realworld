package http

import (
	transport "github.com/go-kit/kit/transport/http"
	"github.com/xesina/go-kit-realworld-example-app/article"
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

func (h UserHandler) updateHandlerFunc() http.HandlerFunc {
	return wrapHandler(transport.NewServer(
		user.UpdateEndpoint(h.service),
		h.decodeUpdateRequest,
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

func (h UserHandler) followHandlerFunc() http.HandlerFunc {
	return wrapHandler(transport.NewServer(
		user.FollowEndpoint(h.service),
		h.decodeProfileRequest,
		h.encodeProfileResponse,
		h.serverOptions...,
	))
}

func (h UserHandler) unfollowHandlerFunc() http.HandlerFunc {
	return wrapHandler(transport.NewServer(
		user.UnfollowEndpoint(h.service),
		h.decodeProfileRequest,
		h.encodeProfileResponse,
		h.serverOptions...,
	))
}

func (h ArticleHandler) createHandlerFunc() http.HandlerFunc {
	return wrapHandler(transport.NewServer(
		article.CreateEndpoint(h.service, h.userService),
		h.decodeCreateRequest,
		h.encodeArticleResponse,
		h.serverOptions...,
	))
}

func (h ArticleHandler) getHandlerFunc() http.HandlerFunc {
	return wrapHandler(transport.NewServer(
		article.GetEndpoint(h.service, h.userService),
		h.decodeGetRequest,
		h.encodeArticleResponse,
		h.serverOptions...,
	))
}

func (h ArticleHandler) listHandlerFunc() http.HandlerFunc {
	return wrapHandler(transport.NewServer(
		article.ListEndpoint(h.service, h.userService),
		h.decodeListRequest,
		h.encodeArticlesResponse,
		h.serverOptions...,
	))
}

func (h ArticleHandler) deleteHandlerFunc() http.HandlerFunc {
	return wrapHandler(transport.NewServer(
		article.DeleteEndpoint(h.service),
		h.decodeDeleteRequest,
		h.encodeDeleteResponse,
		h.serverOptions...,
	))
}

func (h ArticleHandler) favoriteHandlerFunc() http.HandlerFunc {
	return wrapHandler(transport.NewServer(
		article.FavoriteEndpoint(h.service, h.userService),
		h.decodeFavoriteRequest,
		h.encodeArticleResponse,
		h.serverOptions...,
	))
}

func (h ArticleHandler) unfavoriteHandlerFunc() http.HandlerFunc {
	return wrapHandler(transport.NewServer(
		article.UnfavoriteEndpoint(h.service, h.userService),
		h.decodeFavoriteRequest,
		h.encodeArticleResponse,
		h.serverOptions...,
	))
}
