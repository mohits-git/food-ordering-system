package router

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/mohits-git/food-ordering-system/internal/adapters/http/handlers"
	"github.com/stretchr/testify/require"
)

func Test_router_NewRouter(t *testing.T) {
	router := NewRouter(
		handlers.NewAuthMiddleware(nil),
		handlers.NewUserHandler(nil),
		handlers.NewAuthHandler(nil),
		handlers.NewRestaurantHandler(nil),
		handlers.NewMenuItemHandler(nil),
		handlers.NewOrdersHandler(nil),
		handlers.NewInvoiceHandler(nil),
	)
	require.NotNil(t, router, "expected NewRouter to return a non-nil router")

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)
	res := w.Result()
	require.Equal(t, 200, res.StatusCode, "expected status code 200")
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err, "expected no error while reading response body")
	require.Equal(t, "OK", string(body), "expected body to be 'OK'")
}
