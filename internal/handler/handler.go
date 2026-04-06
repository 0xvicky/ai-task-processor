package handler

import (
	"ai-task-processor/internal/model"
	"ai-task-processor/internal/service"
	"ai-task-processor/internal/utils"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{
		service: s,
	}
}

// default route for "/"
func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ai Task Processor running on port 6969")
}

func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	var newUser model.User
	// println(r.Body)
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&newUser)
	if err != nil {
		utils.ErrorHandler(w, err)
		return
	}

	userIdPayload, userCreateErr := h.service.CreateUser(ctx, newUser)
	if userCreateErr != nil {
		utils.ErrorHandler(w, userCreateErr)
		return
	}

	utils.WriteJsonResponse(w, 201, true, "User Created ✅", userIdPayload)
}

func (h *UserHandler) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var userLoginInfo model.UserLogin
	decoder := json.NewDecoder(r.Body)
	decodeErr := decoder.Decode(&userLoginInfo)
	if decodeErr != nil {
		utils.ErrorHandler(w, decodeErr)
		return
	}

	loginPayload, loginErr := h.service.LoginUser(ctx, userLoginInfo)
	if loginErr != nil {
		utils.ErrorHandler(w, loginErr)
		return
	}

	utils.WriteJsonResponse(w, 200, true, "Login Success", loginPayload)

}

// update user info
func (h *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	userId := utils.ExtractUserId(r)
	var updateUserInfo model.UserUpdate

	decoder := json.NewDecoder(r.Body)
	decodeErr := decoder.Decode(&updateUserInfo)
	if decodeErr != nil {
		utils.ErrorHandler(w, decodeErr)
		return
	}

	userUpdateRes, updateErr := h.service.UpdateUser(ctx, userId, updateUserInfo)

	if updateErr != nil {
		utils.ErrorHandler(w, updateErr)
		return
	}

	updatePayload := struct {
		UserId    int       `json:"user_id"`
		Name      string    `json:"user_name"`
		Email     string    `json:"user_email"`
		CreatedAt time.Time `json:"created_at"`
	}{
		UserId:    userUpdateRes.UserId,
		Name:      userUpdateRes.Name,
		Email:     userUpdateRes.Email,
		CreatedAt: userUpdateRes.CreatedAt,
	}

	utils.WriteJsonResponse(w, 200, true, "Update Success", updatePayload)

}

// delete user
func (h *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	userId := utils.ExtractUserId(r)
	deletedUserRes, deleteErr := h.service.DeleteUser(ctx, userId)
	if errors.Is(deleteErr, sql.ErrNoRows) {
		utils.WriteJsonResponse(w, 200, true, "User Delete Success", nil)
		return
	}
	if deleteErr != nil {
		utils.ErrorHandler(w, deleteErr)
		return
	}

	utils.WriteJsonResponse(w, 200, true, "User Deleted", deletedUserRes)
}

func (h *UserHandler) MeHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel() //clean the messup, resources, called when function returns

	userId := utils.ExtractUserId(r)

	userRes, userErr := h.service.FetchUserById(ctx, userId)

	if userErr != nil {
		utils.ErrorHandler(w, userErr)
		return
	}

	utils.WriteJsonResponse(w, 200, true, "User Fetch Success", userRes)

}

func (h *UserHandler) FetchAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	users, usersErr := h.service.FetchAllUsers(ctx)

	if usersErr != nil {
		utils.ErrorHandler(w, usersErr)
		return
	}

	utils.WriteJsonResponse(w, 200, true, "users fetch success", users)
}
