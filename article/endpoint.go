package article

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	realworld "github.com/xesina/go-kit-realworld-example-app"
	"time"
)

type CreateRequest struct {
	UserID      int64
	Title       string
	Description string
	Body        string
	Tags        []string
}

func (r CreateRequest) buildTags() (tt realworld.Tags) {
	tt = make(realworld.Tags)
	for _, t := range r.Tags {
		tt[t] = realworld.Tag{Tag: t}
	}
	return
}

func (r CreateRequest) toArticle() (a realworld.Article) {
	a = realworld.Article{
		Title:       r.Title,
		Description: r.Description,
		Body:        r.Body,
	}
	a.Author = realworld.User{ID: r.UserID}
	a.Tags = r.buildTags()
	a.Slug = a.MakeSlug()

	return
}

type Response struct {
	Slug           string
	Title          string
	Description    string
	Body           string
	Tags           realworld.Tags
	Favorited      bool
	FavoritesCount int
	Author         Author
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Err            error
}

func (r Response) TagsList() (tt []string) {
	for _, t := range r.Tags {
		tt = append(tt, t.Tag)
	}

	return
}

type Author struct {
	Username  string
	Bio       realworld.Bio
	Image     realworld.Image
	Following bool
}

func NewResponse(a *realworld.Article, u *realworld.User, err error) Response {
	return Response{
		Slug:           a.Slug,
		Title:          a.Title,
		Description:    a.Description,
		Body:           a.Body,
		Tags:           a.Tags,
		Favorited:      a.Favorited(u.ID),
		FavoritesCount: len(a.Favorites),
		Author: Author{
			Username:  u.Username,
			Bio:       u.Bio,
			Image:     u.Image,
			Following: u.IsFollower(u),
		},
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
		Err:       err,
	}
}

func (r Response) error() error { return r.Err }

func (r Response) Failed() error { return r.Err }

func CreateEndpoint(a realworld.ArticleService, u realworld.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateRequest)
		article, err := a.Create(req.toArticle())
		if err != nil {
			return nil, err
		}
		user, err := u.Get(realworld.User{ID: req.UserID})
		if err != nil {
			return nil, err
		}
		return NewResponse(article, user, err), nil
	}
}

type DeleteRequest struct {
	UserID int64
	Slug   string
}

func (r DeleteRequest) toArticle() (a realworld.Article) {
	a.Author = realworld.User{ID: r.UserID}
	a.Title = r.Slug
	a.Slug = a.MakeSlug()
	return
}

type DeleteResponse struct {
	Err error
}

func DeleteEndpoint(a realworld.ArticleService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteRequest)
		err = a.Delete(req.toArticle())
		if err != nil {
			return nil, err
		}
		return DeleteResponse{}, nil
	}
}
