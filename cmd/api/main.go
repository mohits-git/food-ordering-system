package main

import (
	"context"
	"log"
	"net/http"

	"github.com/mohits-git/food-ordering-system/internal/adapters/bcrypt"
	"github.com/mohits-git/food-ordering-system/internal/adapters/http/handlers"
	"github.com/mohits-git/food-ordering-system/internal/adapters/http/router"
	"github.com/mohits-git/food-ordering-system/internal/adapters/jwttoken"
	"github.com/mohits-git/food-ordering-system/internal/adapters/sqlite"
	"github.com/mohits-git/food-ordering-system/internal/services"
)

func main() {
	ctx := context.Background()

	db, err := sqlite.Connect(ctx, "file:food-ordering-system.db?cache=shared&mode=rwc")
	if err != nil {
		log.Println("Failed to connect to database:", err)
		panic(err)
	}

	err = sqlite.Migrate(db)
	if err != nil {
		log.Println("Failed to migrate the database:", err)
		panic(err)
	}
	log.Println("Database connected and migrated successfully")

	// Initialize utilities
	tokenProvider := jwttoken.NewJWTService("secret-key", "app-name", "app-audience")
	bcryptHasher := bcrypt.NewBcryptPasswordHasher(12)

	// Initialize repositories
	userRepo := sqlite.NewUserRepository(db)
	restaurantRepo := sqlite.NewRestaurantRepository(db)
	menuItemRepo := sqlite.NewMenuItemRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo, bcryptHasher)
	authService := services.NewAuthenticationService(userRepo, tokenProvider, bcryptHasher)
	restaurantService := services.NewRestaurantService(restaurantRepo)
	menuItemService := services.NewMenuItemsService(menuItemRepo, restaurantRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)
	restaurantHandler := handlers.NewRestaurantHandler(restaurantService)
	menuItemHandler := handlers.NewMenuItemHandler(menuItemService)

	// middlewares
	authMiddleware := handlers.NewAuthMiddleware(tokenProvider)

	// Initialize router
	mux := router.NewRouter(
		authMiddleware,
		userHandler,
		authHandler,
		restaurantHandler,
		menuItemHandler,
	)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
