package users

import (
	"strings"

	"github.com/google/uuid"
)

type FakeUserRepo struct {
	saved map[string]User
	err   error
}

func newFakeRepo() *FakeUserRepo {
	return &FakeUserRepo{saved: map[string]User{}}
}

func (r *FakeUserRepo) Create(user *User) (string, error) {
	user.ID = uuid.New().String()
	r.saved[user.ID] = *user
	return user.ID, nil
}

func (r *FakeUserRepo) Get(name string) ([]User, error) {
	users := []User{}

	for _, user := range r.saved {
		if strings.Contains(user.Name, name) {
			users = append(users, user)
		}
	}

	return users, nil
}

func (r *FakeUserRepo) GetById(userId string, viewerId string) (User, error) {
	panic("unimplemented")
}

func (r *FakeUserRepo) GetByEmail(email string) (User, error) {
	for _, u := range r.saved {
		if u.Email == email {
			return u, nil
		}
	}

	return User{}, nil
}

// Delete implements [UserRepoInterface].
func (r *FakeUserRepo) Delete(userId string) error {
	panic("unimplemented")
}

// Follow implements [UserRepoInterface].
func (r *FakeUserRepo) Follow(userId string, followId string) error {
	panic("unimplemented")
}

// GetPwd implements [UserRepoInterface].
func (r *FakeUserRepo) GetPwd(userId string) (string, error) {
	panic("unimplemented")
}

// Unfollow implements [UserRepoInterface].
func (r *FakeUserRepo) Unfollow(userId string, unfollowId string) error {
	panic("unimplemented")
}

// Update implements [UserRepoInterface].
func (r *FakeUserRepo) Update(userId string, user *User) error {
	panic("unimplemented")
}

// UpdatePwd implements [UserRepoInterface].
func (r *FakeUserRepo) UpdatePwd(userId string, newPwd string) error {
	panic("unimplemented")
}
