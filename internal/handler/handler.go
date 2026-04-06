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
		// Handle errors (e.g., malformed JSON, wrong type)
		utils.WriteJsonResponse(w, 400, false, err.Error(), nil)
		return
	}

	userIdPayload, userCreateErr := h.service.CreateUser(ctx, newUser)
	if userCreateErr != nil {
		if errors.Is(userCreateErr, context.DeadlineExceeded) {
			utils.WriteJsonResponse(w, 504, false, userCreateErr.Error(), nil)
			return
		}
		if errors.Is(userCreateErr, context.Canceled) {
			utils.WriteJsonResponse(w, 500, false, userCreateErr.Error(), nil)
			return
		}
		utils.WriteJsonResponse(w, 500, false, userCreateErr.Error(), nil)
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
		utils.WriteJsonResponse(w, 400, false, decodeErr.Error(), nil)
		return
	}

	loginPayload, loginErr := h.service.LoginUser(ctx, userLoginInfo)
	if loginErr != nil {
		if errors.Is(loginErr, context.DeadlineExceeded) {
			utils.WriteJsonResponse(w, 504, false, loginErr.Error(), nil)
			return
		}
		if errors.Is(loginErr, context.Canceled) {
			utils.WriteJsonResponse(w, 500, false, loginErr.Error(), nil)
			return
		}
		utils.WriteJsonResponse(w, 401, false, loginErr.Error(), nil)
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
		utils.WriteJsonResponse(w, 500, false, "Update decode failed", nil)
		return
	}

	userUpdateRes, updateErr := h.service.UpdateUser(ctx, userId, updateUserInfo)

	if updateErr != nil {
		if errors.Is(updateErr, context.DeadlineExceeded) {
			utils.WriteJsonResponse(w, 504, false, "Context Timeout", nil)
			return
		}
		if errors.Is(updateErr, context.Canceled) {
			utils.WriteJsonResponse(w, 409, false, "Request cancelled", nil)
			return
		}
		utils.WriteJsonResponse(w, 500, false, updateErr.Error(), nil)
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
		if errors.Is(deleteErr, context.DeadlineExceeded) {
			utils.WriteJsonResponse(w, 504, false, "Context Timeout", nil)
			return
		}
		if errors.Is(deleteErr, context.Canceled) {
			utils.WriteJsonResponse(w, 409, false, "Request cancelled", nil)
			return
		}
		utils.WriteJsonResponse(w, 500, false, deleteErr.Error(), nil)
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
		if errors.Is(userErr, context.DeadlineExceeded) {
			utils.WriteJsonResponse(w, 504, false, "Context Timeout", nil)
			return
		}
		if errors.Is(userErr, context.Canceled) {
			utils.WriteJsonResponse(w, 409, false, "Request cancelled", nil)
			return
		}
		utils.WriteJsonResponse(w, 500, false, userErr.Error(), nil)
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
		if errors.Is(usersErr, context.DeadlineExceeded) {
			utils.WriteJsonResponse(w, 504, false, "Context Timeout", nil)
			return
		}
		if errors.Is(usersErr, context.Canceled) {
			utils.WriteJsonResponse(w, 409, false, "Request cancelled", nil)
			return
		}
		utils.WriteJsonResponse(w, 500, false, usersErr.Error(), nil)
		return
	}

	utils.WriteJsonResponse(w, 200, true, "users fetch success", users)
}
