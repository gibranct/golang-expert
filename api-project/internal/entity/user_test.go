package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	// given
	name := "John"
	email := "j@mail.com"
	pass := "123123"

	// when
	u, err := NewUser(name, email, pass)

	// then
	assert.Nil(t, err)
	assert.Equal(t, name, u.Name)
	assert.Equal(t, email, u.Email)
	assert.NotEmpty(t, u.Password)
	assert.NotEmpty(t, u.ID)
}

func TestUser_ValidatePassword(t *testing.T) {
	// given
	name := "John"
	email := "j@mail.com"
	pass := "123123"
	wrongPass := "1231234"
	u, _ := NewUser(name, email, pass)

	// when
	r1 := u.ValidatePassword(pass)
	r2 := u.ValidatePassword(wrongPass)

	// then
	assert.True(t, r1)
	assert.False(t, r2)
	assert.NotEqual(t, pass, u.Password)
}
