package http

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-ozzo/ozzo-validation/v4"
	realworld "github.com/xesina/gokit-realworld"
	"github.com/xesina/gokit-realworld/article"
	httpError "github.com/xesina/gokit-realworld/http/error"
	"github.com/xesina/gokit-realworld/http/middleware"
	"net/http"
	"strconv"
	"time"
)

type addCommentRequest struct {
	userID  int64
	slug    string
	Comment struct {
		Body string `json:"body"`
	} `json:"comment"`
}

func (req *addCommentRequest) bind(r *http.Request) error {
	_, claims, err := middleware.FromContext(r.Context())
	if err != nil {
		return err
	}

	id := claims["id"].(float64)
	req.userID = int64(id)

	req.slug = chi.URLParam(r, "slug")

	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return httpError.NewError(http.StatusUnprocessableEntity, httpError.ErrRequestBody)
	}

	if err := req.validate(); err != nil {
		return err
	}

	return nil
}

func (req *addCommentRequest) validate() error {
	return validation.ValidateStruct(
		&req.Comment,
		validation.Field(&req.Comment.Body, validation.Required),
	)
}

func (req *addCommentRequest) endpointRequest() article.AddCommentRequest {
	return article.AddCommentRequest{
		Slug:   req.slug,
		UserID: req.userID,
		Body:   req.Comment.Body,
	}
}

func (h ArticleHandler) decodeAddCommentRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req addCommentRequest
	if err := req.bind(r); err != nil {
		return nil, err
	}
	er := req.endpointRequest()
	return er, nil
}

type comment struct {
	ID        int64     `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Author    struct {
		Username  string          `json:"username"`
		Bio       realworld.Bio   `json:"bio"`
		Image     realworld.Image `json:"image"`
		Following bool            `json:"following"`
	} `json:"author"`
}

type commentResponse struct {
	Comment *comment `json:"comment"`
}

func newCommentResponse(c *article.CommentResponse) commentResponse {
	return commentResponse{Comment: &comment{
		ID:        c.ID,
		Body:      c.Body,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Author: Author{
			Username:  c.Author.Username,
			Bio:       c.Author.Bio,
			Image:     c.Author.Image,
			Following: c.Author.Following,
		},
	}}
}

func (h ArticleHandler) encodeCommentResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if resp, ok := response.(endpoint.Failer); ok && resp.Failed() != nil {
		httpError.EncodeError(ctx, resp.Failed(), w)
		return nil
	}
	e := response.(article.CommentResponse)
	return jsonResponse(w, newCommentResponse(&e), http.StatusOK)
}

type deleteCommentRequest struct {
	userID int64
	id     int64
	slug   string
}

func (req *deleteCommentRequest) bind(r *http.Request) error {
	_, claims, err := middleware.FromContext(r.Context())
	if err != nil {
		return err
	}

	id := claims["id"].(float64)
	req.userID = int64(id)

	req.slug = chi.URLParam(r, "slug")
	commentID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		return httpError.NewError(http.StatusUnprocessableEntity, httpError.ErrRequestBody)
	}
	req.id = commentID

	return nil
}

func (req *deleteCommentRequest) endpointRequest() article.DeleteCommentRequest {
	return article.DeleteCommentRequest{
		ID:     req.id,
		UserID: req.userID,
		Slug:   req.slug,
	}
}

func (h ArticleHandler) decodeDeleteCommentRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req deleteCommentRequest
	if err := req.bind(r); err != nil {
		return nil, err
	}
	er := req.endpointRequest()
	return er, nil
}

func (h ArticleHandler) encodeDeleteCommentResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if resp, ok := response.(endpoint.Failer); ok && resp.Failed() != nil {
		httpError.EncodeError(ctx, resp.Failed(), w)
		return nil
	}
	return jsonResponse(w, nil, http.StatusOK)
}

type commentsRequest struct {
	userID int64
	slug   string
}

func (req *commentsRequest) bind(r *http.Request) error {
	_, claims, err := middleware.FromContext(r.Context())
	if err != nil {
		return err
	}

	id := claims["id"].(float64)
	req.userID = int64(id)

	req.slug = chi.URLParam(r, "slug")

	return nil
}

func (req *commentsRequest) endpointRequest() article.CommentsRequest {
	return article.CommentsRequest{
		UserID: req.userID,
		Slug:   req.slug,
	}
}

func (h ArticleHandler) decodeCommentsRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req commentsRequest
	if err := req.bind(r); err != nil {
		return nil, err
	}
	er := req.endpointRequest()
	return er, nil
}

type commentsResponse struct {
	Comments []comment `json:"comments"`
}

func newCommentsResponse(list *article.CommentsResponse) (cc commentsResponse) {
	cc.Comments = make([]comment, 0)

	for _, c := range list.Comments {
		resp := comment{
			ID:        c.ID,
			Body:      c.Body,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Author: Author{
				Username:  c.Author.Username,
				Bio:       c.Author.Bio,
				Image:     c.Author.Image,
				Following: c.Author.Following,
			},
		}
		cc.Comments = append(cc.Comments, resp)
	}
	return
}

func (h ArticleHandler) encodeCommentsResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if resp, ok := response.(endpoint.Failer); ok && resp.Failed() != nil {
		httpError.EncodeError(ctx, resp.Failed(), w)
		return nil
	}
	e := response.(article.CommentsResponse)
	return jsonResponse(w, newCommentsResponse(&e), http.StatusCreated)
}
