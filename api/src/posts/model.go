package posts

import (
	"api/src/apierrors"
	"time"
)

type Post struct {
	Id string `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	AuthorId string `json:"authorId,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

func (post *Post) PrepareCreate() error {
	if post.Title == "" || post.Description == "" || post.AuthorId == "" {
		return apierrors.BadRequest("Post title, description or authorId can't be empty")
	}

	if len(post.Title) > 50 {
		return apierrors.BadRequest("Post title can't be bigger than 50 characters")
	}

	if len(post.Description) > 2000 {
		return apierrors.BadRequest("Post description can't be bigger than 2000 characters")
	}

	return nil
}
