package login

import (
	"api/src/apierrors"
	"api/src/responses"
	"api/src/users"
	"encoding/json"
	"io"
	"net/http"
)

type LoginController struct{
	service *LoginService
}

func NewLoginController(service *LoginService) *LoginController {
	return &LoginController{service: service}
}

func (controller LoginController) Login(w http.ResponseWriter, r *http.Request) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, apierrors.UnprocessableEntity(err.Error()))
		return
	}

	var user users.User
	if err := json.Unmarshal(requestBody, &user); err != nil {
		responses.Error(w, apierrors.BadRequest(err.Error()))
		return
	}

	token, err := controller.service.Login(&user)
	if err != nil {
		responses.Error(w, err)
		return
	}

	responses.JSON(w, http.StatusOK, map[string]string{
		"token": token,
	})
}
