package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/mohits-git/food-ordering-system/internal/adapters/bcrypt"
	"github.com/mohits-git/food-ordering-system/internal/adapters/http/handlers"
	"github.com/mohits-git/food-ordering-system/internal/adapters/http/router"
	imageupload "github.com/mohits-git/food-ordering-system/internal/adapters/image-upload"
	"github.com/mohits-git/food-ordering-system/internal/adapters/jwttoken"
	"github.com/mohits-git/food-ordering-system/internal/adapters/sqlite"
	"github.com/mohits-git/food-ordering-system/internal/services"
)

func main() {
	ctx := context.Background()

	// load configurations
	config := LoadConfig()

	// database setup
	db := SetupDB(ctx, config.SQLITE_DSN)

	// Initialize utilities
	tokenProvider := jwttoken.NewJWTService(
		config.JWT_SECRET,
		config.JWT_ISSUER,
		config.JWT_AUDIENCE,
	)
	bcryptHasher := bcrypt.NewBcryptPasswordHasher(12)

	// Initialize repositories
	userRepo := sqlite.NewUserRepository(db)
	restaurantRepo := sqlite.NewRestaurantRepository(db)
	menuItemRepo := sqlite.NewMenuItemRepository(db)
	orderRepo := sqlite.NewOrderRepository(db)
	invoiceRepo := sqlite.NewInvoiceRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo, bcryptHasher)
	authService := services.NewAuthenticationService(userRepo, tokenProvider, bcryptHasher)
	restaurantService := services.NewRestaurantService(restaurantRepo)
	menuItemService := services.NewMenuItemsService(menuItemRepo, restaurantRepo)
	orderService := services.NewOrderService(orderRepo, menuItemRepo)
	invoiceService := services.NewInvoiceService(invoiceRepo, orderRepo, menuItemRepo)
	imageUploadService := imageupload.NewFSImageUpload("http://localhost:8080", config.UPLOAD_DIRECTORY)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)
	restaurantHandler := handlers.NewRestaurantHandler(restaurantService)
	menuItemHandler := handlers.NewMenuItemHandler(menuItemService)
	orderHandler := handlers.NewOrdersHandler(orderService)
	invoiceHandler := handlers.NewInvoiceHandler(invoiceService)
	imageUploadHandler := handlers.NewImageUploadHandler(imageUploadService)

	// middlewares
	authMiddleware := handlers.NewAuthMiddleware(tokenProvider)

	// Initialize router
	mux := router.NewRouter(
		authMiddleware,
		userHandler,
		authHandler,
		restaurantHandler,
		menuItemHandler,
		orderHandler,
		invoiceHandler,
		imageUploadHandler,
	)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func SetupDB(ctx context.Context, dsn string) *sql.DB {
	db, err := sqlite.Connect(ctx, dsn)
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
	return db
}
