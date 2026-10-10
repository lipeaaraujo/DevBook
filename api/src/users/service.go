package users

import "api/src/utils"

type UserRepoInterface interface {
	Create(user *User) (string, error)
	Get(nameQuery string) ([]User, error)
	GetById(userId, viewerId string) (User, error)
	GetByEmail(email string) (User, error)
	Update(userId string, user *User) error
	Delete(userId string) error
	Follow(userId, followId string) error
	Unfollow(userId, unfollowId string) error
	GetFollowers(userId string) ([]User, error)
	GetFollowing(userId string) ([]User, error)
	GetPwd(userId string) (string, error)
	UpdatePwd(userId, newPwd string) error
}

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

	existingUser, err := service.repo.GetByEmail(user.Email)
	if err != nil {
		return err, nil
	}

	if existingUser.ID != "" {
		return ErrUserWithEmailExists, nil
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
		return nil, ErrUserNotFound
	}

	if id != viewerId {
		user.Email = ""
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
		return ErrCannotFollowSelf
	}

	if err := service.repo.Follow(userId, followId); err != nil {
		return err
	}

	return nil
}

func (service UserService) Unfollow(userId, followId string) error {
	if userId == followId {
		return ErrCannotUnfollowSelf
	}

	if err := service.repo.Unfollow(userId, followId); err != nil {
		return err
	}

	return nil
}

func (service UserService) GetFollowers(userId string) ([]User, error) {
	return service.repo.GetFollowers(userId)
}

func (service UserService) GetFollowing(userId string) ([]User, error) {
	return service.repo.GetFollowing(userId)
}

func (service UserService) ChangePassword(userId, currentPwd, newPwd string) error {
	// check currentPwd
	userPwd, err := service.repo.GetPwd(userId)
	if err != nil {
		return err
	}

	err = utils.VerifyHash(currentPwd, userPwd)
	if err != nil {
		return ErrInvalidCredentials
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
