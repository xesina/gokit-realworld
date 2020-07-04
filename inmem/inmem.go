package inmem

import (
	realworld "github.com/xesina/go-kit-realworld-example-app"
	"sync"
	"sync/atomic"
)

func NewMemUserSaver() realworld.UserRepo {
	return &memUserSaver{
		m: map[string]realworld.User{},
	}
}

type memUserSaver struct {
	rwlock  sync.RWMutex
	m       map[string]realworld.User
	counter int64
}

func (store *memUserSaver) Create(u realworld.User) (*realworld.User, error) {
	store.rwlock.Lock()
	defer store.rwlock.Unlock()

	if _, ok := store.m[u.Email]; ok {
		return nil, realworld.UserAlreadyExistsError(u.Email)
	}

	u.ID = atomic.AddInt64(&store.counter, 1)
	u.Followings = make(map[int64]struct{})
	u.Followers = make(map[int64]struct{})
	store.m[u.Email] = u
	return &u, nil
}

func (store *memUserSaver) Update(u realworld.User) (*realworld.User, error) {
	store.rwlock.Lock()
	defer store.rwlock.Unlock()

	old, ok := store.m[u.Email]
	if !ok {
		return nil, realworld.UserNotFoundError()
	}

	// TODO: I'm not sure if this is the best place to prevent followers/followings from change in update
	u.Followers = old.Followers
	u.Followings = old.Followings

	store.m[u.Email] = u

	// If user has changes her email we need to create a new map entry
	// preserving ID as before and also delete the old email and entry
	if u.Email != old.Email {
		delete(store.m, old.Email)
	}

	return &u, nil
}

func (store *memUserSaver) Get(e string) (*realworld.User, error) {
	user, ok := store.m[e]
	if !ok {
		return nil, realworld.UserNotFoundError()
	}
	return &user, nil
}

func (store *memUserSaver) GetByID(id int64) (*realworld.User, error) {
	store.rwlock.RLock()
	defer store.rwlock.RUnlock()
	var email string
	for k, v := range store.m {
		if v.ID == id {
			email = k
			break
		}
	}
	user, ok := store.m[email]
	if !ok {
		return nil, realworld.UserNotFoundError()
	}
	return &user, nil
}

func (store *memUserSaver) GetByUsername(username string) (*realworld.User, error) {
	store.rwlock.RLock()
	defer store.rwlock.RUnlock()
	var email string
	for k, v := range store.m {
		if v.Username == username {
			email = k
			break
		}
	}
	user, ok := store.m[email]
	if !ok {
		return nil, realworld.UserNotFoundError()
	}
	return &user, nil
}

func (store *memUserSaver) AddFollower(follower, followee int64) (*realworld.User, error) {
	followerUser, err := store.GetByID(follower)
	if err != nil {
		return nil, err
	}
	followeeUser, err := store.GetByID(followee)
	if err != nil {
		return nil, err
	}

	// We maybe need to do the following commands within a transaction if we are
	// using a relational storage.

	store.rwlock.Lock()
	defer store.rwlock.Unlock()

	// add follower
	store.m[followeeUser.Email].Followers[follower] = struct{}{}
	// add followee
	store.m[followerUser.Email].Followings[followee] = struct{}{}

	return followeeUser, nil
}

func (store *memUserSaver) RemoveFollower(follower, followee int64) (*realworld.User, error) {
	followerUser, err := store.GetByID(follower)
	if err != nil {
		return nil, err
	}
	followeeUser, err := store.GetByID(followee)
	if err != nil {
		return nil, err
	}

	store.rwlock.Lock()
	defer store.rwlock.Unlock()

	delete(store.m[followeeUser.Email].Followers, follower)
	delete(store.m[followerUser.Email].Followings, followee)

	return followeeUser, nil
}
