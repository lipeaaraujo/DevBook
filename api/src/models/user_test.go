package models_test

import (
	"api/src/models"
	"testing"
)

func TestPrepareLogin(t *testing.T) {
	tests := []struct{
		name string
		email string
		password string
		formattedEmail string
		error bool
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
			user := models.User{Email: tt.email, Password: tt.password}

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
			if user.Email != tt.formattedEmail {
				t.Fatalf("expected email: '%s', but got: '%s'", tt.formattedEmail, user.Email)
			}
		})
	}
}
