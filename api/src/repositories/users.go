package repositories

import (
	"api/src/models"
	"database/sql"
	"fmt"
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

func (repo users) Get(nameQuery string) ([]models.User, error) {
	nameQuery = fmt.Sprintf("%%%s%%", nameQuery)

	rows, err := repo.db.Query(
		"select id, name, nickname, email, created_at from users where name ILIKE $1 or nickname ILIKE $2",
		nameQuery, nameQuery,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Nickname,
			&user.Email,
			&user.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
