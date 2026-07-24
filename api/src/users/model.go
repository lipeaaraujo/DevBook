package users

import (
	"api/src/apierrors"
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
		return apierrors.ValidationError("User name can't be null or empty")
	}

	if user.Nickname == "" {
		return apierrors.ValidationError("User nickname can't be null or empty")
	}

	if user.Email == "" {
		return apierrors.ValidationError("User email can't be null or empty")
	}

	if err := checkmail.ValidateFormat(user.Email); err != nil {
		return apierrors.ValidationError("Invalid email format")
	}

	if !isUpdating && user.Password == "" {
		return apierrors.ValidationError("User password can't be null or empty")
	}

	return nil
}

func (user *User) PrepareLogin() error {
	user.Email = strings.TrimSpace(user.Email)

	if user.Email == "" {
		return apierrors.ValidationError("User email is required")
	}

	if err := checkmail.ValidateFormat(user.Email); err != nil {
		return apierrors.ValidationError("Invalid email format")
	}

	if user.Password == "" {
		return apierrors.ValidationError("User password is required")
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
			return apierrors.Internal(err)
		}

		user.Password = string(hashPassword)
	}

	return nil
}
