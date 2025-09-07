package router

import (
	"net/http"

	"github.com/mohits-git/food-ordering-system/internal/adapters/http/handlers"
)

func NewRouter(
	userHandler *handlers.UserHandler,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// users routes
	mux.HandleFunc("POST /api/users", userHandler.HandleCreateUser)
	mux.HandleFunc("GET /api/users/{id}", userHandler.HandleGetUserById)

	// auth routes

	// restaurants routes

	// menu items routes

	// orders routes

	// invoice routes

	return mux
}
