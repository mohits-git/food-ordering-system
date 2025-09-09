package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/mohits-git/food-ordering-system/internal/adapters/http/dtos"
	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/mohits-git/food-ordering-system/internal/utils/apperr"
	mockservice "github.com/mohits-git/food-ordering-system/tests/mock_service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_handlers_MenuItemHandlers_NewMenuItemHandler(t *testing.T) {
	mockservice := &mockservice.MenuItemService{}
	handler := NewMenuItemHandler(mockservice)
	require.NotNil(t, handler, "expected NewMenuItemHandler to return a non-nil handler")
}

func Test_handlers_MenuItemHandler_HandleAddMenuItemToRestaurant(t *testing.T) {
	mockservice := &mockservice.MenuItemService{}
	handler := NewMenuItemHandler(mockservice)
	require.NotNil(t, handler, "expected NewMenuItemHandler to return a non-nil handler")

	addRequest := dtos.AddMenuItemRequest{
		Name:      "Test Menu Item",
		Price:     9.99,
		Available: true,
	}
	requestBody, err := json.Marshal(addRequest)
	require.NoError(t, err, "expected no error while marshalling request body")

	mockservice.On("CreateMenuItemForRestaurant", mock.Anything, mock.MatchedBy(func(menuItem domain.MenuItem) bool {
		return menuItem.Name == addRequest.Name &&
			menuItem.Price == addRequest.Price &&
			menuItem.Available == addRequest.Available &&
			menuItem.RestaurantID == 1
	})).Return(1, nil).Once()

	req := httptest.NewRequest("POST", "/api/restaurants/1/items", bytes.NewReader(requestBody))
	req.SetPathValue("id", "1")

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.HandleAddMenuItemToRestaurant(w, req)
	res := w.Result()

	require.Equal(t, 201, res.StatusCode, "expected status code 201")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	addResponse, err := decodeResponse[dtos.AddMenuItemResponse](res)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 1, addResponse.ID, "expected menu item ID to be 1")
}

func Test_handlers_MenuItemHandler_HandleAddMenuItemToRestaurant_InvalidRequestBody(t *testing.T) {
	mockservice := &mockservice.MenuItemService{}
	handler := NewMenuItemHandler(mockservice)
	require.NotNil(t, handler, "expected NewMenuItemHandler to return a non-nil handler")

	invalidRequestBody := []byte(`{"name": 1, "price": -5, "available": true}`)

	req := httptest.NewRequest("POST", "/api/restaurants/1/items", bytes.NewReader(invalidRequestBody))
	req.SetPathValue("id", "1")

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.HandleAddMenuItemToRestaurant(w, req)
	res := w.Result()

	require.Equal(t, 400, res.StatusCode, "expected status code 400")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
}

func Test_handlers_MenuItemHandler_HandleAddMenuItemToRestaurant_InvalidId(t *testing.T) {
	mockservice := &mockservice.MenuItemService{}
	handler := NewMenuItemHandler(mockservice)
	require.NotNil(t, handler, "expected NewMenuItemHandler to return a non-nil handler")

	req := httptest.NewRequest("POST", "/api/restaurants/abc/items", nil)
	req.SetPathValue("id", "abc")

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.HandleAddMenuItemToRestaurant(w, req)
	res := w.Result()

	require.Equal(t, 400, res.StatusCode, "expected status code 400")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
}

func Test_handlers_MenuItemHandler_HandleAddMenuItemToRestaurant_InvalidRequest(t *testing.T) {
	mockservice := &mockservice.MenuItemService{}
	handler := NewMenuItemHandler(mockservice)
	require.NotNil(t, handler, "expected NewMenuItemHandler to return a non-nil handler")

	mockservice.On("CreateMenuItemForRestaurant", mock.Anything, mock.Anything).Return(
		0, apperr.NewAppError(apperr.ErrInvalid, "invalid menu item data", nil)).Once()

	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, dtos.AddMenuItemRequest{
		Name:  "",
		Price: 10,
	})
	require.NoError(t, err, "expected no error while encoding request body")

	req := httptest.NewRequest("POST", "/api/restaurants/1/items", buf)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleAddMenuItemToRestaurant(w, req)
	res := w.Result()

	require.Equal(t, 400, res.StatusCode, "expected status code 400")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 400, errorResponse.Status, "expected error status to be 400")
	require.Contains(t, errorResponse.Message, "invalid menu item", "expected error message to contain 'invalid menu item'")
}

func Test_handlers_MenuItemHandler_HandleAddMenuItemToRestaurant_ServiceError(t *testing.T) {
	mockservice := &mockservice.MenuItemService{}
	handler := NewMenuItemHandler(mockservice)
	require.NotNil(t, handler, "expected NewMenuItemHandler to return a non-nil handler")

	addRequest := dtos.AddMenuItemRequest{
		Name:      "Test Menu Item",
		Price:     9.99,
		Available: true,
	}
	requestBody, err := json.Marshal(addRequest)
	require.NoError(t, err, "expected no error while marshalling request body")

	mockservice.On("CreateMenuItemForRestaurant", mock.Anything, mock.MatchedBy(func(menuItem domain.MenuItem) bool {
		return menuItem.Name == addRequest.Name &&
			menuItem.Price == addRequest.Price &&
			menuItem.Available == addRequest.Available &&
			menuItem.RestaurantID == 1
	})).Return(
		0, apperr.NewAppError(apperr.ErrInternal, "failed to create menu item", nil)).Once()

	req := httptest.NewRequest("POST", "/api/restaurants/1/items", bytes.NewReader(requestBody))
	req.SetPathValue("id", "1")

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.HandleAddMenuItemToRestaurant(w, req)
	res := w.Result()

	require.Equal(t, 500, res.StatusCode, "expected status code 500")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 500, errorResponse.Status, "expected error status to be 500")
	require.Contains(t, errorResponse.Message, "internal server error", "expected error message to contain 'internal server error'")
}

func Test_handlers_MenuItemHandler_HandleAddMenuItemToRestaurant_Unauthorized(t *testing.T) {
	mockservice := &mockservice.MenuItemService{}
	handler := NewMenuItemHandler(mockservice)
	require.NotNil(t, handler, "expected NewMenuItemHandler to return a non-nil handler")

	addRequest := dtos.AddMenuItemRequest{
		Name:      "Test Menu Item",
		Price:     9.99,
		Available: true,
	}
	requestBody, err := json.Marshal(addRequest)
	require.NoError(t, err, "expected no error while marshalling request body")

	mockservice.On("CreateMenuItemForRestaurant", mock.Anything, mock.MatchedBy(func(menuItem domain.MenuItem) bool {
		return menuItem.Name == addRequest.Name &&
			menuItem.Price == addRequest.Price &&
			menuItem.Available == addRequest.Available &&
			menuItem.RestaurantID == 1
	})).Return(
		0, apperr.NewAppError(apperr.ErrUnauthorized, "unauthenticated user", nil)).Once()

	req := httptest.NewRequest("POST", "/api/restaurants/1/items", bytes.NewReader(requestBody))
	req.SetPathValue("id", "1")

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.HandleAddMenuItemToRestaurant(w, req)
	res := w.Result()

	require.Equal(t, 401, res.StatusCode, "expected status code 401")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 401, errorResponse.Status, "expected error status to be 401")
	require.Contains(t, errorResponse.Message, "unauthenticated user", "expected error message to contain 'unauthenticated user'")
}

func Test_handlers_MenuItemHandler_HandleAddMenuItemToRestaurant_Forbidden(t *testing.T) {
	mockservice := &mockservice.MenuItemService{}
	handler := NewMenuItemHandler(mockservice)
	require.NotNil(t, handler, "expected NewMenuItemHandler to return a non-nil handler")

	addRequest := dtos.AddMenuItemRequest{
		Name:      "Test Menu Item",
		Price:     9.99,
		Available: true,
	}
	requestBody, err := json.Marshal(addRequest)
	require.NoError(t, err, "expected no error while marshalling request body")

	mockservice.On("CreateMenuItemForRestaurant", mock.Anything, mock.MatchedBy(func(menuItem domain.MenuItem) bool {
		return menuItem.Name == addRequest.Name &&
			menuItem.Price == addRequest.Price &&
			menuItem.Available == addRequest.Available &&
			menuItem.RestaurantID == 1
	})).Return(
		0, apperr.NewAppError(apperr.ErrForbidden, "only restaurant owners can add menu items", nil)).Once()

	req := httptest.NewRequest("POST", "/api/restaurants/1/items", bytes.NewReader(requestBody))
	req.SetPathValue("id", "1")

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.HandleAddMenuItemToRestaurant(w, req)
	res := w.Result()

	require.Equal(t, 403, res.StatusCode, "expected status code 403")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 403, errorResponse.Status, "expected error status to be 403")
	require.Contains(t, errorResponse.Message, "only restaurant owners can add menu items", "expected error message to contain 'only restaurant owners can add menu items'")
}

func Test_handlers_MenuItemHandler_HandleGetRestaurantMenuItems(t *testing.T) {
	mockservice := &mockservice.MenuItemService{}
	handler := NewMenuItemHandler(mockservice)
	require.NotNil(t, handler, "expected NewMenuItemHandler to return a non-nil handler")

	menuItems := []domain.MenuItem{
		domain.NewMenuItem(1, "Item 1", 10.0, true, 1),
		domain.NewMenuItem(2, "Item 2", 15.0, false, 1),
	}

	mockservice.On("GetAllMenuItemsByRestaurantId", mock.Anything, 1).Return(menuItems, nil).Once()

	req := httptest.NewRequest("GET", "/api/restaurants/1/items", nil)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleGetRestaurantMenuItems(w, req)
	res := w.Result()

	require.Equal(t, 200, res.StatusCode, "expected status code 200")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	getResponse, err := decodeResponse[dtos.GetMenuItemsResponse](res)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 1, getResponse.RestaurantID, "expected restaurant ID to be 1")
	require.Len(t, getResponse.Items, 2, "expected 2 menu items in response")

	require.Equal(t, 1, getResponse.Items[0].ID)
	require.Equal(t, "Item 1", getResponse.Items[0].Name)
	require.Equal(t, 10.0, getResponse.Items[0].Price)
	require.True(t, getResponse.Items[0].Available)

	require.Equal(t, 2, getResponse.Items[1].ID)
	require.Equal(t, "Item 2", getResponse.Items[1].Name)
	require.Equal(t, 15.0, getResponse.Items[1].Price)
	require.False(t, getResponse.Items[1].Available)
}

func Test_handlers_MenuItemHandler_HandleGetRestaurantMenuItems_InvalidId(t *testing.T) {
	mockservice := &mockservice.MenuItemService{}
	handler := NewMenuItemHandler(mockservice)
	require.NotNil(t, handler, "expected NewMenuItemHandler to return a non-nil handler")

	req := httptest.NewRequest("GET", "/api/restaurants/abc/items", nil)
	req.SetPathValue("id", "abc")

	w := httptest.NewRecorder()
	handler.HandleGetRestaurantMenuItems(w, req)
	res := w.Result()

	require.Equal(t, 400, res.StatusCode, "expected status code 400")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
}

func Test_handlers_MenuItemHandler_HandleGetRestaurantMenuItems_ServiceError(t *testing.T) {
	mockservice := &mockservice.MenuItemService{}
	handler := NewMenuItemHandler(mockservice)
	require.NotNil(t, handler, "expected NewMenuItemHandler to return a non-nil handler")

	mockservice.On("GetAllMenuItemsByRestaurantId", mock.Anything, 1).Return(
		[]domain.MenuItem{}, apperr.NewAppError(apperr.ErrInternal, "failed to fetch menu items", nil)).Once()

	req := httptest.NewRequest("GET", "/api/restaurants/1/items", nil)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleGetRestaurantMenuItems(w, req)
	res := w.Result()

	require.Equal(t, 500, res.StatusCode, "expected status code 500")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 500, errorResponse.Status, "expected error status to be 500")
	require.Contains(t, errorResponse.Message, "internal server error", "expected error message to contain 'internal server error'")
}

func Test_handlers_MenuItemHandler_HandleGetRestaurantMenuItems_NoItems(t *testing.T) {
	mockservice := &mockservice.MenuItemService{}
	handler := NewMenuItemHandler(mockservice)
	require.NotNil(t, handler, "expected NewMenuItemHandler to return a non-nil handler")

	mockservice.On("GetAllMenuItemsByRestaurantId", mock.Anything, 1).Return([]domain.MenuItem{}, nil).Once()

	req := httptest.NewRequest("GET", "/api/restaurants/1/items", nil)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleGetRestaurantMenuItems(w, req)
	res := w.Result()

	require.Equal(t, 200, res.StatusCode, "expected status code 200")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	getResponse, err := decodeResponse[dtos.GetMenuItemsResponse](res)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 1, getResponse.RestaurantID, "expected restaurant ID to be 1")
	require.Len(t, getResponse.Items, 0, "expected 0 menu items in response")
}

func Test_handlers_MenuItemHandler_HandleUpdateAvailability(t *testing.T) {
	mockservice := &mockservice.MenuItemService{}
	handler := NewMenuItemHandler(mockservice)
	require.NotNil(t, handler, "expected NewMenuItemHandler to return a non-nil handler")

	updateRequest := dtos.UpdateMenuItemAvailabilityRequest{
		Available: false,
	}
	requestBody, err := json.Marshal(updateRequest)
	require.NoError(t, err, "expected no error while marshalling request body")

	mockservice.On("UpdateAvailability", mock.Anything, 1, false).Return(nil).Once()

	req := httptest.NewRequest("PATCH", "/api/restaurants/1/items/1/availability", bytes.NewReader(requestBody))
	req.SetPathValue("id", "1")
	req.SetPathValue("itemId", "1")

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.HandleUpdateAvailability(w, req)
	res := w.Result()

	require.Equal(t, 200, res.StatusCode, "expected status code 200")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	updateResponse, err := decodeJson[dtos.BaseResponse](res.Body)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 200, updateResponse.Status, "expected status to be 200")
	require.Equal(t, "menu item availability updated successfully", updateResponse.Message, "expected success message")
}

func Test_handlers_MenuItemHandler_HandleUpdateAvailability_InvalidId(t *testing.T) {
	mockservice := &mockservice.MenuItemService{}
	handler := NewMenuItemHandler(mockservice)
	require.NotNil(t, handler, "expected NewMenuItemHandler to return a non-nil handler")

	req := httptest.NewRequest("PATCH", "/api/restaurants/1/items/abc/availability", nil)
	req.SetPathValue("id", "1")
	req.SetPathValue("itemId", "abc")

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.HandleUpdateAvailability(w, req)
	res := w.Result()

	require.Equal(t, 400, res.StatusCode, "expected status code 400")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
}

func Test_handlers_MenuItemHandler_HandleUpdateAvailability_InvalidRequestBody(t *testing.T) {
	mockservice := &mockservice.MenuItemService{}
	handler := NewMenuItemHandler(mockservice)
	require.NotNil(t, handler, "expected NewMenuItemHandler to return a non-nil handler")

	invalidRequestBody := []byte(`{"available": "yes"}`)

	req := httptest.NewRequest("PATCH", "/api/restaurants/1/items/1/availability", bytes.NewReader(invalidRequestBody))
	req.SetPathValue("id", "1")
	req.SetPathValue("itemId", "1")

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.HandleUpdateAvailability(w, req)
	res := w.Result()

	require.Equal(t, 400, res.StatusCode, "expected status code 400")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
}

func Test_handlers_MenuItemHandler_HandleUpdateAvailability_ServiceError(t *testing.T) {
	mockservice := &mockservice.MenuItemService{}
	handler := NewMenuItemHandler(mockservice)
	require.NotNil(t, handler, "expected NewMenuItemHandler to return a non-nil handler")

	updateRequest := dtos.UpdateMenuItemAvailabilityRequest{
		Available: false,
	}
	requestBody, err := json.Marshal(updateRequest)
	require.NoError(t, err, "expected no error while marshalling request body")

	mockservice.On("UpdateAvailability", mock.Anything, 1, false).Return(
		apperr.NewAppError(apperr.ErrInternal, "failed to update availability", nil)).Once()

	req := httptest.NewRequest("PATCH", "/api/restaurants/1/items/1/availability", bytes.NewReader(requestBody))
	req.SetPathValue("id", "1")
	req.SetPathValue("itemId", "1")

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.HandleUpdateAvailability(w, req)
	res := w.Result()

	require.Equal(t, 500, res.StatusCode, "expected status code 500")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 500, errorResponse.Status, "expected error status to be 500")
	require.Contains(t, errorResponse.Message, "internal server error", "expected error message to contain 'internal server error'")
}

func Test_handlers_MenuItemHandler_HandleUpdateAvailability_Unauthorized(t *testing.T) {
	mockservice := &mockservice.MenuItemService{}
	handler := NewMenuItemHandler(mockservice)
	require.NotNil(t, handler, "expected NewMenuItemHandler to return a non-nil handler")

	updateRequest := dtos.UpdateMenuItemAvailabilityRequest{
		Available: false,
	}
	requestBody, err := json.Marshal(updateRequest)
	require.NoError(t, err, "expected no error while marshalling request body")

	mockservice.On("UpdateAvailability", mock.Anything, 1, false).Return(
		apperr.NewAppError(apperr.ErrUnauthorized, "unauthenticated user", nil)).Once()

	req := httptest.NewRequest("PATCH", "/api/restaurants/1/items/1/availability", bytes.NewReader(requestBody))
	req.SetPathValue("id", "1")
	req.SetPathValue("itemId", "1")

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.HandleUpdateAvailability(w, req)
	res := w.Result()

	require.Equal(t, 401, res.StatusCode, "expected status code 401")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 401, errorResponse.Status, "expected error status to be 401")
	require.Contains(t, errorResponse.Message, "unauthenticated user", "expected error message to contain 'unauthenticated user'")
}

func Test_handlers_MenuItemHandler_HandleUpdateAvailability_Forbidden(t *testing.T) {
	mockservice := &mockservice.MenuItemService{}
	handler := NewMenuItemHandler(mockservice)
	require.NotNil(t, handler, "expected NewMenuItemHandler to return a non-nil handler")

	updateRequest := dtos.UpdateMenuItemAvailabilityRequest{
		Available: false,
	}
	requestBody, err := json.Marshal(updateRequest)
	require.NoError(t, err, "expected no error while marshalling request body")

	mockservice.On("UpdateAvailability", mock.Anything, 1, false).Return(
		apperr.NewAppError(apperr.ErrForbidden, "only restaurant owners can update menu items", nil)).Once()

	req := httptest.NewRequest("PATCH", "/api/restaurants/1/items/1/availability", bytes.NewReader(requestBody))
	req.SetPathValue("id", "1")
	req.SetPathValue("itemId", "1")

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.HandleUpdateAvailability(w, req)
	res := w.Result()

	require.Equal(t, 403, res.StatusCode, "expected status code 403")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 403, errorResponse.Status, "expected error status to be 403")
	require.Contains(t, errorResponse.Message, "only restaurant owners can update menu items", "expected error message to contain 'only restaurant owners can update menu items'")
}

func Test_handlers_MenuItemHandler_HandleUpdateAvailability_NotFound(t *testing.T) {
	mockservice := &mockservice.MenuItemService{}
	handler := NewMenuItemHandler(mockservice)
	require.NotNil(t, handler, "expected NewMenuItemHandler to return a non-nil handler")

	updateRequest := dtos.UpdateMenuItemAvailabilityRequest{
		Available: false,
	}
	requestBody, err := json.Marshal(updateRequest)
	require.NoError(t, err, "expected no error while marshalling request body")

	mockservice.On("UpdateAvailability", mock.Anything, 1, false).Return(
		apperr.NewAppError(apperr.ErrNotFound, "menu item not found", nil)).Once()

	req := httptest.NewRequest("PATCH", "/api/restaurants/1/items/1/availability", bytes.NewReader(requestBody))
	req.SetPathValue("id", "1")
	req.SetPathValue("itemId", "1")

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.HandleUpdateAvailability(w, req)
	res := w.Result()

	require.Equal(t, 404, res.StatusCode, "expected status code 404")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 404, errorResponse.Status, "expected error status to be 404")
	require.Contains(t, errorResponse.Message, "menu item not found", "expected error message to contain 'menu item not found'")
}
