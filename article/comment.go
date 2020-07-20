package article

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	realworld "github.com/xesina/gokit-realworld"
	"time"
)

type AddCommentRequest struct {
	Slug   string
	UserID int64
	Body   string
}

func (req AddCommentRequest) toComment() realworld.Comment {
	return realworld.Comment{
		Article: realworld.Article{Slug: req.Slug},
		UserID:  req.UserID,
		Body:    req.Body,
	}
}

type Comment struct {
	ID        int64
	Body      string
	Author    Author
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CommentResponse struct {
	Comment
	Err error
}

func NewCommentResponse(c *realworld.Comment, u *realworld.User, userSrv realworld.UserService, err error) CommentResponse {
	author, err := userSrv.Get(realworld.User{ID: c.UserID})
	if err != nil {
		return CommentResponse{
			Err: err,
		}
	}

	return CommentResponse{
		Comment: Comment{
			ID:   c.ID,
			Body: c.Body,
			Author: Author{
				Username:  author.Username,
				Bio:       author.Bio,
				Image:     author.Image,
				Following: author.IsFollower(u),
			},
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		},
		Err: err,
	}
}

func (r CommentResponse) error() error { return r.Err }

func (r CommentResponse) Failed() error { return r.Err }

func AddCommentEndpoint(a realworld.ArticleService, u realworld.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AddCommentRequest)
		comment, err := a.AddComment(req.toComment())
		if err != nil {
			return nil, err
		}
		return NewCommentResponse(comment, &realworld.User{ID: req.UserID}, u, err), nil
	}
}

type DeleteCommentRequest struct {
	ID     int64
	Slug   string
	UserID int64
}

func (req DeleteCommentRequest) toComment() realworld.Comment {
	return realworld.Comment{
		ID:      req.ID,
		Article: realworld.Article{Slug: req.Slug},
		UserID:  req.UserID,
	}
}

func DeleteCommentEndpoint(a realworld.ArticleService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteCommentRequest)
		err = a.DeleteComment(req.toComment())
		if err != nil {
			return nil, err
		}
		return DeleteResponse{}, nil
	}
}

type CommentsRequest struct {
	UserID int64
	Slug   string
}

func (req CommentsRequest) toArticle() realworld.Article {
	return realworld.Article{
		Slug: req.Slug,
	}
}

type CommentsResponse struct {
	Comments []Comment
	Err      error
}

func NewCommentsResponse(
	cc []*realworld.Comment, u *realworld.User, userSrv realworld.UserService, err error,
) CommentsResponse {
	var comments CommentsResponse
	for _, comment := range cc {
		author, err := userSrv.Get(realworld.User{ID: comment.UserID})
		if err != nil {
			return CommentsResponse{nil, err}
		}

		resp := Comment{
			ID:   comment.ID,
			Body: comment.Body,
			Author: Author{
				Username:  author.Username,
				Bio:       author.Bio,
				Image:     author.Image,
				Following: author.IsFollower(u),
			},
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
		}

		comments.Comments = append(comments.Comments, resp)
	}
	comments.Err = err

	return comments
}

func CommentsEndpoint(a realworld.ArticleService, u realworld.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CommentsRequest)
		cc, err := a.Comments(req.toArticle())
		if err != nil {
			return nil, err
		}
		return NewCommentsResponse(cc, &realworld.User{ID: req.UserID}, u, err), nil
	}
}
