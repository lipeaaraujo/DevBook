package posts

import (
	"api/src/apierrors"
	"api/src/responses"
	"api/src/utils/auth"
	"encoding/json"
	"io"
	"net/http"
)

type PostController struct {
	service *PostService
}

func NewPostController(service *PostService) *PostController {
	return &PostController{service: service}
}

func (controller PostController) CreatePost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, apierrors.ValidationError(err.Error()))
		return 
	}

	var post Post
	if err := json.Unmarshal(body, &post); err != nil {
		responses.Error(w, apierrors.BadRequest(err.Error()))
		return
	}

	userId, err := auth.ExtractUserId(r)
	if err != nil {
		responses.Error(w, err)
		return
	}
	post.AuthorId = userId

	createdId, err := controller.service.CreatePost(post)
	if err != nil {
		responses.Error(w, err)
		return
	}

	responses.JSON(w, http.StatusCreated, map[string]string{
		"id": createdId,
	})
}

func (controller PostController) GetPosts(w http.ResponseWriter, r *http.Request) {

}

func (controller PostController) GetPostById(w http.ResponseWriter, r *http.Request) {

}

func (controller PostController) UpdatePost(w http.ResponseWriter, r *http.Request) {

}

func (controller PostController) DeletePost(w http.ResponseWriter, r *http.Request) {

}
