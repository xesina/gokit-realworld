package user

import (
	realworld "github.com/xesina/gokit-realworld"
)

type Service struct {
	UserRepo realworld.UserRepo
}

func (s Service) Register(u realworld.User) (*realworld.User, error) {
	hashed, err := u.HashPassword(u.Password)
	if err != nil {
		return nil, realworld.InternalError(err)
	}
	u.Password = hashed
	return s.UserRepo.Create(u)
}

func (s Service) Login(u realworld.User) (*realworld.User, error) {
	found, err := s.UserRepo.Get(u.Email)
	if err != nil {
		return nil, err
	}

	if !found.CheckPassword(u.Password) {
		return nil, realworld.ErrIncorrectPasswordError
	}

	return found, nil
}

func (s Service) Get(u realworld.User) (*realworld.User, error) {
	return s.UserRepo.GetByID(u.ID)
}

func (s Service) Update(u realworld.User) (*realworld.User, error) {
	// TODO: check: this is a full update. should I consider patching instead?
	// TODO: check: where should I check if this user exists at all? in store or service impl?
	if u.Password != "" {
		hashed, err := u.HashPassword(u.Password)
		if err != nil {
			return nil, realworld.InternalError(err)
		}
		u.Password = hashed
	}

	return s.UserRepo.Update(u)
}

func (s Service) GetProfile(user realworld.User) (*realworld.User, error) {
	return s.UserRepo.GetByUsername(user.Username)
}

func (s Service) Follow(req realworld.FollowRequest) (*realworld.User, error) {
	followee, err := s.UserRepo.GetByUsername(req.Followee)
	if err != nil {
		return nil, err
	}

	return s.UserRepo.AddFollower(req.Follower, followee.ID)
}

func (s Service) Unfollow(req realworld.FollowRequest) (*realworld.User, error) {
	followee, err := s.UserRepo.GetByUsername(req.Followee)
	if err != nil {
		return nil, err
	}

	return s.UserRepo.RemoveFollower(req.Follower, followee.ID)
}
