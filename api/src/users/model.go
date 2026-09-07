package users

import (
	"api/src/utils"
	"strings"
	"time"

	"github.com/badoux/checkmail"
)

type User struct {
	ID        string     `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Nickname  string     `json:"nickname,omitempty"`
	Email     string     `json:"email,omitempty"`
	Password  string     `json:"password,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`

	FollowersCount int  `json:"followersCount"`
	FollowingCount int  `json:"followingCount"`
	IsFollowing    bool `json:"isFollowing"`
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
		return ErrNameEmpty
	}

	if user.Nickname == "" {
		return ErrNicknameEmpty
	}

	if user.Email == "" {
		return ErrEmailEmpty
	}

	if err := checkmail.ValidateFormat(user.Email); err != nil {
		return ErrInvalidEmailFormat
	}

	if !isUpdating && user.Password == "" {
		return ErrPasswordEmpty
	}

	return nil
}

func (user *User) PrepareLogin() error {
	user.Email = strings.TrimSpace(user.Email)

	if user.Email == "" {
		return ErrEmailEmpty
	}

	if err := checkmail.ValidateFormat(user.Email); err != nil {
		return ErrInvalidEmailFormat
	}

	if user.Password == "" {
		return ErrPasswordEmpty
	}

	return nil
}

func (user *User) format(isUpdating bool) error {
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)
	user.Nickname = strings.TrimSpace(user.Nickname)

	if !isUpdating && user.Password != "" {
		hashPassword, err := utils.Hash(user.Password)
		if err != nil {
			return err
		}

		user.Password = string(hashPassword)
	}

	return nil
}
