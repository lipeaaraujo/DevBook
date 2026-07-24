package posts

import (
	"api/src/router/routes"
	"net/http"
)

func CreatePostRoutes(controller *PostController) []routes.Route {
	return []routes.Route{
		{
			URI:          "/post",
			Method:       http.MethodPost,
			Handler:      controller.CreatePost,
			RequiresAuth: true,
		},
		{
			URI:          "/post",
			Method:       http.MethodGet,
			Handler:      controller.GetPosts,
			RequiresAuth: true,
		},
		{
			URI:          "/post/{postId}",
			Method:       http.MethodGet,
			Handler:      controller.GetPostById,
			RequiresAuth: true,
		},
		{
			URI:          "/post/user/{authorId}",
			Method:       http.MethodGet,
			Handler:      controller.GetByAuthor,
			RequiresAuth: true,
		},
		{
			URI:          "/post/{postId}",
			Method:       http.MethodPut,
			Handler:      controller.UpdatePost,
			RequiresAuth: true,
		},
		{
			URI:          "/post/{postId}",
			Method:       http.MethodDelete,
			Handler:      controller.DeletePost,
			RequiresAuth: true,
		},
	}
}
