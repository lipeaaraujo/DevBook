package users

import (
	"api/src/apierrors"
	"api/src/responses"
	"api/src/utils/auth"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type UserController struct {
	service *UserService
}

func NewUserController(s *UserService) *UserController {
	return &UserController{s}
}

type changePasswordRequest struct {
	CurrentPwd string `json:"currentPassword"`
	NewPwd     string `json:"newPassword"`
}

func userErrorResponse(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrEmailEmpty),
		errors.Is(err, ErrInvalidEmailFormat),
		errors.Is(err, ErrNicknameEmpty),
		errors.Is(err, ErrPasswordEmpty),
		errors.Is(err, ErrNameEmpty):
		responses.Error(w, apierrors.ValidationError(err.Error()))

	case errors.Is(err, ErrUserWithEmailExists):
		responses.Error(w, apierrors.ResourceAlreadyExists("User with that email"))

	case errors.Is(err, ErrUserNotFound):
		responses.Error(w, apierrors.NotFound("User"))

	case errors.Is(err, ErrCannotFollowSelf), errors.Is(err, ErrCannotUnfollowSelf):
		responses.Error(w, apierrors.Forbidden(err.Error()))

	case errors.Is(err, ErrAlreadyFollowingUser):
		responses.Error(w, apierrors.BadRequest(err.Error()))

	case errors.Is(err, ErrInvalidCredentials):
		responses.Error(w, apierrors.Unauthorized(err.Error()))

	default:
		responses.Error(w, err)
	}
}

func (controller UserController) CreateUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, apierrors.ValidationError(err.Error()))
		return
	}

	var user User
	if err := json.Unmarshal(requestBody, &user); err != nil {
		responses.Error(w, apierrors.BadRequest(err.Error()))
		return
	}

	err, createdUser := controller.service.Create(&user)
	if err != nil {
		userErrorResponse(w, err)
		return
	}

	createdUser.Password = ""
	responses.JSON(w, http.StatusCreated, createdUser)
}

func (controller UserController) GetUsers(
	w http.ResponseWriter,
	r *http.Request,
) {
	nameQuery := strings.ToLower(r.URL.Query().Get("user"))

	users, err := controller.service.Get(nameQuery)
	if err != nil {
		userErrorResponse(w, err)
		return
	}

	responses.JSON(w, http.StatusOK, users)
}

func (controller UserController) GetUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	params := mux.Vars(r)
	userId := params["userId"]

	viewerId, err := auth.ExtractUserId(r)
	if err != nil {
		responses.Error(w, err)
		return
	}

	user, err := controller.service.GetById(userId, viewerId)
	if err != nil {
		userErrorResponse(w, err)
		return
	}

	responses.JSON(w, http.StatusOK, user)
}

func (controller UserController) UpdateUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	params := mux.Vars(r)
	userId := params["userId"]

	autenticatedUserId, err := auth.ExtractUserId(r)
	if err != nil {
		responses.Error(w, err)
		return
	}

	if userId != autenticatedUserId {
		responses.Error(w, apierrors.Forbidden("Can't edit a user different from yours"))
		return
	}

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, apierrors.ValidationError(err.Error()))
		return
	}

	var user User
	if err := json.Unmarshal(requestBody, &user); err != nil {
		responses.Error(w, apierrors.Forbidden(err.Error()))
		return
	}

	if err := controller.service.Update(userId, &user); err != nil {
		userErrorResponse(w, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func (controller UserController) DeleteUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	params := mux.Vars(r)
	userId := params["userId"]

	autenticatedUserId, err := auth.ExtractUserId(r)
	if err != nil {
		responses.Error(w, err)
		return
	}

	if userId != autenticatedUserId {
		responses.Error(w, apierrors.Forbidden("Can't delete a different user from yours"))
		return
	}

	if err := controller.service.Delete(userId); err != nil {
		userErrorResponse(w, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func (controller UserController) FollowUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	params := mux.Vars(r)
	userToFollowId := params["userId"]

	authenticatedUserId, err := auth.ExtractUserId(r)
	if err != nil {
		responses.Error(w, err)
		return
	}

	if err = controller.service.Follow(authenticatedUserId, userToFollowId); err != nil {
		userErrorResponse(w, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func (controller UserController) UnfollowUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	params := mux.Vars(r)
	unfollowId := params["userId"]

	authUserId, err := auth.ExtractUserId(r)
	if err != nil {
		responses.Error(w, err)
		return
	}

	if err := controller.service.Unfollow(authUserId, unfollowId); err != nil {
		userErrorResponse(w, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func (controller UserController) ChangePassword(
	w http.ResponseWriter,
	r *http.Request,
) {
	params := mux.Vars(r)
	userId := params["userId"]

	authUserId, err := auth.ExtractUserId(r)
	if err != nil {
		responses.Error(w, err)
		return
	}

	if authUserId != userId {
		responses.Error(w, apierrors.Forbidden("Can't change a different user's password"))
		return
	}

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, apierrors.ValidationError(err.Error()))
		return
	}

	var fields changePasswordRequest
	if err := json.Unmarshal(requestBody, &fields); err != nil {
		responses.Error(w, apierrors.BadRequest(err.Error()))
		return
	}

	if err = controller.service.ChangePassword(
		userId,
		fields.CurrentPwd,
		fields.NewPwd,
	); err != nil {
		userErrorResponse(w, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}
