package handlers

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/mohits-git/food-ordering-system/internal/adapters/http/dtos"
	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/mohits-git/food-ordering-system/internal/utils/apperr"
	mockservice "github.com/mohits-git/food-ordering-system/tests/mock_service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_handlers_RestaurantHandler_NewRestaurantHandler(t *testing.T) {
	service := &mockservice.RestaurantService{}
	handler := NewRestaurantHandler(service)
	require.NotNil(t, handler, "expected NewRestaurantHandler to return a non-nil handler")
}

func Test_handlers_RestaurantHandler_HandleGetRestaurants(t *testing.T) {
	mockservice := &mockservice.RestaurantService{}
	handler := NewRestaurantHandler(mockservice)
	require.NotNil(t, handler, "expected NewRestaurantHandler to return a non-nil handler")

	mockservice.On("GetAllRestaurants", mock.Anything).Return([]domain.Restaurant{
		{
			ID:       1,
			Name:     "Test Restaurant",
			OwnerID:  1,
			ImageURL: "file.com",
		},
	}, nil).Once()

	req := httptest.NewRequest("GET", "/api/restaurants", nil)
	w := httptest.NewRecorder()
	handler.HandleGetRestaurants(w, req)
	res := w.Result()

	require.Equal(t, 200, res.StatusCode, "expected status code 200")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	restaurants, err := decodeResponse[dtos.GetRestaurantsResponse](res)

	require.NoError(t, err, "expected no error while decoding response")
	require.Len(t, restaurants.Restaurants, 1, "expected one restaurant in response")
	require.Equal(t, "Test Restaurant", restaurants.Restaurants[0].Name, "expected restaurant name to be 'Test Restaurant'")
	require.Equal(t, 1, restaurants.Restaurants[0].ID, "expected restaurant ID to be 1")
	require.Equal(t, 1, restaurants.Restaurants[0].OwnerID, "expected restaurant OwnerID to be 1")
}

func Test_handlers_RestaurantHandler_HandleGetRestaurants_Error(t *testing.T) {
	mockservice := &mockservice.RestaurantService{}
	handler := NewRestaurantHandler(mockservice)
	require.NotNil(t, handler, "expected NewRestaurantHandler to return a non-nil handler")

	mockservice.On("GetAllRestaurants", mock.Anything).Return(
		[]domain.Restaurant{}, apperr.NewAppError(apperr.ErrInternal, "failed to fetch restaurants", nil)).Once()

	req := httptest.NewRequest("GET", "/api/restaurants", nil)
	w := httptest.NewRecorder()
	handler.HandleGetRestaurants(w, req)
	res := w.Result()

	require.Equal(t, 500, res.StatusCode, "expected status code 500")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 500, errorResponse.Status, "expected error status to be 500")
	require.Contains(t, errorResponse.Message, "failed to fetch restaurants", "expected error message to contain 'failed to fetch restaurants'")
}

func Test_handlers_RestaurantHandler_HandleCreateRestaurant(t *testing.T) {
	mockservice := &mockservice.RestaurantService{}
	handler := NewRestaurantHandler(mockservice)
	require.NotNil(t, handler, "expected NewRestaurantHandler to return a non-nil handler")

	mockservice.On("CreateRestaurant", mock.Anything, "New Restaurant", "file.com").Return(
		1, nil).Once()

	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, dtos.CreateRestaurantRequest{
		Name:     "New Restaurant",
		ImageURL: "file.com",
	})
	require.NoError(t, err, "expected no error while encoding request body")

	req := httptest.NewRequest("POST", "/api/restaurants", buf)

	w := httptest.NewRecorder()
	handler.HandleCreateRestaurant(w, req)
	res := w.Result()

	require.Equal(t, 201, res.StatusCode, "expected status code 201")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	response, err := decodeResponse[dtos.CreateRestaurantResponse](res)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 1, response.ID, "expected restaurant ID to be 1")
}

func Test_handlers_RestaurantHandler_HandleCreateRestaurant_Unauthorized(t *testing.T) {
	mockservice := &mockservice.RestaurantService{}
	handler := NewRestaurantHandler(mockservice)
	require.NotNil(t, handler, "expected NewRestaurantHandler to return a non-nil handler")

	mockservice.On("CreateRestaurant", mock.Anything, "New Restaurant", "file.com").Return(
		0, apperr.NewAppError(apperr.ErrUnauthorized, "unauthorized", nil)).Once()

	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, dtos.CreateRestaurantRequest{
		Name:     "New Restaurant",
		ImageURL: "file.com",
	})
	require.NoError(t, err, "expected no error while encoding request body")

	req := httptest.NewRequest("POST", "/api/restaurants", buf)
	w := httptest.NewRecorder()
	handler.HandleCreateRestaurant(w, req)
	res := w.Result()

	require.Equal(t, 401, res.StatusCode, "expected status code 401")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 401, errorResponse.Status, "expected error status to be 401")
	require.Contains(t, errorResponse.Message, "unauthorized", "expected error message to contain 'unauthorized'")
}

func Test_handlers_RestaurantHandler_HandleCreateRestaurant_Error(t *testing.T) {
	mockservice := &mockservice.RestaurantService{}
	handler := NewRestaurantHandler(mockservice)
	require.NotNil(t, handler, "expected NewRestaurantHandler to return a non-nil handler")

	mockservice.On("CreateRestaurant", mock.Anything, "New Restaurant", "file.com").Return(
		0, apperr.NewAppError(apperr.ErrInternal, "failed to create restaurant", nil)).Once()

	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, dtos.CreateRestaurantRequest{
		Name:     "New Restaurant",
		ImageURL: "file.com",
	})
	require.NoError(t, err, "expected no error while encoding request body")

	req := httptest.NewRequest("POST", "/api/restaurants", buf)

	w := httptest.NewRecorder()
	handler.HandleCreateRestaurant(w, req)
	res := w.Result()

	require.Equal(t, 500, res.StatusCode, "expected status code 500")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 500, errorResponse.Status, "expected error status to be 500")
	require.Contains(t, errorResponse.Message, "failed to create restaurant", "expected error message to contain 'failed to create restaurant'")
}

func Test_handlers_RestaurantHandler_HandleCreateRestaurant_InvalidRequest(t *testing.T) {
	mockservice := &mockservice.RestaurantService{}
	handler := NewRestaurantHandler(mockservice)
	require.NotNil(t, handler, "expected NewRestaurantHandler to return a non-nil handler")

	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, struct {
		Name int
	}{
		Name: 1,
	})
	require.NoError(t, err, "expected no error while encoding request body")

	req := httptest.NewRequest("POST", "/api/restaurants", buf)

	w := httptest.NewRecorder()
	handler.HandleCreateRestaurant(w, req)
	res := w.Result()

	require.Equal(t, 400, res.StatusCode, "expected status code 400")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 400, errorResponse.Status, "expected error status to be 400")
	require.Contains(t, errorResponse.Message, "invalid request payload", "expected error message to contain 'invalid request payload'")
}

func Test_handlers_RestaurantHandler_HandleCreateRestaurant_EmptyName(t *testing.T) {
	mockservice := &mockservice.RestaurantService{}
	handler := NewRestaurantHandler(mockservice)
	require.NotNil(t, handler, "expected NewRestaurantHandler to return a non-nil handler")

	mockservice.On("CreateRestaurant", mock.Anything, "", "").Return(
		0, apperr.NewAppError(apperr.ErrInvalid, "invalid empty restaurant name", nil)).Once()

	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, dtos.CreateRestaurantRequest{
		Name:     "",
		ImageURL: "",
	})
	require.NoError(t, err, "expected no error while encoding request body")

	req := httptest.NewRequest("POST", "/api/restaurants", buf)

	w := httptest.NewRecorder()
	handler.HandleCreateRestaurant(w, req)
	res := w.Result()

	require.Equal(t, 400, res.StatusCode, "expected status code 400")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 400, errorResponse.Status, "expected error status to be 400")
	require.Contains(t, errorResponse.Message, "invalid empty restaurant name", "expected error message to contain 'restaurant name cannot be empty'")
}
