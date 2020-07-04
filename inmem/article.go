package inmem

import (
	realworld "github.com/xesina/go-kit-realworld-example-app"
	"sync"
	"sync/atomic"
	"time"
)

func NewMemArticleRepo() realworld.ArticleRepo {
	return &memArticleRepo{
		m: map[string]realworld.Article{},
	}
}

type memArticleRepo struct {
	rwlock  sync.RWMutex
	m       map[string]realworld.Article
	counter int64
}

func (store *memArticleRepo) Create(a realworld.Article) (*realworld.Article, error) {
	store.rwlock.Lock()
	defer store.rwlock.Unlock()

	if _, ok := store.m[a.Slug]; ok {
		return nil, realworld.ArticleAlreadyExistsError(a.Slug)
	}

	a.ID = atomic.AddInt64(&store.counter, 1)
	a.Favorites = make(realworld.Favorites, 0)
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()

	store.m[a.Slug] = a
	return &a, nil
}

func (store *memArticleRepo) Delete(a realworld.Article) error {
	store.rwlock.Lock()
	defer store.rwlock.Unlock()
	delete(store.m, a.Slug)
	return nil
}

func (store *memArticleRepo) Get(slug string) (*realworld.Article, error) {
	store.rwlock.RLock()
	defer store.rwlock.RUnlock()

	article, ok := store.m[slug]

	if !ok {
		return nil, realworld.ArticleNotFoundError()
	}

	return &article, nil
}
