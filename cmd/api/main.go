package main

import (
	"context"
	"log"
	"net/http"

	"github.com/mohits-git/food-ordering-system/internal/adapters/http/handlers"
	"github.com/mohits-git/food-ordering-system/internal/adapters/http/router"
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

	// Initialize repositories
	userRepo := sqlite.NewUserRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)

	mux := router.NewRouter(
		userHandler,
	)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
