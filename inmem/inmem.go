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
