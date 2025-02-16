package controllers

import (
	"awesomeProject/dao"
	"awesomeProject/dto"
	"awesomeProject/services"
	u "awesomeProject/utils"
	"encoding/json"
	_ "github.com/gorilla/mux"
	"net/http"
)

// Authenticate - эндпоинт для авторизации пользователя.
var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	auth := &dto.AuthRequest{}
	err := json.NewDecoder(r.Body).Decode(auth)

	// Если неправильное тело запроса.
	if err != nil {
		u.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Если пустые поля в теле запроса.
	if auth.Username == "" || auth.Password == "" {
		u.RespondWithError(w, http.StatusBadRequest, "username and password shouldn't be blank strings")
		return
	}

	user, err := services.GetUser(auth.Username)
	if err == nil {
		// Если неправильный пароль.
		if user.Password != auth.Password {
			u.RespondWithError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
	} else {
		newUser := dao.User{Username: auth.Username, Password: auth.Password}
		if err := services.GetDB().Create(&newUser).Error; err != nil {
			u.RespondWithError(w, http.StatusInternalServerError, "error while creating a new user")
			return
		}
	}

	token, err := services.GenerateJWT(auth.Username)
	if err != nil {
		u.RespondWithError(w, http.StatusInternalServerError, "error while generating jwt-token")
		return
	}

	u.RespondJSON(w, http.StatusOK, dto.AuthResponse{Token: token})
}
