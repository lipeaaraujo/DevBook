package posts

import "errors"

var ErrPostTitleDescriptionAuthorIDEmpty = errors.New("Post title, description or authorId can't be empty")
var ErrPostTitleTooLong = errors.New("Post title can't be bigger than 50 characters")
var ErrPostDescriptionTooLong = errors.New("Post description can't be bigger than 2000 characters")
var ErrPostTitleDescriptionEmpty = errors.New("Post title or description can't be empty")
var ErrPostIDRequired = errors.New("You need to pass the postId")
var ErrAuthorIDRequired = errors.New("You need to pass the authorId")
var ErrInvalidUserID = errors.New("Invalid userId")
var ErrPostIDRequiredForDelete = errors.New("You have to pass the postId")
var ErrPostNotFound = errors.New("Post not found")
