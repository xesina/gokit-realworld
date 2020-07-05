package http

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-kit/kit/endpoint"
	transport "github.com/go-kit/kit/transport/http"
	"github.com/go-ozzo/ozzo-validation/v4"
	realworld "github.com/xesina/go-kit-realworld-example-app"
	"github.com/xesina/go-kit-realworld-example-app/article"
	httpError "github.com/xesina/go-kit-realworld-example-app/http/error"
	"github.com/xesina/go-kit-realworld-example-app/http/middleware"
	"net/http"
	"strconv"
	"time"
)

type ArticleHandler struct {
	service       realworld.ArticleService
	userService   realworld.UserService
	serverOptions []transport.ServerOption
}

func NewArticleHandler(c Context) ArticleHandler {
	return ArticleHandler{
		service:       c.articleService,
		userService:   c.userService,
		serverOptions: c.serverOptions,
	}
}

func (h ArticleHandler) decodeCreateRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req articleCreateRequest
	if err := req.bind(r); err != nil {
		return nil, err
	}
	er := req.endpointRequest()
	return er, nil
}

type articleCreateRequest struct {
	userID  int64
	Article struct {
		Title       string   `json:"title" validate:"required"`
		Description string   `json:"description" validate:"required"`
		Body        string   `json:"body" validate:"required"`
		Tags        []string `json:"tagList, omitempty"`
	} `json:"article"`
}

func (req *articleCreateRequest) bind(r *http.Request) error {
	_, claims, err := middleware.FromContext(r.Context())
	if err != nil {
		return err
	}

	id := claims["id"].(float64)
	req.userID = int64(id)

	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return httpError.NewError(http.StatusUnprocessableEntity, httpError.ErrRequestBody)
	}

	if err := req.validate(); err != nil {
		return err
	}

	return nil
}

func (req *articleCreateRequest) validate() error {
	return validation.ValidateStruct(
		&req.Article,
		validation.Field(&req.Article.Title, validation.Required),
		validation.Field(&req.Article.Description, validation.Required),
		validation.Field(&req.Article.Body, validation.Required),
	)
}

func (req *articleCreateRequest) endpointRequest() article.CreateRequest {
	return article.CreateRequest{
		UserID:      req.userID,
		Title:       req.Article.Title,
		Description: req.Article.Description,
		Body:        req.Article.Body,
		Tags:        req.Article.Tags,
	}
}

type Author struct {
	Username  string          `json:"username"`
	Bio       realworld.Bio   `json:"bio"`
	Image     realworld.Image `json:"image"`
	Following bool            `json:"following"`
}

type articleResponse struct {
	Slug           string    `json:"slug"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Body           string    `json:"body"`
	Tags           []string  `json:"tagList"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Favorited      bool      `json:"favorited"`
	FavoritesCount int       `json:"favoritesCount"`
	Author         Author    `json:"author"`
}

type singleArticleResponse struct {
	Article *articleResponse `json:"article"`
}

func newArticleResponse(a *article.Response) singleArticleResponse {
	return singleArticleResponse{Article: &articleResponse{
		Slug:           a.Slug,
		Title:          a.Title,
		Description:    a.Description,
		Body:           a.Body,
		Tags:           a.TagsList(),
		CreatedAt:      a.CreatedAt,
		UpdatedAt:      a.UpdatedAt,
		Favorited:      a.Favorited,
		FavoritesCount: a.FavoritesCount,
		Author: Author{
			Username:  a.Author.Username,
			Bio:       a.Author.Bio,
			Image:     a.Author.Image,
			Following: a.Author.Following,
		},
	}}
}

func (h ArticleHandler) encodeArticlesResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if resp, ok := response.(endpoint.Failer); ok && resp.Failed() != nil {
		httpError.EncodeError(ctx, resp.Failed(), w)
		return nil
	}
	e := response.(article.ListResponse)
	return jsonResponse(w, newArticlesResponse(&e), http.StatusCreated)
}

type deleteRequest struct {
	userID int64
	slug   string
}

func (req *deleteRequest) bind(r *http.Request) error {
	_, claims, err := middleware.FromContext(r.Context())
	if err != nil {
		return err
	}

	id := claims["id"].(float64)
	req.userID = int64(id)

	req.slug = chi.URLParam(r, "slug")

	if err := req.validate(); err != nil {
		return err
	}

	return nil
}

func (req *deleteRequest) validate() error {
	return validation.ValidateStruct(
		req,
		validation.Field(&req.slug, validation.Required),
	)
}

func (req *deleteRequest) endpointRequest() article.DeleteRequest {
	return article.DeleteRequest{
		UserID: req.userID,
		Slug:   req.slug,
	}
}

func (h ArticleHandler) decodeDeleteRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req deleteRequest
	if err := req.bind(r); err != nil {
		return nil, err
	}
	er := req.endpointRequest()
	return er, nil
}

func (h ArticleHandler) encodeDeleteResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if resp, ok := response.(endpoint.Failer); ok && resp.Failed() != nil {
		httpError.EncodeError(ctx, resp.Failed(), w)
		return nil
	}
	return jsonResponse(w, nil, http.StatusOK)
}

type getRequest struct {
	userID int64
	slug   string
}

func (req *getRequest) bind(r *http.Request) error {
	_, claims, err := middleware.FromContext(r.Context())
	if err != nil {
		return err
	}

	id := claims["id"].(float64)
	req.userID = int64(id)

	req.slug = chi.URLParam(r, "slug")

	if err := req.validate(); err != nil {
		return err
	}

	return nil
}

func (req *getRequest) validate() error {
	return validation.ValidateStruct(
		req,
		validation.Field(&req.slug, validation.Required),
	)
}

func (req *getRequest) endpointRequest() article.GetRequest {
	return article.GetRequest{
		UserID: req.userID,
		Slug:   req.slug,
	}
}

func (h ArticleHandler) decodeGetRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req getRequest
	if err := req.bind(r); err != nil {
		return nil, err
	}
	er := req.endpointRequest()
	return er, nil
}

type favoriteRequest struct {
	userID int64
	slug   string
}

func (req *favoriteRequest) bind(r *http.Request) error {
	_, claims, err := middleware.FromContext(r.Context())
	if err != nil {
		return err
	}

	id := claims["id"].(float64)
	req.userID = int64(id)

	req.slug = chi.URLParam(r, "slug")

	if err := req.validate(); err != nil {
		return err
	}

	return nil
}

func (req *favoriteRequest) validate() error {
	return validation.ValidateStruct(
		req,
		validation.Field(&req.slug, validation.Required),
	)
}

func (req *favoriteRequest) endpointRequest() article.FavoriteRequest {
	return article.FavoriteRequest{
		UserID: req.userID,
		Slug:   req.slug,
	}
}

func (h ArticleHandler) decodeFavoriteRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req favoriteRequest
	if err := req.bind(r); err != nil {
		return nil, err
	}
	er := req.endpointRequest()
	return er, nil
}

type listRequest struct {
	userID      int64
	tag         string
	author      string
	authorID    int64
	favoritedBy string
	favoriterID int64
	limit       int
	offset      int
}

func (req *listRequest) bind(r *http.Request) error {
	_, claims, err := middleware.FromContext(r.Context())
	if err != nil {
		return err
	}

	id := claims["id"].(float64)
	req.userID = int64(id)

	req.tag = r.URL.Query().Get("tag")
	req.author = r.URL.Query().Get("author")
	req.favoritedBy = r.URL.Query().Get("favorited")

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 0
	}
	req.limit = limit

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}
	req.offset = offset

	if err := req.validate(); err != nil {
		return err
	}

	return nil
}

func (req *listRequest) validate() error {
	return nil
}

func (req *listRequest) endpointRequest() article.ListRequest {
	return article.ListRequest{
		UserID:      req.userID,
		Tag:         req.tag,
		AuthorID:    req.authorID,
		FavoriterID: req.favoriterID,
		Limit:       req.limit,
		Offset:      req.offset,
	}
}

func (h ArticleHandler) decodeListRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req listRequest
	if err := req.bind(r); err != nil {
		return nil, err
	}

	if req.favoritedBy != "" {
		favoriter, err := h.userService.GetProfile(realworld.User{Username: req.favoritedBy})
		if err != nil {
			return nil, err
		}
		req.favoriterID = favoriter.ID
	}

	if req.author != "" {
		author, err := h.userService.GetProfile(realworld.User{Username: req.author})
		if err != nil {
			return nil, err
		}
		req.authorID = author.ID
	}

	er := req.endpointRequest()
	return er, nil
}

type articleListResponse struct {
	Articles      []*articleResponse `json:"articles"`
	ArticlesCount int                `json:"articlesCount"`
}

func newArticlesResponse(list *article.ListResponse) (aa articleListResponse) {
	aa.Articles = make([]*articleResponse, 0)

	for _, a := range list.Articles {
		resp := articleResponse{
			Slug:           a.Slug,
			Title:          a.Title,
			Description:    a.Description,
			Body:           a.Body,
			Tags:           a.Tags.TagsList(),
			CreatedAt:      a.CreatedAt,
			UpdatedAt:      a.UpdatedAt,
			Favorited:      a.Favorited,
			FavoritesCount: a.FavoritesCount,
			Author: Author{
				Username:  a.Author.Username,
				Bio:       a.Author.Bio,
				Image:     a.Author.Image,
				Following: a.Author.Following,
			},
		}
		aa.Articles = append(aa.Articles, &resp)
	}
	aa.ArticlesCount = len(list.Articles)
	return
}

func (h ArticleHandler) encodeArticleResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if resp, ok := response.(endpoint.Failer); ok && resp.Failed() != nil {
		httpError.EncodeError(ctx, resp.Failed(), w)
		return nil
	}
	e := response.(article.Response)
	return jsonResponse(w, newArticleResponse(&e), http.StatusCreated)
}