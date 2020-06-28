package inmem

import (
	"errors"
	realworld "github.com/xesina/go-kit-realworld-example-app"
)

func NewMemUserSaver() realworld.UserRepo {
	return &memUserSaver{
		m: map[string]realworld.User{},
	}
}

type memUserSaver struct {
	m map[string]realworld.User
}

func (store *memUserSaver) Create(u realworld.User) (*realworld.User, error) {
	if _, ok := store.m[u.Email]; ok {
		return nil, errors.New("user already exists")
	}
	// TODO: set id
	// u.ID = random()
	store.m[u.Email] = u
	return &u, nil
}

func (store *memUserSaver) Get(e string) (*realworld.User, error) {
	user, ok := store.m[e]
	if !ok {
		return nil, errors.New("user not found")
	}
	return &user, nil
}
