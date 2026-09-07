package users

import (
	"api/src/utils"
	"errors"
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
}
