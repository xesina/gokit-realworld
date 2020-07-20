package inmem

import (
	realworld "github.com/xesina/go-kit-realworld-example-app"
	"sync/atomic"
	"time"
)

func (store *memArticleRepo) AddComment(c realworld.Comment) (comment *realworld.Comment, err error) {
	article, ok := store.m[c.Article.Slug]
	if !ok {
		return nil, realworld.ErrArticleNotFound
	}

	c.ID = atomic.AddInt64(&store.counter, 1)
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	article.Comments[c.ID] = c

	return &c, nil
}

func (store *memArticleRepo) DeleteComment(c realworld.Comment) error {
	article, ok := store.m[c.Article.Slug]
	if !ok {
		return realworld.ErrArticleNotFound
	}

	delete(article.Comments, c.ID)

	return nil
}

func (store *memArticleRepo) Comments(a realworld.Article) ([]*realworld.Comment, error) {
	article, ok := store.m[a.Slug]
	if !ok {
		return nil, realworld.ErrArticleNotFound
	}

	var comments []*realworld.Comment
	for k, _ := range article.Comments {
		c := article.Comments[k]
		comments = append(comments, &c)
	}

	return comments, nil
}

