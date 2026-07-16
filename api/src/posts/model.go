package posts

import "time"

type Post struct {
	Id string `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
	Description string `json:"description.omitempty"`
	AuthorId string `json:"authorId,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt"`
}
