package go_kit_realworld_example_app

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type Bio struct {
	Value string
	Valid bool
}

type Image struct {
	Value string
	Valid bool
}

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
	Register(user User) error
	Get(email string) (*User, error)
}

// ==========================================================================================
// Store
// ==========================================================================================

type UserRepo interface {
	// TODO: should this return user? What if we assume this should only be a **write** command
	Create(u User) (*User, error)
	Get(e string) (*User, error)
}
