package handlers

import (
	"fmt"
	"net/http"
	"strconv"

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
	createUserReq, err := decodeJson[dtos.CreateUserRequest](r)
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

	resp := dtos.NewResponse(
		http.StatusCreated,
		"user created successfully",
		dtos.CreateUserResponse{
			UserID: userId,
		},
	)
	encodeJson(w, http.StatusCreated, resp)
}

func (h *UserHandler) HandleGetUserById(w http.ResponseWriter, r *http.Request) {
	userIdParam := r.PathValue("id")
	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid user id: %s", userIdParam))
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

	resp := dtos.NewResponse(
		http.StatusOK,
		"user fetched successfully",
		dtos.GetUserResponse{
			UserID: user.ID,
			Name:   user.Name,
			Email:  user.Email,
			Role:   string(user.Role),
		},
	)
	encodeJson(w, http.StatusOK, resp)
}
