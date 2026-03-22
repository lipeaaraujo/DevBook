package repositories

import (
	"api/src/models"
	"database/sql"
)

type users struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *users {
	return &users{db}
}

func (repo users) Create(user models.User) (string, error) {
	statement, err := repo.db.Prepare(
		"insert into users (name, nickname, email, password) values($1, $2, $3, $4) returning id",
	)
	if err != nil {
		return "", err
	}
	defer statement.Close()

	var insertedId string
	err = statement.QueryRow(user.Name, user.Nickname, user.Email, user.Password).Scan(&insertedId)
	if err != nil {
		return "", err
	}

	return insertedId, nil
}
