package posts

import "database/sql"

type PostRepo struct{
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

func (repo PostRepo) Get() ([]Post, error) {
	return []Post{}, nil
}
