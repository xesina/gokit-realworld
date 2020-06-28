package user

import (
	realworld "github.com/xesina/go-kit-realworld-example-app"
)

type Service struct {
	Store realworld.UserRepo
}

func (s Service) Register(u realworld.User) (*realworld.User, error) {
	return s.Store.Create(u)
}

func (s Service) Get(email string) (*realworld.User, error) {
	return s.Store.Get(email)
}
