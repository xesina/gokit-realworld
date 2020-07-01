package go_kit_realworld_example_app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_CheckPassword(t *testing.T) {
	plain := "password"
	u := User{}
	password, err := u.HashPassword(plain)
	u.Password = password
	assert.NoError(t, err)
	assert.True(t, u.CheckPassword(plain))
}
