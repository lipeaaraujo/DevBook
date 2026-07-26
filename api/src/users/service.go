package users

import (
	"api/src/apierrors"
	"api/src/utils"
	"net/http"
)

type UserService struct {
	repo UserRepoInterface
}

func NewUserService(repo UserRepoInterface) *UserService {
	return &UserService{repo: repo}
}

func (service UserService) Create(user *User) (error, *User) {
	if err := user.Prepare(false); err != nil {
		return err, nil
	}

	userId, err := service.repo.Create(user)
	if err != nil {
		return err, nil
	}
	user.ID = userId

	return nil, user
}

func (service UserService) Get(name string) ([]User, error) {
	users, err := service.repo.Get(name)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (service UserService) GetById(id, viewerId string) (*User, error) {
	user, err := service.repo.GetById(id, viewerId)
	if err != nil {
		return nil, err
	}

	if user == (User{}) {
		return nil, apierrors.NotFound("User")
	}

	return &user, err
}

func (service UserService) Update(id string, user *User) error {
	if err := user.Prepare(true); err != nil {
		return err
	}

	if err := service.repo.Update(id, user); err != nil {
		return err
	}

	return nil
}

func (service UserService) Delete(id string) error {
	if err := service.repo.Delete(id); err != nil {
		return err
	}
	return nil
}

func (service UserService) Follow(userId, followId string) error {
	if userId == followId {
		return apierrors.New(http.StatusForbidden, "You can't follow your own user.", nil)
	}

	if err := service.repo.Follow(userId, followId); err != nil {
		return err
	}

	return nil
}

func (service UserService) Unfollow(userId, followId string) error {
	if userId == followId {
		return apierrors.New(http.StatusForbidden, "You can't unfollow your own user.", nil)
	}

	if err := service.repo.Unfollow(userId, followId); err != nil {
		return err
	}

	return nil
}

func (service UserService) ChangePassword(userId, currentPwd, newPwd string) error {
	// check currentPwd
	userPwd, err := service.repo.GetPwd(userId)
	if err != nil {
		return err
	}

	err = utils.VerifyHash(currentPwd, userPwd)
	if err != nil {
		return apierrors.Unauthorized("Invalid credentials")
	}

	// hash newPwd
	hash, err := utils.Hash(newPwd)
	if err != nil {
		return err
	}

	// save newPwd
	if err := service.repo.UpdatePwd(userId, string(hash)); err != nil {
		return err
	}

	return nil
}
