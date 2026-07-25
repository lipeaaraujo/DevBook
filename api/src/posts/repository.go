package posts

import (
	"api/src/apierrors"
	"database/sql"
	"fmt"
	"time"
)

type PostRepo struct {
	db *sql.DB
}

func NewPostRepo(db *sql.DB) *PostRepo {
	return &PostRepo{db: db}
}

func (repo PostRepo) Create(post Post) (string, error) {
	statement, err := repo.db.Prepare(
		"insert into posts (title, description, author_id) values ($1, $2, $3) returning id",
	)
	if err != nil {
		return "", err
	}
	defer statement.Close()

	var insertedId string
	err = statement.QueryRow(post.Title, post.Description, post.AuthorId).Scan(&insertedId)
	if err != nil {
		return "", err
	}

	return insertedId, nil
}

// Every read returns the author's name and nickname, so clients don't have to
// resolve author_id into a person themselves.
const selectPost = `select p.id, p.title, p.description, p.author_id, u.name, u.nickname, p.created_at, p.updated_at
	from posts p inner join users u on u.id = p.author_id`

func scanPost(rows *sql.Rows, post *Post) error {
	return rows.Scan(
		&post.Id,
		&post.Title,
		&post.Description,
		&post.AuthorId,
		&post.AuthorName,
		&post.AuthorNickname,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
}

func scanPosts(rows *sql.Rows) ([]Post, error) {
	posts := []Post{}
	for rows.Next() {
		var post Post
		if err := scanPost(rows, &post); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, rows.Err()
}

func (repo PostRepo) GetById(id string) (Post, error) {
	rows, err := repo.db.Query(selectPost+" where p.id = $1", id)
	if err != nil {
		return Post{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		return Post{}, apierrors.NotFound("Post")
	}

	var post Post
	if err := scanPost(rows, &post); err != nil {
		return Post{}, err
	}

	return post, nil
}

func (repo PostRepo) Get(title string) ([]Post, error) {
	titleQuery := fmt.Sprintf("%%%s%%", title)

	rows, err := repo.db.Query(selectPost+" where p.title ILIKE $1", titleQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPosts(rows)
}

func (repo PostRepo) GetByAuthor(authorId string, title string) ([]Post, error) {
	titleQuery := fmt.Sprintf("%%%s%%", title)

	rows, err := repo.db.Query(
		selectPost+` where p.author_id = $1 and p.title ILIKE $2
			order by p.created_at desc`,
		authorId,
		titleQuery,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPosts(rows)
}

func (repo PostRepo) GetFromFollowers(userId string) ([]Post, error) {
	rows, err := repo.db.Query(
		selectPost+` inner join followers f on f.user_id = p.author_id
		 where f.follower_id = $1
		 order by p.created_at desc`,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPosts(rows)
}

func (repo PostRepo) Update(post *Post) error {
	statement, err := repo.db.Prepare(
		`update posts set title = $1, description = $2, updated_at = $3 where id = $4`,
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	updatedAt := time.Now()

	if _, err := statement.Exec(
		post.Title,
		post.Description,
		updatedAt,
		post.Id,
	); err != nil {
		return err
	}

	return nil
}

func (repo PostRepo) Delete(postId string) error {
	statement, err := repo.db.Prepare(`delete from posts where id = $1`)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err := statement.Exec(postId); err != nil {
		return err
	}

	return nil
}
