package router

import (
	"net/http"

	"github.com/mohits-git/food-ordering-system/internal/adapters/http/handlers"
)

func NewRouter(
	authMiddleware *handlers.AuthMiddleware,
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
	restaurantHandler *handlers.RestaurantHandler,
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
	mux.HandleFunc("POST /api/auth/login", authHandler.HandleLogin)
	mux.HandleFunc("POST /api/auth/logout", authMiddleware.WithToken(authHandler.HandleLogout))

	// restaurants routes
	mux.HandleFunc("GET /api/restaurants", restaurantHandler.HandleGetRestaurants)
	mux.HandleFunc("POST /api/restaurants", authMiddleware.Authenticated(restaurantHandler.HandleCreateRestaurant))

	// menu items routes

	// orders routes

	// invoice routes

	return mux
}
