package users

import (
	"api/src/utils"
	"errors"
	"strings"
	"testing"
)

func TestService_Create(t *testing.T) {
	tests := []struct {
		name      string
		user      User
		wantUser  User
		wantErr   error
		wantSaved bool
	}{
		{
			name:      "creates valid user",
			user:      User{Name: "Teste", Nickname: "teste", Email: "teste@teste.com", Password: "teste"},
			wantUser:  User{Name: "Teste", Nickname: "teste", Email: "teste@teste.com", Password: "teste"},
			wantErr:   nil,
			wantSaved: true,
		},
		{
			name:    "rejects empty name",
			user:    User{Name: "", Nickname: "teste", Email: "teste@teste.com", Password: "teste"},
			wantErr: ErrNameEmpty,
		},
		{
			name:    "rejects empty nickname",
			user:    User{Name: "Teste", Nickname: "", Email: "teste@teste.com", Password: "teste"},
			wantErr: ErrNicknameEmpty,
		},
		{
			name:    "rejects empty email",
			user:    User{Name: "Teste", Nickname: "teste", Email: "", Password: "teste"},
			wantErr: ErrEmailEmpty,
		},
		{
			name:    "rejects empty password",
			user:    User{Name: "Teste", Nickname: "teste", Email: "teste@teste.com", Password: ""},
			wantErr: ErrPasswordEmpty,
		},
		{
			name:    "rejects invalid email",
			user:    User{Name: "Teste", Nickname: "teste", Email: "teste", Password: "teste"},
			wantErr: ErrInvalidEmailFormat,
		},
		{
			name:      "trims empty space",
			user:      User{Name: "   Teste Espaço  ", Nickname: "teste    ", Email: "    teste@teste.com", Password: "   teste    "},
			wantUser:  User{Name: "Teste Espaço", Nickname: "teste", Email: "teste@teste.com", Password: "   teste    "},
			wantErr:   nil,
			wantSaved: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeRepo()
			service := NewUserService(repo)
			user := tt.user

			err, savedUser := service.Create(&user)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected err: %v, got %v", tt.wantErr, err)
			}

			if tt.wantSaved {
				if len(repo.saved) != 1 {
					t.Fatalf("expected user to be wantSaved, but wantSaved len is %d", len(repo.saved))
				}

				if tt.wantUser.Name != savedUser.Name {
					t.Fatalf("expected wantSaved name to be: %s, but found: %s", tt.wantUser.Name, savedUser.Name)
				}

				if tt.wantUser.Nickname != savedUser.Nickname {
					t.Fatalf("expected wantSaved nickname to be: %s, but found: %s", tt.wantUser.Nickname, savedUser.Nickname)
				}

				if tt.wantUser.Email != savedUser.Email {
					t.Fatalf("expected wantSaved email to be: %s, but found: %s", tt.wantUser.Email, savedUser.Email)
				}

				if err := utils.VerifyHash(tt.wantUser.Password, savedUser.Password); err != nil {
					t.Fatalf("expected hash to match original password: %s", tt.wantUser.Password)
				}
			}

			if !tt.wantSaved && len(repo.saved) != 0 {
				t.Fatalf("expected user to not be wantSaved, but wantSaved len is %d", len(repo.saved))
			}
		})
	}

	t.Run("returns repository errors", func(t *testing.T) {
		validUser := User{Name: "User", Nickname: "user", Email: "user@test.com", Password: "password"}

		for _, tt := range []struct {
			name string
			repo *FakeUserRepo
		}{
			{name: "get by email", repo: &FakeUserRepo{getByEmailErr: errors.New("lookup failed")}},
			{name: "create", repo: &FakeUserRepo{createErr: errors.New("insert failed")}},
		} {
			t.Run(tt.name, func(t *testing.T) {
				user := validUser
				err, created := NewUserService(tt.repo).Create(&user)

				if !errors.Is(err, tt.repo.getByEmailErr) && !errors.Is(err, tt.repo.createErr) {
					t.Fatalf("expected repository error, got: %v", err)
				}
				if created != nil {
					t.Fatalf("expected no created user, got: %+v", created)
				}
			})
		}
	})

	t.Run("doesn't save two users with the same email", func(t *testing.T) {
		repo := newFakeRepo()
		service := NewUserService(repo)
		u1 := User{Name: "u1", Nickname: "u1", Email: "teste@teste.com", Password: "12345"}
		u2 := User{Name: "u2", Nickname: "u2", Email: "teste@teste.com", Password: "12345"}

		err, _ := service.Create(&u1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err, _ = service.Create(&u2)
		if !errors.Is(err, ErrUserWithEmailExists) {
			t.Fatalf("expected err: %v, but got: %v", ErrUserWithEmailExists, err)
		}

		if len(repo.saved) != 1 {
			t.Fatalf("expected one saved user, got %d", len(repo.saved))
		}
	})

	t.Run("saves two users with different emails", func(t *testing.T) {
		repo := newFakeRepo()
		service := NewUserService(repo)
		u1 := User{Name: "u1", Nickname: "u1", Email: "teste1@teste.com", Password: "12345"}
		u2 := User{Name: "u2", Nickname: "u2", Email: "teste2@teste.com", Password: "12345"}

		err, _ := service.Create(&u1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err, _ = service.Create(&u2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(repo.saved) != 2 {
			t.Fatalf("expected wantSaved lenght to be: %d, but got: %d", 2, len(repo.saved))
		}
	})
}

func TestService_Get(t *testing.T) {
	t.Run("returns empty list when no users are registered", func(t *testing.T) {
		repo := newFakeRepo()
		service := NewUserService(repo)

		users, err := service.Get("")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(users) != 0 {
			t.Fatalf("expected empty user list, got len equal to: %d", len(users))
		}
	})

	t.Run("returns all users when name query is empty", func(t *testing.T) {
		repo := newFakeRepo()
		service := NewUserService(repo)
		existing := []User{
			{Name: "User1", Email: "user1@test.com", Nickname: "u1", Password: "12345"},
			{Name: "User2", Email: "user2@test.com", Nickname: "u2", Password: "12345"},
			{Name: "User3", Email: "user3@test.com", Nickname: "u3", Password: "12345"},
		}

		for _, u := range existing {
			err, _ := service.Create(&u)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}

		users, err := service.Get("")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(users) != len(existing) {
			t.Fatalf("expected len: %d, got len equal to: %d", len(existing), len(users))
		}
	})

	t.Run("filters result by name when name query is not empty", func(t *testing.T) {
		repo := newFakeRepo()
		service := NewUserService(repo)
		existing := []User{
			{Name: "testuser1", Email: "user1@test.com", Nickname: "u1", Password: "12345"},
			{Name: "testuser2", Email: "user2@test.com", Nickname: "u2", Password: "12345"},
			{Name: "user3", Email: "user3@test.com", Nickname: "u3", Password: "12345"},
		}

		for _, u := range existing {
			err, _ := service.Create(&u)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}

		users, err := service.Get("test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(users) != 2 {
			t.Fatalf("expected len: %d, got len equal to: %d", 2, len(users))
		}
	})

	t.Run("returns repository error", func(t *testing.T) {
		wantErr := errors.New("database unavailable")
		service := NewUserService(&FakeUserRepo{getErr: wantErr})

		_, err := service.Get("")
		if !errors.Is(err, wantErr) {
			t.Fatalf("expected err: %v, got: %v", wantErr, err)
		}
	})
}

func TestService_GetByID(t *testing.T) {
	t.Run("returns user", func(t *testing.T) {
		repo := newFakeRepo()
		repo.saved["user-id"] = User{ID: "user-id", Name: "User"}

		user, err := NewUserService(repo).GetById("user-id", "viewer-id")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if user.ID != "user-id" {
			t.Fatalf("expected user ID user-id, got %s", user.ID)
		}
	})

	t.Run("returns not found when repository returns no user", func(t *testing.T) {
		user, err := NewUserService(newFakeRepo()).GetById("missing", "viewer-id")
		if !errors.Is(err, ErrUserNotFound) {
			t.Fatalf("expected err: %v, got: %v", ErrUserNotFound, err)
		}
		if user != nil {
			t.Fatalf("expected no user, got: %+v", user)
		}
	})

	t.Run("returns repository error", func(t *testing.T) {
		wantErr := errors.New("lookup failed")
		user, err := NewUserService(&FakeUserRepo{getByIDErr: wantErr}).GetById("user-id", "viewer-id")
		if !errors.Is(err, wantErr) {
			t.Fatalf("expected err: %v, got: %v", wantErr, err)
		}
		if user != nil {
			t.Fatalf("expected no user, got: %+v", user)
		}
	})
}

func TestService_Update(t *testing.T) {
	t.Run("formats and sends user to repository", func(t *testing.T) {
		repo := newFakeRepo()
		user := User{Name: " User ", Nickname: " user ", Email: " user@test.com "}

		if err := NewUserService(repo).Update("user-id", &user); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if repo.updatedUserID != "user-id" {
			t.Fatalf("expected update for user-id, got %s", repo.updatedUserID)
		}
		if repo.updatedUser.Name != "User" || repo.updatedUser.Nickname != "user" || repo.updatedUser.Email != "user@test.com" {
			t.Fatalf("expected formatted user, got: %+v", repo.updatedUser)
		}
	})

	t.Run("returns validation error without updating", func(t *testing.T) {
		repo := newFakeRepo()
		err := NewUserService(repo).Update("user-id", &User{Nickname: "user", Email: "user@test.com"})
		if !errors.Is(err, ErrNameEmpty) {
			t.Fatalf("expected err: %v, got: %v", ErrNameEmpty, err)
		}
		if repo.updatedUserID != "" {
			t.Fatalf("repository should not have been called")
		}
	})

	t.Run("returns repository error", func(t *testing.T) {
		wantErr := errors.New("update failed")
		repo := &FakeUserRepo{updateErr: wantErr}
		err := NewUserService(repo).Update("user-id", &User{Name: "User", Nickname: "user", Email: "user@test.com"})
		if !errors.Is(err, wantErr) {
			t.Fatalf("expected err: %v, got: %v", wantErr, err)
		}
	})
}

func TestService_Delete(t *testing.T) {
	t.Run("deletes user", func(t *testing.T) {
		repo := newFakeRepo()
		if err := NewUserService(repo).Delete("user-id"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if repo.deletedUserID != "user-id" {
			t.Fatalf("expected deleted user ID user-id, got %s", repo.deletedUserID)
		}
	})

	t.Run("returns repository error", func(t *testing.T) {
		wantErr := errors.New("delete failed")
		err := NewUserService(&FakeUserRepo{deleteErr: wantErr}).Delete("user-id")
		if !errors.Is(err, wantErr) {
			t.Fatalf("expected err: %v, got: %v", wantErr, err)
		}
	})
}

func TestService_Follow(t *testing.T) {
	t.Run("follows another user", func(t *testing.T) {
		repo := newFakeRepo()
		if err := NewUserService(repo).Follow("user-id", "follow-id"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if repo.followedUserID != "follow-id" {
			t.Fatalf("expected followed user ID follow-id, got %s", repo.followedUserID)
		}
	})

	t.Run("rejects following self", func(t *testing.T) {
		repo := newFakeRepo()
		err := NewUserService(repo).Follow("user-id", "user-id")
		if !errors.Is(err, ErrCannotFollowSelf) {
			t.Fatalf("expected err: %v, got: %v", ErrCannotFollowSelf, err)
		}
		if repo.followedUserID != "" {
			t.Fatalf("repository should not have been called")
		}
	})

	t.Run("returns repository error", func(t *testing.T) {
		wantErr := errors.New("follow failed")
		err := NewUserService(&FakeUserRepo{followErr: wantErr}).Follow("user-id", "follow-id")
		if !errors.Is(err, wantErr) {
			t.Fatalf("expected err: %v, got: %v", wantErr, err)
		}
	})
}

func TestService_Unfollow(t *testing.T) {
	t.Run("unfollows another user", func(t *testing.T) {
		repo := newFakeRepo()
		if err := NewUserService(repo).Unfollow("user-id", "unfollow-id"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if repo.unfollowedID != "unfollow-id" {
			t.Fatalf("expected unfollowed user ID unfollow-id, got %s", repo.unfollowedID)
		}
	})

	t.Run("rejects unfollowing self", func(t *testing.T) {
		repo := newFakeRepo()
		err := NewUserService(repo).Unfollow("user-id", "user-id")
		if !errors.Is(err, ErrCannotUnfollowSelf) {
			t.Fatalf("expected err: %v, got: %v", ErrCannotUnfollowSelf, err)
		}
		if repo.unfollowedID != "" {
			t.Fatalf("repository should not have been called")
		}
	})

	t.Run("returns repository error", func(t *testing.T) {
		wantErr := errors.New("unfollow failed")
		err := NewUserService(&FakeUserRepo{unfollowErr: wantErr}).Unfollow("user-id", "unfollow-id")
		if !errors.Is(err, wantErr) {
			t.Fatalf("expected err: %v, got: %v", wantErr, err)
		}
	})
}

func TestService_FollowLists(t *testing.T) {
	users := []User{{ID: "a"}, {ID: "b"}}
	repoErr := errors.New("lookup failed")
	followers := func(s *UserService) ([]User, error) { return s.GetFollowers("user-id") }
	following := func(s *UserService) ([]User, error) { return s.GetFollowing("user-id") }

	for _, tt := range []struct {
		name    string
		list    func(*UserService) ([]User, error)
		repo    *FakeUserRepo
		wantLen int
		wantErr error
	}{
		{name: "followers", list: followers, repo: &FakeUserRepo{followers: users}, wantLen: 2},
		{name: "following", list: following, repo: &FakeUserRepo{following: users}, wantLen: 2},
		{name: "followers error", list: followers, repo: &FakeUserRepo{followersErr: repoErr}, wantErr: repoErr},
		{name: "following error", list: following, repo: &FakeUserRepo{followingErr: repoErr}, wantErr: repoErr},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.list(NewUserService(tt.repo))
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected err: %v, got: %v", tt.wantErr, err)
			}
			if len(got) != tt.wantLen {
				t.Fatalf("expected %d users, got %d", tt.wantLen, len(got))
			}
		})
	}
}

func TestService_ChangePassword(t *testing.T) {
	t.Run("updates password", func(t *testing.T) {
		hash, err := utils.Hash("current-password")
		if err != nil {
			t.Fatalf("failed to hash password: %v", err)
		}
		repo := newFakeRepo()
		repo.passwords["user-id"] = string(hash)

		if err := NewUserService(repo).ChangePassword("user-id", "current-password", "new-password"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := utils.VerifyHash("new-password", repo.passwords["user-id"]); err != nil {
			t.Fatalf("expected stored password to match new password")
		}
	})

	t.Run("returns hash error", func(t *testing.T) {
		hash, err := utils.Hash("current-password")
		if err != nil {
			t.Fatalf("failed to hash password: %v", err)
		}
		repo := newFakeRepo()
		repo.passwords["user-id"] = string(hash)

		err = NewUserService(repo).ChangePassword("user-id", "current-password", strings.Repeat("a", 73))
		if err == nil {
			t.Fatal("expected hash error")
		}
	})

	getPwdErr := errors.New("lookup failed")
	updatePwdErr := errors.New("update failed")

	for _, tt := range []struct {
		name string
		repo *FakeUserRepo
		pwd  string
		want error
	}{
		{name: "get password error", repo: &FakeUserRepo{getPwdErr: getPwdErr}, want: getPwdErr},
		{name: "invalid current password", repo: &FakeUserRepo{passwords: map[string]string{}}, pwd: "invalid", want: ErrInvalidCredentials},
		{name: "update password error", repo: &FakeUserRepo{passwords: map[string]string{}, updatePwdErr: updatePwdErr}, pwd: "current-password", want: updatePwdErr},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if tt.pwd == "current-password" {
				hash, err := utils.Hash(tt.pwd)
				if err != nil {
					t.Fatalf("failed to hash password: %v", err)
				}
				tt.repo.passwords["user-id"] = string(hash)
			}

			err := NewUserService(tt.repo).ChangePassword("user-id", tt.pwd, "new-password")
			if !errors.Is(err, tt.want) {
				t.Fatalf("expected err: %v, got: %v", tt.want, err)
			}
		})
	}
}
