package main

import (
	"context"
	"log"

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
	_ = userService
}
