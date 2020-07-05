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

type Article struct {
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
}

type Response struct {
	Article
	Err error
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

func NewResponse(a *realworld.Article, u realworld.User, userSrv realworld.UserService, err error) Response {
	viewer, err := userSrv.Get(u)
	if err != nil {
		return Response{
			Err: err,
		}
	}

	author, err := userSrv.Get(a.Author)
	if err != nil {
		return Response{
			Err: err,
		}
	}

	return Response{
		Article{
			Slug:           a.Slug,
			Title:          a.Title,
			Description:    a.Description,
			Body:           a.Body,
			Tags:           a.Tags,
			Favorited:      a.Favorited(viewer.ID),
			FavoritesCount: len(a.Favorites),
			Author: Author{
				Username:  author.Username,
				Bio:       author.Bio,
				Image:     author.Image,
				Following: author.IsFollower(viewer),
			},
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
		},
		err,
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
		return NewResponse(article, realworld.User{ID: req.UserID}, u, err), nil
	}
}

type UpdateRequest struct {
	TargetSlug  string
	Slug        string
	UserID      int64
	Title       string
	Description string
	Body        string
}

func (r UpdateRequest) toArticle() (a realworld.Article) {
	a = realworld.Article{
		Title:       r.Title,
		Description: r.Description,
		Body:        r.Body,
	}
	a.Author = realworld.User{ID: r.UserID}
	a.Slug = a.MakeSlug()
	return
}

func UpdateEndpoint(a realworld.ArticleService, u realworld.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateRequest)
		article, err := a.Update(req.TargetSlug, req.toArticle())
		if err != nil {
			return nil, err
		}
		return NewResponse(article, realworld.User{ID: req.UserID}, u, err), nil
	}
}

type GetRequest struct {
	UserID int64
	Slug   string
}

func (r GetRequest) toArticle() (a realworld.Article) {
	a = realworld.Article{
		Slug: r.Slug,
	}
	return
}

func GetEndpoint(a realworld.ArticleService, u realworld.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetRequest)
		article, err := a.Get(req.toArticle())
		if err != nil {
			return nil, err
		}
		return NewResponse(article, realworld.User{ID: req.UserID}, u, err), nil
	}
}

type ListResponse struct {
	Articles []Article
	Err      error
}

func NewListResponse(
	articles []*realworld.Article, u *realworld.User, userSrv realworld.UserService, err error,
) ListResponse {
	var listResponse ListResponse
	for _, article := range articles {
		author, err := userSrv.Get(realworld.User{ID: article.Author.ID})
		if err != nil {
			return ListResponse{nil, err}
		}

		resp := Article{
			Slug:           article.Slug,
			Title:          article.Title,
			Description:    article.Description,
			Body:           article.Body,
			Tags:           article.Tags,
			Favorited:      article.Favorited(u.ID),
			FavoritesCount: len(article.Favorites),
			Author: Author{
				Username:  author.Username,
				Bio:       author.Bio,
				Image:     author.Image,
				Following: author.IsFollower(u),
			},
			CreatedAt: article.CreatedAt,
			UpdatedAt: article.UpdatedAt,
		}

		listResponse.Articles = append(listResponse.Articles, resp)
	}
	listResponse.Err = err

	return listResponse
}

func (r ListResponse) error() error { return r.Err }

func (r ListResponse) Failed() error { return r.Err }

type ListRequest struct {
	UserID      int64
	Tag         string
	AuthorID    int64
	FavoriterID int64
	Limit       int
	Offset      int
}

func (req ListRequest) serviceRequest() realworld.ListRequest {
	return realworld.ListRequest{
		Tag:         req.Tag,
		AuthorID:    req.AuthorID,
		FavoriterID: req.FavoriterID,
		Offset:      req.Offset,
		Limit:       req.Limit,
	}
}

func ListEndpoint(a realworld.ArticleService, u realworld.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ListRequest)
		aa, err := a.List(req.serviceRequest())
		if err != nil {
			return nil, err
		}
		user, err := u.Get(realworld.User{ID: req.UserID})
		if err != nil {
			return nil, err
		}
		return NewListResponse(aa, user, u, err), nil
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

type FavoriteRequest struct {
	UserID int64
	Slug   string
}

func (r FavoriteRequest) toArticle() (a realworld.Article) {
	a.Slug = r.Slug
	return
}

func (r FavoriteRequest) toUser() (u realworld.User) {
	u.ID = r.UserID
	return
}

func FavoriteEndpoint(a realworld.ArticleService, u realworld.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(FavoriteRequest)
		article, err := a.Favorite(req.toArticle(), req.toUser())
		if err != nil {
			return nil, err
		}
		return NewResponse(article, realworld.User{ID: req.UserID}, u, err), nil
	}
}

func UnfavoriteEndpoint(a realworld.ArticleService, u realworld.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(FavoriteRequest)
		article, err := a.Unfavorite(req.toArticle(), req.toUser())
		if err != nil {
			return nil, err
		}
		return NewResponse(article, realworld.User{ID: req.UserID}, u, err), nil
	}
}
