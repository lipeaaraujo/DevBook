package models

import "time"

type User struct {
	ID string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Email string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}
