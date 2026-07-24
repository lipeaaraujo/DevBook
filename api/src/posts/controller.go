package posts

import (
	"api/src/apierrors"
	"api/src/responses"
	"api/src/utils/auth"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
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
	titleQuery := strings.ToLower(r.URL.Query().Get("title"))

	posts, err := controller.service.GetPosts(titleQuery)
	if err != nil {
		responses.Error(w, err)
		return
	}

	responses.JSON(w, http.StatusOK, posts)
}

func (controller PostController) GetPostById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postId := params["postId"]

	post, err := controller.service.GetById(postId)
	if err != nil {
		responses.Error(w, err)
		return
	}

	responses.JSON(w, http.StatusOK, post)
}

func (controller PostController) GetByAuthor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	authorId := params["authorId"]
	title := strings.ToLower(r.URL.Query().Get("title"))

	posts, err := controller.service.GetByAuthor(authorId, title)
	if err != nil {
		responses.Error(w, err)
		return
	}

	responses.JSON(w, http.StatusOK, posts)
}

func (controller PostController) UpdatePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postId := params["postId"]

	userId, err := auth.ExtractUserId(r)
	if err != nil {
		responses.Error(w, err)
		return
	}

	existingPost, err := controller.service.GetById(postId)
	if err != nil {
		responses.Error(w, err)
		return
	}

	if userId != existingPost.AuthorId {
		fmt.Println(userId, existingPost.AuthorId)
		responses.Error(w, apierrors.Forbidden("Can't edit a post that isn't yours"))
		return
	}

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
	post.Id = postId

	if err := controller.service.UpdatePost(&post); err != nil {
		responses.Error(w, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func (controller PostController) DeletePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postId := params["postId"]

	userId, err := auth.ExtractUserId(r)
	if err != nil {
		responses.Error(w, err)
		return
	}

	existingPost, err := controller.service.GetById(postId)
	if err != nil {
		responses.Error(w, err)
		return
	}

	if userId != existingPost.AuthorId {
		responses.Error(w, apierrors.Forbidden("Can't delete a post that isn't yours"))
		return
	}

	if err := controller.service.DeletePost(postId); err != nil {
		responses.Error(w, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}
