package login

import (
	"api/src/users"
	"api/src/utils"
	"api/src/utils/auth"
)

type UserLookup interface {
	GetByEmail(email string) (users.User, error)
}

type LoginService struct {
	userRepo UserLookup
}

func NewLoginService(userRepo UserLookup) *LoginService {
	return &LoginService{userRepo: userRepo}
}

func (service LoginService) Login(user *users.User) (string, error) {
	if err := user.PrepareLogin(); err != nil {
		return "", err
	}

	existingUser, err := service.userRepo.GetByEmail(user.Email)
	if err != nil {
		return "", err
	}

	err = utils.VerifyHash(user.Password, existingUser.Password)
	if err != nil {
		return "", users.ErrInvalidCredentials
	}

	token, err := auth.CreateToken(existingUser.ID)
	if err != nil {
		return "", err
	}

	return token, err
}
