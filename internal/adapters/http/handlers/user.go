package handlers

import (
	"net/http"

	"github.com/mohits-git/food-ordering-system/internal/adapters/http/dtos"
	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/mohits-git/food-ordering-system/internal/ports"
	"github.com/mohits-git/food-ordering-system/internal/utils/apperr"
)

type UserHandler struct {
	userService ports.UserService
}

func NewUserHandler(userService ports.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	createUserReq, err := decodeRequest[dtos.CreateUserRequest](r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	userId, err := h.userService.CreateUser(
		r.Context(),
		domain.NewUser(
			0,
			createUserReq.Name,
			createUserReq.Email,
			createUserReq.Password,
			domain.UserRole(createUserReq.Role)),
	)
	if err != nil {
		if apperr.IsConflictError(err) {
			writeError(w, http.StatusConflict, "user already exists")
		} else {
			writeError(w, http.StatusInternalServerError, "failed to create user")
		}
		return
	}

	resp := dtos.CreateUserResponse{UserID: userId}
	writeResponse(w, http.StatusCreated, "user created successfully", resp)
}

func (h *UserHandler) HandleGetUserById(w http.ResponseWriter, r *http.Request) {
	userId := getIdFromPath(r, "id")
	if userId <= 0 {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.userService.GetUserById(r.Context(), userId)
	if err != nil {
		if apperr.IsNotFoundError(err) {
			writeError(w, http.StatusNotFound, "user not found")
		} else {
			writeError(w, http.StatusInternalServerError, "failed to fetch user")
		}
	}

	resp := dtos.GetUserResponse{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Role:   string(user.Role),
	}
	writeResponse(w, http.StatusOK, "user fetched successfully", resp)
}
