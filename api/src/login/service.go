package login

import (
	"api/src/apierrors"
	"api/src/users"
	"api/src/utils"
	"api/src/utils/auth"
)

type LoginService struct{
	userRepo users.UserRepoInterface
}

func NewLoginService(userRepo users.UserRepoInterface) *LoginService {
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
		return "", apierrors.Unauthorized("Invalid credentials")
	}

	token, err := auth.CreateToken(existingUser.ID)
	if err != nil {
		return "", err
	}

	return token, err
}
