package users

import (
	"strings"

	"github.com/google/uuid"
)

type FakeUserRepo struct {
	saved          map[string]User
	passwords      map[string]string
	createErr      error
	getErr         error
	getByIDErr     error
	getByEmailErr  error
	updateErr      error
	deleteErr      error
	followErr      error
	unfollowErr    error
	getPwdErr      error
	updatePwdErr   error
	updatedUser    User
	updatedUserID  string
	deletedUserID  string
	followedUserID string
	unfollowedID   string
}

func newFakeRepo() *FakeUserRepo {
	return &FakeUserRepo{
		saved:     map[string]User{},
		passwords: map[string]string{},
	}
}

func (r *FakeUserRepo) Create(user *User) (string, error) {
	if r.createErr != nil {
		return "", r.createErr
	}

	user.ID = uuid.New().String()
	r.saved[user.ID] = *user
	return user.ID, nil
}

func (r *FakeUserRepo) Get(name string) ([]User, error) {
	if r.getErr != nil {
		return nil, r.getErr
	}

	users := []User{}
	for _, user := range r.saved {
		if name == "" || strings.Contains(strings.ToLower(user.Name), strings.ToLower(name)) {
			users = append(users, user)
		}
	}

	return users, nil
}

func (r *FakeUserRepo) GetById(userId string, viewerId string) (User, error) {
	if r.getByIDErr != nil {
		return User{}, r.getByIDErr
	}

	return r.saved[userId], nil
}

func (r *FakeUserRepo) GetByEmail(email string) (User, error) {
	if r.getByEmailErr != nil {
		return User{}, r.getByEmailErr
	}

	for _, u := range r.saved {
		if u.Email == email {
			return u, nil
		}
	}

	return User{}, nil
}

func (r *FakeUserRepo) Update(userId string, user *User) error {
	if r.updateErr != nil {
		return r.updateErr
	}

	r.updatedUserID = userId
	r.updatedUser = *user
	return nil
}

func (r *FakeUserRepo) Delete(userId string) error {
	if r.deleteErr != nil {
		return r.deleteErr
	}

	r.deletedUserID = userId
	delete(r.saved, userId)
	return nil
}

func (r *FakeUserRepo) Follow(userId string, followId string) error {
	if r.followErr != nil {
		return r.followErr
	}

	r.followedUserID = followId
	return nil
}

func (r *FakeUserRepo) Unfollow(userId string, unfollowId string) error {
	if r.unfollowErr != nil {
		return r.unfollowErr
	}

	r.unfollowedID = unfollowId
	return nil
}

func (r *FakeUserRepo) GetPwd(userId string) (string, error) {
	if r.getPwdErr != nil {
		return "", r.getPwdErr
	}

	return r.passwords[userId], nil
}

func (r *FakeUserRepo) UpdatePwd(userId string, newPwd string) error {
	if r.updatePwdErr != nil {
		return r.updatePwdErr
	}

	r.passwords[userId] = newPwd
	return nil
}
