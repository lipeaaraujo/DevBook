package users_test

import (
	"api/src/users"
	"fmt"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func createDefaultUser() users.User {
	return users.User{
		Name:     "test",
		Nickname: "test_test",
		Email:    "tt@tt.com",
		Password: "12345",
	}
}

func TestPrepare(t *testing.T) {
	t.Run("trim spaces from email, name and nickname", func(t *testing.T) {
		tests := []struct {
			email            string
			name             string
			nickname         string
			expectedEmail    string
			expectedName     string
			expectedNickname string
		}{
			{" tt@tt.com", " test", " test_test", "tt@tt.com", "test", "test_test"},
			{"tt@tt.com ", "test ", "test_test ", "tt@tt.com", "test", "test_test"},
			{"  tt@tt.com   ", "  test   ", "  test_test   ", "tt@tt.com", "test", "test_test"},
			{"tt@  tt.com", "te  st", "test_  test", "tt@  tt.com", "te  st", "test_  test"},
			{"tt@tt.com", "test", "test_test", "tt@tt.com", "test", "test_test"},
		}

		for _, tt := range tests {
			t.Run(
				fmt.Sprintf("email: %s, name: %s, password: %s",
					tt.email, tt.name, tt.nickname,
				), func(t *testing.T) {
					user := createDefaultUser()
					user.Email = tt.email
					user.Name = tt.name
					user.Nickname = tt.nickname

					err := user.Prepare(false)
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}

					if user.Email != tt.expectedEmail {
						t.Fatalf("expected '%s', but got '%s'", tt.expectedEmail, user.Email)
					}
					if user.Name != tt.expectedName {
						t.Fatalf("expected '%s', but got '%s'", tt.expectedName, user.Name)
					}
					if user.Nickname != tt.expectedNickname {
						t.Fatalf("expected '%s', but got '%s'", tt.expectedNickname, user.Nickname)
					}
				})
		}
	})

	t.Run("hash user password", func(t *testing.T) {
		user := createDefaultUser()
		user.Password = "12345"

		user.Prepare(false)
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("12345"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("null, invalid and valid params", func(t *testing.T) {
		tests := []struct {
			email      string
			name       string
			nickname   string
			password   string
			isUpdating bool
			error      bool
		}{
			{"", "test", "test_test", "12345", false, true},
			{"tt@tt.com", "", "test_test", "12345", false, true},
			{"tt@tt.com", "test", "", "12345", false, true},
			{"tt@tt.com", "test", "test_test", "", false, true},
			{"tt@", "test", "test_test", "12345", false, true},
			{"tt", "test", "test_test", "12345", false, true},
			{"tt@tt.com", "test", "test_test", "12345", false, false},
			{"tt@tt.com", "test", "test_test", "", true, false},
		}

		for _, tt := range tests {
			t.Run(fmt.Sprintf(
				"email: %s, name: %s, nickname: %s, password: %s",
				tt.email, tt.name, tt.nickname, tt.password,
			), func(t *testing.T) {
				user := users.User{
					Email:    tt.email,
					Name:     tt.name,
					Nickname: tt.nickname,
					Password: tt.password,
				}

				err := user.Prepare(tt.isUpdating)
				if err == nil && tt.error {
					t.Fatalf("expected error, got nil, password: %s", user.Password)
				} else if err != nil && !tt.error {
					t.Fatalf("unexpected error: %v", err)
				}
			})
		}
	})
}

func TestPrepareLogin(t *testing.T) {
	tests := []struct {
		name           string
		email          string
		password       string
		formattedEmail string
		error          bool
	}{
		{"trim left email spaces", "  teste@teste.com", "12345", "teste@teste.com", false},
		{"trim right email spaces", "teste@teste.com   ", "12345", "teste@teste.com", false},
		{"don't trim inside email spaces", "teste@  teste.com", "12345", "teste@  teste.com", false},
		{"do nothing when theres no spaces", "teste@teste.com", "12345", "teste@teste.com", false},
		{"reject empty email", "", "12345", "", true},
		{"reject email with only spaces", "   ", "12345", "", true},
		{"reject invalid email 1", "teste", "12345", "teste", true},
		{"reject invalid email 2", "teste@", "12345", "teste@", true},
		{"reject empty password", "teste@teste.com", "", "teste@teste.com", true},
		{"accept valid email and password", "teste@teste.com", "12345", "teste@teste.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := users.User{Email: tt.email, Password: tt.password}

			err := user.PrepareLogin()
			if tt.error {
				if err == nil {
					t.Fatalf(
						"expected error for user with email: %s and password: %s, got nil",
						tt.email, tt.password,
					)
				}
				return
			}
			if err != nil {
				t.Fatalf(
					"unexpected error for email: %s and password: %s",
					tt.email, tt.password,
				)
			}
			if user.Email != tt.formattedEmail {
				t.Fatalf("expected email: '%s', but got: '%s'", tt.formattedEmail, user.Email)
			}
		})
	}
}
