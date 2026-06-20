package users

import (
	"api/src/apierrors"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type UserRepoInterface interface {
	Create(user *User) (string, error)
	Get(nameQuery string) ([]User, error)
	GetById(userId string) (User, error)
	GetByEmail(email string) (User, error)
	Update(userId string, user *User) (error)
	Delete(userId string) (error)
	Follow(userId, followId string) (error)
	Unfollow(userId, unfollowId string) (error)
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo UserRepository) Create(user *User) (string, error) {
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
		if strings.Contains(err.Error(), "duplicate key") {
			err = apierrors.ResourceAlreadyExists("User with that email")
		}
		return "", err
	}

	return insertedId, nil
}

func (repo UserRepository) Get(nameQuery string) ([]User, error) {
	nameQuery = fmt.Sprintf("%%%s%%", nameQuery)

	rows, err := repo.db.Query(
		"select id, name, nickname, email, created_at, updated_at from users where name ILIKE $1 or nickname ILIKE $2",
		nameQuery, nameQuery,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User

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

func (repo UserRepository) GetById(userId string) (User, error) {
	rows, err := repo.db.Query(
		"select id, name, nickname, email, created_at, updated_at from users where id = $1",
		userId,
	)

	if err != nil {
		return User{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		return User{}, nil
	}

	var user User
	err = rows.Scan(
		&user.ID,
		&user.Name,
		&user.Nickname,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (repo UserRepository) GetByEmail(email string) (User, error) {
	rows, err := repo.db.Query(
		"select id, password from users where email = $1",
		email,
	)
	if err != nil {
		return User{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		return User{}, err
	}

	var user User
	err = rows.Scan(
		&user.ID,
		&user.Password,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (repo UserRepository) Update(userId string, user *User) (error) {
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

func (repo UserRepository) Delete(userId string) error {
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

func (repo UserRepository) Follow(userId, followId string) error {
	statement, err := repo.db.Prepare(
		"insert into followers (follower_id, user_id) values ($1, $2)",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err := statement.Exec(userId, followId); err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			err = apierrors.New(http.StatusBadRequest, "Already following user.", nil)
			return err
		}
		return err
	}

	return nil
}

func (repo UserRepository) Unfollow(userId, unfollowId string) error {
	statement, err := repo.db.Prepare(
		"delete from followers where follower_id = $1 and user_id = $2",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err := statement.Exec(userId, unfollowId); err != nil {
		return err
	}

	return nil
}
