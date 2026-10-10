package users

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

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
			err = ErrUserWithEmailExists
		}
		return "", err
	}

	return insertedId, nil
}

func (repo UserRepository) Get(nameQuery string) ([]User, error) {
	nameQuery = fmt.Sprintf("%%%s%%", nameQuery)

	rows, err := repo.db.Query(
		"select id, name, nickname, created_at, updated_at from users where name ILIKE $1 or nickname ILIKE $2",
		nameQuery, nameQuery,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User

		if err := rows.Err(); err != nil {
			return nil, err
		}

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Nickname,
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

func (repo UserRepository) GetById(userId, viewerId string) (User, error) {
	rows, err := repo.db.Query(
		`select u.id, u.name, u.nickname, u.email, u.created_at, u.updated_at,
			(select count(*) from followers where user_id = u.id) as followers_count,
			(select count(*) from followers where follower_id = u.id) as following_count,
			exists(select 1 from followers where user_id = u.id and follower_id = $2) as is_following
		 from users u where u.id = $1`,
		userId,
		viewerId,
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
		&user.FollowersCount,
		&user.FollowingCount,
		&user.IsFollowing,
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

func (repo UserRepository) Update(userId string, user *User) error {
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
			err = ErrAlreadyFollowingUser
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

func (repo UserRepository) GetFollowers(userId string) ([]User, error) {
	return repo.listUsers(
		`select u.id, u.name, u.nickname from users u
		 inner join followers f on f.follower_id = u.id
		 where f.user_id = $1 order by u.name`,
		userId,
	)
}

func (repo UserRepository) GetFollowing(userId string) ([]User, error) {
	return repo.listUsers(
		`select u.id, u.name, u.nickname from users u
		 inner join followers f on f.user_id = u.id
		 where f.follower_id = $1 order by u.name`,
		userId,
	)
}

func (repo UserRepository) listUsers(query string, args ...any) ([]User, error) {
	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Nickname); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func (repo UserRepository) GetPwd(userId string) (string, error) {
	rows, err := repo.db.Query(
		"select password from users where id = $1",
		userId,
	)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var password string
	if !rows.Next() {
		return "", ErrUserNotFound
	}

	if err := rows.Scan(&password); err != nil {
		return "", err
	}

	return password, nil
}

func (repo UserRepository) UpdatePwd(userId, hashedPassword string) error {
	statement, err := repo.db.Prepare(
		"update users set password = $1 where id = $2",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err := statement.Exec(hashedPassword, userId); err != nil {
		return err
	}

	return nil
}
