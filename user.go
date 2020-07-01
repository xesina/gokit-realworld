package go_kit_realworld_example_app

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type Bio struct {
	Value string
	Valid bool
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

func (i Image) MarshalJSON() ([]byte, error) {
	if i.Valid {
		return json.Marshal(i.Value)
	}

	return json.Marshal(nil)
}

// TODO: where should we put the validation? Should we have separate validation per domain model and transports?
type User struct {
	ID       int64
	Username string
	Email    string
	Password string
	Bio      Bio
	Image    Image
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

type UserService interface {
	Register(user User) (*User, error)
	Login(user User) (*User, error)
	Get(user User) (*User, error)
}

type UserRepo interface {
	// TODO: should this return user? What if we assume this should only be a **write** command
	Create(u User) (*User, error)
	Get(e string) (*User, error)
	GetByID(id int64) (*User, error)
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
