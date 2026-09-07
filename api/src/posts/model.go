package posts

import (
	"time"
)

type Post struct {
	Id             string     `json:"id,omitempty"`
	Title          string     `json:"title,omitempty"`
	Description    string     `json:"description,omitempty"`
	AuthorId       string     `json:"authorId,omitempty"`
	AuthorName     string     `json:"authorName,omitempty"`
	AuthorNickname string     `json:"authorNickname,omitempty"`
	CreatedAt      *time.Time `json:"createdAt,omitempty"`
	UpdatedAt      *time.Time `json:"updatedAt,omitempty"`
}

func (post *Post) PrepareCreate() error {
	if post.Title == "" || post.Description == "" || post.AuthorId == "" {
		return ErrPostTitleDescriptionAuthorIDEmpty
	}

	if len(post.Title) > 50 {
		return ErrPostTitleTooLong
	}

	if len(post.Description) > 2000 {
		return ErrPostDescriptionTooLong
	}

	return nil
}

func (post *Post) PrepareUpdate() error {
	if post.Title == "" || post.Description == "" {
		return ErrPostTitleDescriptionEmpty
	}

	if len(post.Title) > 50 {
		return ErrPostTitleTooLong
	}

	if len(post.Description) > 2000 {
		return ErrPostDescriptionTooLong
	}

	return nil
}
