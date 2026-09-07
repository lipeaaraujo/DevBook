package users

import (
	"errors"
	"testing"
)

func TestService_Create(t *testing.T) {
	tests := []struct {
		name     string
		username string
		nickname string
		email    string
		password string
		err      error
		saved    bool
	}{
		{name: "creates valid user", username: "Teste", nickname: "teste", email: "teste@teste.com", password: "teste", err: nil, saved: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeRepo()
			service := NewUserService(repo)
			user := User{
				Name:     tt.username,
				Nickname: tt.nickname,
				Password: tt.password,
				Email:    tt.email,
			}

			err, _ := service.Create(&user)
			if !errors.Is(tt.err, err) {
				t.Fatalf("expected err: %v, got %v", tt.err, err)
			}

			if tt.saved && len(repo.saved) != 1 {
				t.Fatalf("expected user to be saved, but saved len is %d", len(repo.saved))
			}

			if !tt.saved && len(repo.saved) != 0 {
				t.Fatalf("expected user to not be saved, but saved len is %d", len(repo.saved))
			}
		})
	}
}
