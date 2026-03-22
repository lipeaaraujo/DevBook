package repositories

import (
	"api/src/models"
	"database/sql"
	"fmt"
	"time"
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
		"select id, name, nickname, email, created_at, updated_at from users where name ILIKE $1 or nickname ILIKE $2",
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
			&user.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (repo users) GetById(userId string) (models.User, error) {
	rows, err := repo.db.Query(
		"select id, name, nickname, email, created_at, updated_at from users where id = $1",
		userId,
	)

	if err != nil {
		return models.User{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		return models.User{}, nil
	}

	var user models.User
	err = rows.Scan(
		&user.ID,
		&user.Name,
		&user.Nickname,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (repo users) UpdateUser(userId string, user models.User) (error) {
	statement, err := repo.db.Prepare(
		"update users set name = $1, nickname = $2, email = $3, updated_at = $4 where id = $5",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	updatedAt := time.Now()

	_, err = statement.Exec(user.Name, user.Nickname, user.Email, updatedAt, userId)
	if err != nil {
		return err
	}

	return nil
}

func (repo users) Delete(userId string) error {
	statement, err := repo.db.Prepare("delete from users where id = $1")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(userId); err != nil {
		return err
	}

	return nil
}
