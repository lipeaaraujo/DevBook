package models

import (
	"api/src/utils"
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

func (user *User) Prepare(isUpdating bool) error {
	if err := user.format(isUpdating); err != nil {
		return err
	}
	if err := user.validate(isUpdating); err != nil {
		return err
	}
	return nil
}

func (user *User) validate(isUpdating bool) error {
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

	if !isUpdating && user.Password == "" {
		return errors.New("User password can't be null or empty")
	}

	return nil
}

func (user *User) format(isUpdating bool) error {
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)
	user.Nickname = strings.TrimSpace(user.Nickname)

	if (!isUpdating) {
		hashPassword, err := utils.Hash(user.Password)
		if err != nil {
			return err
		}	

		user.Password = string(hashPassword)
	}

	return nil
}
