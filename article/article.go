package article

import (
	realworld "github.com/xesina/gokit-realworld"
)

type Service struct {
	Repo realworld.ArticleRepo
}

func (s Service) Create(a realworld.Article) (*realworld.Article, error) {
	return s.Repo.Create(a)
}

func (s Service) Delete(a realworld.Article) error {
	return s.Repo.Delete(a)
}

func (s Service) Get(a realworld.Article) (*realworld.Article, error) {
	return s.Repo.Get(a.Slug)
}

func (s Service) List(req realworld.ListRequest) ([]*realworld.Article, int, error) {
	switch {
	case req.Tag != "":
		return s.Repo.ListByTag(req.Tag, req.Offset, req.Limit)
	case req.FavoriterID != 0:
		return s.Repo.ListByFavoriterID(req.FavoriterID, req.Offset, req.Limit)
	case req.AuthorID != 0:
		return s.Repo.ListByAuthorID(req.AuthorID, req.Offset, req.Limit)
	default:
		return s.Repo.List(req.Offset, req.Limit)
	}
}

func (s Service) Feed(req realworld.FeedRequest) ([]*realworld.Article, int, error) {
	return s.Repo.Feed(req)
}

func (s Service) Favorite(a realworld.Article, u realworld.User) (*realworld.Article, error) {
	return s.Repo.AddFavorite(a, u)
}

func (s Service) Unfavorite(a realworld.Article, u realworld.User) (*realworld.Article, error) {
	return s.Repo.RemoveFavorite(a, u)
}

func (s Service) AddComment(c realworld.Comment) (*realworld.Comment, error) {
	return s.Repo.AddComment(c)
}

func (s Service) DeleteComment(c realworld.Comment) error {
	return s.Repo.DeleteComment(c)
}

func (s Service) Comments(a realworld.Article) ([]*realworld.Comment, error) {
	return s.Repo.Comments(a)
}

func (s Service) Update(slug string, a realworld.Article) (*realworld.Article, error) {
	return s.Repo.Update(slug, a)
}

func (s Service) Tags() ([]*realworld.Tag, error) {
	return s.Repo.Tags()
}
