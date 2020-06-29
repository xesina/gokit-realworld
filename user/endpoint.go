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
	ID       int64
	Username string
	Email    string
	Bio      realworld.Bio
	Image    realworld.Image
	Err      error
}

func NewResponse(u *realworld.User, err error) Response {
	return Response{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Bio:      realworld.Bio{},
		Image:    realworld.Image{},
		Err:      err,
	}
}

func (r Response) error() error { return r.Err }

func (r Response) Failed() error { return r.Err }

func RegisterEndpoint(s realworld.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(RegisterRequest)
		u, err := s.Register(req.toUser())
		if err != nil {
			return nil, err
		}
		return NewResponse(u, err), nil
	}
}

type LoginRequest struct {
	Email    string
	Password string
}

func (r LoginRequest) toUser() realworld.User {
	return realworld.User{
		Email:    r.Email,
		Password: r.Password,
	}
}

func LoginEndpoint(s realworld.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(LoginRequest)
		u, err := s.Login(req.toUser())
		if err != nil {
			return nil, err
		}
		return NewResponse(u, err), nil
	}
}

type GetRequest struct {
	ID int64
}

func (r GetRequest) toUser() realworld.User {
	return realworld.User{
		ID: r.ID,
	}
}

func GetEndpoint(s realworld.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetRequest)
		u, err := s.Get(req.toUser())
		if err != nil {
			return nil, err
		}
		return NewResponse(u, err), nil
	}
}
