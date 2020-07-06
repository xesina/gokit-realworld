package go_kit_realworld_example_app

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type Bio struct {
	Value string
	Valid bool
}

func (b *Bio) UnmarshalJSON(bytes []byte) error {
	bio := string(bytes)
	bio = strings.Trim(bio, `"`)
	if strings.ToLower(bio) == "null" {
		b.Valid = false
		return nil
	}

	b.Valid = true
	b.Value = bio
	return nil
}

func (b Bio) MarshalJSON() ([]byte, error) {
	if b.Valid {
		return json.Marshal(b.Value)
	}

	return json.Marshal(nil)
}

type Image struct {
	Value string
	Valid bool
}

// TODO: check: if we can extract these methods to a generic methods/type and attach these to the type
func (i *Image) UnmarshalJSON(bytes []byte) error {
	image := string(bytes)
	image = strings.Trim(image, `"`)
	if strings.ToLower(image) == "null" {
		i.Valid = false
		return nil
	}

	i.Valid = true
	i.Value = image
	return nil
}

func (i Image) MarshalJSON() ([]byte, error) {
	if i.Valid {
		return json.Marshal(i.Value)
	}

	return json.Marshal(nil)
}

type FollowRequest struct {
	Followee string
	Follower int64
}

type Follows map[int64]struct{}

func (ff Follows) List() (l []int64) {
	l = make([]int64, 0)
	for f, _ := range ff {
		l = append(l, f)
	}
	return
}

// TODO: where should we put the validation? Should we have separate validation per domain model and transports?
type User struct {
	ID         int64
	Username   string
	Email      string
	Password   string
	Bio        Bio
	Image      Image
	Followers  Follows
	Followings Follows
}

func (u *User) HashPassword(plain string) (string, error) {
	if len(plain) == 0 {
		return "", errors.New("password should not be empty")
	}
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(h), err
}

func (u *User) CheckPassword(plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain))
	return err == nil
}

func (u *User) IsFollower(follower *User) bool {
	if u.Followers == nil {
		return false
	}

	if _, ok := u.Followers[follower.ID]; ok {
		return true
	}

	return false
}

type UserService interface {
	Register(user User) (*User, error)
	Login(user User) (*User, error)
	Get(user User) (*User, error)
	Update(user User) (*User, error)
	GetProfile(user User) (*User, error)
	Follow(req FollowRequest) (*User, error)
	Unfollow(req FollowRequest) (*User, error)
}

type UserRepo interface {
	// TODO: should this return user? What if we assume this should only be a **write** command
	Create(u User) (*User, error)
	Update(u User) (*User, error)
	Get(e string) (*User, error)
	GetByID(id int64) (*User, error)
	GetByUsername(u string) (*User, error)
	AddFollower(follower, followee int64) (*User, error)
	RemoveFollower(follower, followee int64) (*User, error)
}

func UserAlreadyExistsError(email string) error {
	return Error{
		Code: EConflict,
		Err:  fmt.Errorf("user with email: %s already exists", email),
	}
}

func UserNotFoundError() error {
	return Error{
		Code: ENotFound,
		Err:  errors.New("user not found"),
	}
}

func IncorrectPasswordError() error {
	return Error{
		Code: EIncorrectPassword,
		Err:  errors.New("incorrect password"),
	}
}
