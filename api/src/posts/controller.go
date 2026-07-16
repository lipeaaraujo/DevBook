package posts

import "net/http"

type PostController struct {
	service *PostService
}

func NewPostController(service *PostService) *PostController {
	return &PostController{service: service}
}

func (controller PostController) CreatePost(w http.ResponseWriter, r *http.Request) {

}

func (controller PostController) GetPosts(w http.ResponseWriter, r *http.Request) {

}

func (controller PostController) GetPostById(w http.ResponseWriter, r *http.Request) {

}

func (controller PostController) UpdatePost(w http.ResponseWriter, r *http.Request) {

}

func (controller PostController) DeletePost(w http.ResponseWriter, r *http.Request) {

}
