package user

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	realworld "github.com/xesina/go-kit-realworld-example-app"
)

type RegisterRequest struct {
	Username string
	Email    string
	Password string
}

func (r RegisterRequest) toUser() realworld.User {
	return realworld.User{
		Username: r.Username,
		Email:    r.Email,
		Password: r.Password,
	}
}

type Response struct {
	Username string
	Email    string
	Bio      realworld.Bio
	Image    realworld.Image
	Err      error
}

func NewResponse(u *realworld.User, err error) Response {
	return Response{
		Username: u.Username,
		Email:    u.Email,
		Bio:      realworld.Bio{},
		Image:    realworld.Image{},
		Err:      err,
	}
}

func (r Response) error() error { return r.Err }

func RegisterEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(RegisterRequest)
		u, err := s.Register(req.toUser())
		return NewResponse(u, err), nil
	}
}
