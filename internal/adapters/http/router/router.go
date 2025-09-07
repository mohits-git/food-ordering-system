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
	menuItemHandler *handlers.MenuItemHandler,
	orderHandler *handlers.OrdersHandler,
	invoiceHandler *handlers.InvoiceHandler,
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
	mux.HandleFunc("GET /api/restaurants/{id}/items", menuItemHandler.HandleGetRestaurantMenuItems)
	mux.HandleFunc("POST /api/restaurants/{id}/items", authMiddleware.Authenticated(menuItemHandler.HandleAddMenuItemToRestaurant))
	mux.HandleFunc("PATCH /api/items/{id}", authMiddleware.Authenticated(menuItemHandler.HandleUpdateAvailability))

	// orders routes
	mux.HandleFunc("POST /api/orders", authMiddleware.Authenticated(orderHandler.HandleCreateOrder))
	mux.HandleFunc("GET /api/orders/{id}", authMiddleware.Authenticated(orderHandler.HandleGetOrderById))
	mux.HandleFunc("POST /api/orders/{id}/items", authMiddleware.Authenticated(orderHandler.HandleAddOrderItem))

	// invoice routes
	mux.HandleFunc("GET /api/invoices/{id}", authMiddleware.Authenticated(invoiceHandler.HandleGetInvoice))
	mux.HandleFunc("POST /api/orders/{id}/invoices", authMiddleware.Authenticated(invoiceHandler.HandleCreateInvoice))
	mux.HandleFunc("POST /api/invoices/{id}/pay", authMiddleware.Authenticated(invoiceHandler.HandleInvoicePayment))

	return mux
}
