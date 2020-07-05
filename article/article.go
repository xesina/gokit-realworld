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

func (s Service) Get(a realworld.Article) (*realworld.Article, error) {
	return s.Store.Get(a.Slug)
}

func (s Service) List(req realworld.ListRequest) ([]*realworld.Article, error) {
	return s.Store.List(req)
}

func (s Service) Favorite(a realworld.Article, u realworld.User) (*realworld.Article, error) {
	return s.Store.AddFavorite(a, u)
}

func (s Service) Unfavorite(a realworld.Article, u realworld.User) (*realworld.Article, error) {
	return s.Store.RemoveFavorite(a, u)
}
