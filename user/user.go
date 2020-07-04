package user

import (
	realworld "github.com/xesina/go-kit-realworld-example-app"
)

type Service struct {
	Store realworld.UserRepo
}

func (s Service) Register(u realworld.User) (*realworld.User, error) {
	hashed, err := u.HashPassword(u.Password)
	if err != nil {
		return nil, realworld.InternalError(err)
	}
	u.Password = hashed
	return s.Store.Create(u)
}

func (s Service) Login(u realworld.User) (*realworld.User, error) {
	found, err := s.Store.Get(u.Email)
	if err != nil {
		return nil, err
	}

	if !found.CheckPassword(u.Password) {
		return nil, realworld.IncorrectPasswordError()
	}

	return found, nil
}

func (s Service) Get(u realworld.User) (*realworld.User, error) {
	return s.Store.GetByID(u.ID)
}

func (s Service) Update(u realworld.User) (*realworld.User, error) {
	// TODO: check: this is a full update. should I consider patching instead?
	// TODO: check: where should I check if this user exists at all? in store or service impl?
	hashed, err := u.HashPassword(u.Password)
	if err != nil {
		return nil, realworld.InternalError(err)
	}
	u.Password = hashed

	return s.Store.Update(u)
}

func (s Service) GetProfile(user realworld.User) (*realworld.User, error) {
	return s.Store.GetByUsername(user.Username)
}

func (s Service) Follow(req realworld.FollowRequest) (*realworld.User, error) {
	followee, err := s.Store.GetByUsername(req.Followee)
	if err != nil {
		return nil, err
	}

	return s.Store.AddFollower(req.Follower, followee.ID)
}

func (s Service) Unfollow(req realworld.FollowRequest) (*realworld.User, error) {
	followee, err := s.Store.GetByUsername(req.Followee)
	if err != nil {
		return nil, err
	}

	return s.Store.RemoveFollower(req.Follower, followee.ID)
}
