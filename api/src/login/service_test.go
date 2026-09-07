package login

import (
	"api/src/config"
	"api/src/users"
	"api/src/utils"
	"errors"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

type fakeUserLookup struct {
	user  users.User
	err   error
	email string
}

func (f *fakeUserLookup) GetByEmail(email string) (users.User, error) {
	f.email = email
	return f.user, f.err
}

func TestLoginService_Login(t *testing.T) {
	hashedPassword, err := utils.Hash("password")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	lookupErr := errors.New("database unavailable")

	tests := []struct {
		name       string
		user       users.User
		repo       *fakeUserLookup
		wantErr    error
		wantEmail  string
		wantUserID string
	}{
		{
			name:    "rejects empty email",
			user:    users.User{Password: "password"},
			wantErr: users.ErrEmailEmpty,
		},
		{
			name:    "rejects invalid email",
			user:    users.User{Email: "invalid", Password: "password"},
			wantErr: users.ErrInvalidEmailFormat,
		},
		{
			name:    "rejects empty password",
			user:    users.User{Email: "user@example.com"},
			wantErr: users.ErrPasswordEmpty,
		},
		{
			name:      "returns lookup error",
			user:      users.User{Email: "user@example.com", Password: "password"},
			repo:      &fakeUserLookup{err: lookupErr},
			wantErr:   lookupErr,
			wantEmail: "user@example.com",
		},
		{
			name:      "rejects unknown user credentials",
			user:      users.User{Email: "user@example.com", Password: "password"},
			repo:      &fakeUserLookup{},
			wantErr:   users.ErrInvalidCredentials,
			wantEmail: "user@example.com",
		},
		{
			name:      "rejects incorrect password",
			user:      users.User{Email: "user@example.com", Password: "incorrect"},
			repo:      &fakeUserLookup{user: users.User{ID: "user-id", Password: string(hashedPassword)}},
			wantErr:   users.ErrInvalidCredentials,
			wantEmail: "user@example.com",
		},
		{
			name:       "creates a token for valid credentials",
			user:       users.User{Email: " user@example.com ", Password: "password"},
			repo:       &fakeUserLookup{user: users.User{ID: "user-id", Password: string(hashedPassword)}},
			wantEmail:  "user@example.com",
			wantUserID: "user-id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.repo
			if repo == nil {
				repo = &fakeUserLookup{}
			}

			token, err := NewLoginService(repo).Login(&tt.user)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
			if tt.wantErr != nil && token != "" {
				t.Fatalf("expected no token, got %q", token)
			}
			if tt.wantErr == nil && token == "" {
				t.Fatal("expected token")
			}
			if tt.wantUserID != "" {
				parsed, err := jwt.Parse(token, func(*jwt.Token) (any, error) {
					return []byte(config.JwtTokenSecret), nil
				})
				if err != nil || !parsed.Valid {
					t.Fatalf("expected valid token, got %v", err)
				}
				if parsed.Claims.(jwt.MapClaims)["userId"] != tt.wantUserID {
					t.Fatalf("expected token for %q", tt.wantUserID)
				}
			}
			if repo.email != tt.wantEmail {
				t.Fatalf("expected lookup email %q, got %q", tt.wantEmail, repo.email)
			}
		})
	}
}
