package go_kit_realworld_example_app

import (
	"errors"
	"fmt"
	"github.com/gosimple/slug"
	"time"
)

type Favorites map[int64]struct{}

func (ff Favorites) FavoritedBy(id int64) bool {
	_, ok := ff[id]
	return ok
}

type Tags map[string]Tag

func (tt Tags) HasTag(t string) bool {
	_, ok := tt[t]
	return ok
}

func (tt Tags) TagsList() (tagList []string) {
	for k, _ := range tt {
		tagList = append(tagList, k)
	}
	return
}

type Comments map[int64]Comment

type Article struct {
	ID          int64
	Slug        string
	Title       string
	Description string
	Body        string
	Author      User
	Comments    Comments
	Favorites   Favorites
	Tags        Tags
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (a Article) MakeSlug() string {
	return slug.Make(a.Title)
}

func (a Article) Favorited(id int64) bool {
	if a.Favorites == nil {
		return false
	}

	if _, ok := a.Favorites[id]; ok {
		return true
	}

	return false
}

type ListRequest struct {
	Tag         string
	AuthorID    int64
	FavoriterID int64
	Offset      int
	Limit       int
}

type FeedRequest struct {
	UserID       int64
	FollowingIDs []int64
	Limit        int
	Offset       int
}

type ArticleService interface {
	Create(a Article) (*Article, error)
	Update(slug string, a Article) (*Article, error)
	Get(a Article) (*Article, error)
	List(r ListRequest) ([]*Article, error)
	Feed(r FeedRequest) ([]*Article, error)
	Delete(a Article) error
	Favorite(a Article, u User) (*Article, error)
	Unfavorite(a Article, u User) (*Article, error)
	AddComment(c Comment) (*Comment, error)
	DeleteComment(c Comment) error
	Comments(a Article) ([]*Comment, error)
	Tags() ([]*Tag, error)
}

type ArticleRepo interface {
	Get(slug string) (*Article, error)
	List(req ListRequest) ([]*Article, error)
	Feed(req FeedRequest) ([]*Article, error)
	Create(u Article) (*Article, error)
	Update(slug string, u Article) (*Article, error)
	Delete(u Article) error
	AddFavorite(a Article, u User) (*Article, error)
	RemoveFavorite(a Article, u User) (*Article, error)
	AddComment(c Comment) (*Comment, error)
	DeleteComment(c Comment) error
	Comments(a Article) ([]*Comment, error)
	Tags() ([]*Tag, error)
}

type Comment struct {
	ID        int64
	Article   Article
	ArticleID uint
	UserID    int64
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Tag struct {
	ID       int64
	Tag      string
	Articles []Article
}

func ArticleAlreadyExistsError(slug string) error {
	return Error{
		Code: EConflict,
		Err:  fmt.Errorf("article with slug: %s already exists", slug),
	}
}

func ArticleNotFoundError() error {
	return Error{
		Code: ENotFound,
		Err:  errors.New("article not found"),
	}
}
