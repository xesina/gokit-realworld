package article

import (
	realworld "github.com/xesina/go-kit-realworld-example-app"
)

type Service struct {
	Store realworld.ArticleRepo
}

func (s Service) Create(a realworld.Article) (*realworld.Article, error) {
	return s.Store.Create(a)
}

func (s Service) Delete(a realworld.Article) error {
	article, err := s.Store.Get(a.Slug)
	if err != nil {
		return err
	}

	if article.Author.ID != a.Author.ID {
		return realworld.ArticleNotFoundError()
	}

	return s.Store.Delete(a)
}
