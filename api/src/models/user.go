package models

import (
	"errors"
	"strings"
	"time"

	"github.com/badoux/checkmail"
)

type User struct {
	ID string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Email string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (user *User) Prepare(isUpdate bool) error {
	user.format()
	if err := user.validate(isUpdate); err != nil {
		return err
	}
	return nil
}

func (user *User) validate(isUpdate bool) error {
	if user.Name == "" {
		return errors.New("User name can't be null or empty")
	}

	if user.Nickname == "" {
		return errors.New("User nickname can't be null or empty")
	}

	if user.Email == "" {
		return errors.New("User email can't be null or empty")
	}

	if err := checkmail.ValidateFormat(user.Email); err != nil {
		return errors.New("Invalid email format")
	}

	if !isUpdate && user.Password == "" {
		return errors.New("User password can't be null or empty")
	}

	return nil
}

func (user *User) format() {
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)
	user.Nickname = strings.TrimSpace(user.Nickname)
}
