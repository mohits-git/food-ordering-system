package dtos

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	UserID int `json:"user_id"`
}

type GetUserResponse struct {
	UserID int `json:"user_id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
