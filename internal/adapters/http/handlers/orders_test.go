package handlers

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/mohits-git/food-ordering-system/internal/adapters/http/dtos"
	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/mohits-git/food-ordering-system/internal/utils/apperr"
	"github.com/mohits-git/food-ordering-system/internal/utils/authctx"
	mockservice "github.com/mohits-git/food-ordering-system/tests/mock_service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_handlers_OrdersHandler_NewOrdersHandler(t *testing.T) {
	mockOrderService := &mockservice.OrderService{}
	handler := NewOrdersHandler(mockOrderService)
	require.NotNil(t, handler, "expected NewOrdersHandler to return a non-nil handler")
}

func Test_handlers_OrdersHandler_HandleCreateOrder(t *testing.T) {
	mockOrderService := &mockservice.OrderService{}
	handler := NewOrdersHandler(mockOrderService)
	require.NotNil(t, handler, "expected NewOrdersHandler to return a non-nil handler")

	mockOrderService.On("CreateOrder", mock.Anything, domain.Order{
		CustomerID:   1,
		RestaurantID: 1,
		OrderItems: []domain.OrderItem{
			{MenuItemID: 1, Quantity: 2},
			{MenuItemID: 2, Quantity: 1},
		},
	}).Return(1, nil).Once()

	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, dtos.CreateOrderRequest{
		RestaurantID: 1,
		OrderItems: []dtos.OrderItemsDTO{
			{MenuItemID: 1, Quantity: 2},
			{MenuItemID: 2, Quantity: 1},
		},
	})
	require.NoError(t, err, "expected no error while encoding request")
	req := httptest.NewRequest("POST", "/api/orders", buf)

	userClaim := &authctx.UserClaims{
		UserID: 1,
		Role:   "customer",
	}
	ctx := authctx.WithUserClaims(req.Context(), userClaim)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.HandleCreateOrder(w, req)
	res := w.Result()

	require.Equal(t, 201, res.StatusCode, "expected status code 201")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	createOrderResp, err := decodeResponse[dtos.CreateOrderResponse](res)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 1, createOrderResp.ID, "expected order ID to be 1")
}

func Test_handlers_OrdersHandler_HandleCreateOrder_Unauthorized(t *testing.T) {
	mockOrderService := &mockservice.OrderService{}
	handler := NewOrdersHandler(mockOrderService)
	require.NotNil(t, handler, "expected NewOrdersHandler to return a non-nil handler")

	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, dtos.CreateOrderRequest{
		RestaurantID: 1,
		OrderItems: []dtos.OrderItemsDTO{
			{MenuItemID: 1, Quantity: 2},
			{MenuItemID: 2, Quantity: 1},
		},
	})
	require.NoError(t, err, "expected no error while encoding request")
	req := httptest.NewRequest("POST", "/api/orders", buf)

	w := httptest.NewRecorder()
	handler.HandleCreateOrder(w, req)
	res := w.Result()

	require.Equal(t, 401, res.StatusCode, "expected status code 401")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	_, err = decodeResponse[dtos.CreateOrderResponse](res)

	require.NoError(t, err, "expected no error while decoding response")
}

func Test_handlers_OrdersHandler_HandleCreateOrder_Forbidden(t *testing.T) {
	mockOrderService := &mockservice.OrderService{}
	handler := NewOrdersHandler(mockOrderService)
	require.NotNil(t, handler, "expected NewOrdersHandler to return a non-nil handler")

	mockOrderService.On("CreateOrder", mock.Anything, domain.Order{
		CustomerID:   2,
		RestaurantID: 1,
		OrderItems: []domain.OrderItem{
			{MenuItemID: 1, Quantity: 2},
			{MenuItemID: 2, Quantity: 1},
		},
	}).Return(0, apperr.NewAppError(apperr.ErrForbidden, "forbidden", nil)).Once()

	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, dtos.CreateOrderRequest{
		RestaurantID: 1,
		OrderItems: []dtos.OrderItemsDTO{
			{MenuItemID: 1, Quantity: 2},
			{MenuItemID: 2, Quantity: 1},
		},
	})
	require.NoError(t, err, "expected no error while encoding request")
	req := httptest.NewRequest("POST", "/api/orders", buf)

	userClaim := &authctx.UserClaims{
		UserID: 2,
		Role:   "owner",
	}
	ctx := authctx.WithUserClaims(req.Context(), userClaim)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.HandleCreateOrder(w, req)
	res := w.Result()

	require.Equal(t, 403, res.StatusCode, "expected status code 403")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 403, errorResponse.Status, "expected error status to be 403")
	require.Contains(t, errorResponse.Message, "forbidden", "expected error message to contain 'forbidden'")
}

func Test_handlers_OrdersHandler_HandleCreateOrder_InvalidRequest(t *testing.T) {
	mockOrderService := &mockservice.OrderService{}
	handler := NewOrdersHandler(mockOrderService)
	require.NotNil(t, handler, "expected NewOrdersHandler to return a non-nil handler")

	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, struct {
		RestaurantID string               `json:"restaurant_id"`
		OrderItems   []dtos.OrderItemsDTO `json:"order_items"`
	}{
		RestaurantID: "one",
		OrderItems: []dtos.OrderItemsDTO{
			{MenuItemID: 1, Quantity: 2},
			{MenuItemID: 2, Quantity: 1},
		},
	})
	require.NoError(t, err, "expected no error while encoding request")
	req := httptest.NewRequest("POST", "/api/orders", buf)

	userClaim := &authctx.UserClaims{
		UserID: 1,
		Role:   "customer",
	}
	ctx := authctx.WithUserClaims(req.Context(), userClaim)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.HandleCreateOrder(w, req)
	res := w.Result()

	require.Equal(t, 400, res.StatusCode, "expected status code 400")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 400, errorResponse.Status, "expected error status to be 400")
	require.Contains(t, errorResponse.Message, "invalid", "expected error message to contain 'invalid'")
}

func Test_handlers_OrdersHandler_HandleCreateOrder_InternalServerError(t *testing.T) {
	mockOrderService := &mockservice.OrderService{}
	handler := NewOrdersHandler(mockOrderService)
	require.NotNil(t, handler, "expected NewOrdersHandler to return a non-nil handler")

	mockOrderService.On("CreateOrder", mock.Anything, domain.Order{
		CustomerID:   1,
		RestaurantID: 1,
		OrderItems: []domain.OrderItem{
			{MenuItemID: 1, Quantity: 2},
			{MenuItemID: 2, Quantity: 1},
		},
	}).Return(0, apperr.NewAppError(apperr.ErrInternal, "failed to create order", nil)).Once()
	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, dtos.CreateOrderRequest{
		RestaurantID: 1,
		OrderItems: []dtos.OrderItemsDTO{
			{MenuItemID: 1, Quantity: 2},
			{MenuItemID: 2, Quantity: 1},
		},
	})

	require.NoError(t, err, "expected no error while encoding request")
	req := httptest.NewRequest("POST", "/api/orders", buf)
	userClaim := &authctx.UserClaims{
		UserID: 1,
		Role:   "customer",
	}
	ctx := authctx.WithUserClaims(req.Context(), userClaim)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	handler.HandleCreateOrder(w, req)

	res := w.Result()
	require.Equal(t, 500, res.StatusCode, "expected status code 500")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)
	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 500, errorResponse.Status, "expected error status to be 500")
	require.Contains(t, errorResponse.Message, "internal server error", "expected error message to contain 'internal server error'")
}

func Test_handlers_OrdersHandler_HandleCreateOrder_BadRequestFromService(t *testing.T) {
	mockOrderService := &mockservice.OrderService{}
	handler := NewOrdersHandler(mockOrderService)
	require.NotNil(t, handler, "expected NewOrdersHandler to return a non-nil handler")

	mockOrderService.On("CreateOrder", mock.Anything, domain.Order{
		CustomerID:   1,
		RestaurantID: 1,
		OrderItems: []domain.OrderItem{
			{MenuItemID: 1, Quantity: 2},
			{MenuItemID: 2, Quantity: 1},
		},
	}).Return(0, apperr.NewAppError(apperr.ErrInvalid, "invalid order data", nil)).Once()
	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, dtos.CreateOrderRequest{
		RestaurantID: 1,
		OrderItems: []dtos.OrderItemsDTO{
			{MenuItemID: 1, Quantity: 2},
			{MenuItemID: 2, Quantity: 1},
		},
	})

	require.NoError(t, err, "expected no error while encoding request")
	req := httptest.NewRequest("POST", "/api/orders", buf)
	userClaim := &authctx.UserClaims{
		UserID: 1,
		Role:   "customer",
	}
	ctx := authctx.WithUserClaims(req.Context(), userClaim)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	handler.HandleCreateOrder(w, req)

	res := w.Result()
	require.Equal(t, 400, res.StatusCode, "expected status code 400")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)
	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 400, errorResponse.Status, "expected error status to be 400")
	require.Contains(t, errorResponse.Message, "invalid order data", "expected error message to contain 'invalid order data'")
}

func Test_handlers_OrdersHandler_HandleCreateOrder_UnauthorizedFromService(t *testing.T) {
	mockOrderService := &mockservice.OrderService{}
	handler := NewOrdersHandler(mockOrderService)
	require.NotNil(t, handler, "expected NewOrdersHandler to return a non-nil handler")

	mockOrderService.On("CreateOrder", mock.Anything, domain.Order{
		CustomerID:   1,
		RestaurantID: 1,
		OrderItems: []domain.OrderItem{
			{MenuItemID: 1, Quantity: 2},
			{MenuItemID: 2, Quantity: 1},
		},
	}).Return(0, apperr.NewAppError(apperr.ErrUnauthorized, "unauthorized", nil)).Once()
	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, dtos.CreateOrderRequest{
		RestaurantID: 1,
		OrderItems: []dtos.OrderItemsDTO{
			{MenuItemID: 1, Quantity: 2},
			{MenuItemID: 2, Quantity: 1},
		},
	})

	require.NoError(t, err, "expected no error while encoding request")
	req := httptest.NewRequest("POST", "/api/orders", buf)
	userClaim := &authctx.UserClaims{
		UserID: 1,
		Role:   "customer",
	}
	ctx := authctx.WithUserClaims(req.Context(), userClaim)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	handler.HandleCreateOrder(w, req)

	res := w.Result()
	require.Equal(t, 401, res.StatusCode, "expected status code 401")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)
	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 401, errorResponse.Status, "expected error status to be 401")
	require.Contains(t, errorResponse.Message, "unauthorized", "expected error message to contain 'unauthorized'")
}

func Test_handlers_OrdersHandler_HandleGetOrderById(t *testing.T) {
	mockOrderService := &mockservice.OrderService{}
	handler := NewOrdersHandler(mockOrderService)
	require.NotNil(t, handler, "expected NewOrdersHandler to return a non-nil handler")

	mockOrderService.On("GetOrderById", mock.Anything, 1).Return(domain.Order{
		ID:           1,
		CustomerID:   1,
		RestaurantID: 1,
		OrderItems: []domain.OrderItem{
			{MenuItemID: 1, Quantity: 2},
			{MenuItemID: 2, Quantity: 1},
		},
	}, nil).Once()

	req := httptest.NewRequest("GET", "/api/orders", nil)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleGetOrderById(w, req)
	res := w.Result()

	require.Equal(t, 200, res.StatusCode, "expected status code 200")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	getOrderResp, err := decodeResponse[dtos.GetOrderByIdResponse](res)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 1, getOrderResp.ID, "expected order ID to be 1")
	require.Equal(t, 1, getOrderResp.CustomerID, "expected customer ID to be 1")
	require.Equal(t, 1, getOrderResp.RestaurantID, "expected restaurant ID to be 1")
	require.Len(t, getOrderResp.OrderItems, 2, "expected 2 order items")
}

func Test_handlers_OrdersHandler_HandleGetOrderById_NotFound(t *testing.T) {
	mockOrderService := &mockservice.OrderService{}
	handler := NewOrdersHandler(mockOrderService)
	require.NotNil(t, handler, "expected NewOrdersHandler to return a non-nil handler")

	mockOrderService.On("GetOrderById", mock.Anything, 1).Return(domain.Order{}, apperr.NewAppError(apperr.ErrNotFound, "order not found", nil)).Once()

	req := httptest.NewRequest("GET", "/api/orders", nil)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleGetOrderById(w, req)
	res := w.Result()

	require.Equal(t, 404, res.StatusCode, "expected status code 404")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 404, errorResponse.Status, "expected error status to be 404")
	require.Contains(t, errorResponse.Message, "order not found", "expected error message to contain 'order not found'")
}

func Test_handlers_OrdersHandler_HandleGetOrderById_Unauthorized(t *testing.T) {
	mockOrderService := &mockservice.OrderService{}
	handler := NewOrdersHandler(mockOrderService)
	require.NotNil(t, handler, "expected NewOrdersHandler to return a non-nil handler")

	mockOrderService.On("GetOrderById", mock.Anything, 1).Return(domain.Order{}, apperr.NewAppError(apperr.ErrUnauthorized, "unauthorized", nil)).Once()

	req := httptest.NewRequest("GET", "/api/orders", nil)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleGetOrderById(w, req)
	res := w.Result()

	require.Equal(t, 401, res.StatusCode, "expected status code 401")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 401, errorResponse.Status, "expected error status to be 401")
	require.Contains(t, errorResponse.Message, "unauthorized", "expected error message to contain 'unauthorized'")
}

func Test_handlers_OrdersHandler_HandleGetOrderById_InvalidId(t *testing.T) {
	mockOrderService := &mockservice.OrderService{}
	handler := NewOrdersHandler(mockOrderService)
	require.NotNil(t, handler, "expected NewOrdersHandler to return a non-nil handler")

	req := httptest.NewRequest("GET", "/api/orders/abc", nil)

	w := httptest.NewRecorder()
	handler.HandleGetOrderById(w, req)
	res := w.Result()

	require.Equal(t, 400, res.StatusCode, "expected status code 400")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 400, errorResponse.Status, "expected error status to be 400")
	require.Contains(t, errorResponse.Message, "invalid order id", "expected error message to contain 'invalid order id'")
}

func Test_handlers_OrdersHandler_HandleGetOrderById_InternalServerError(t *testing.T) {
	mockOrderService := &mockservice.OrderService{}
	handler := NewOrdersHandler(mockOrderService)
	require.NotNil(t, handler, "expected NewOrdersHandler to return a non-nil handler")

	mockOrderService.On("GetOrderById", mock.Anything, 1).Return(domain.Order{}, apperr.NewAppError(apperr.ErrInternal, "failed to get order", nil)).Once()

	req := httptest.NewRequest("GET", "/api/orders", nil)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleGetOrderById(w, req)
	res := w.Result()

	require.Equal(t, 500, res.StatusCode, "expected status code 500")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 500, errorResponse.Status, "expected error status to be 500")
	require.Contains(t, errorResponse.Message, "internal server error", "expected error message to contain 'internal server error'")
}

func Test_handlers_OrdersHandler_HandleGetOrderById_Forbidden(t *testing.T) {
	mockOrderService := &mockservice.OrderService{}
	handler := NewOrdersHandler(mockOrderService)
	require.NotNil(t, handler, "expected NewOrdersHandler to return a non-nil handler")

	mockOrderService.On("GetOrderById", mock.Anything, 1).Return(domain.Order{}, apperr.NewAppError(apperr.ErrForbidden, "forbidden", nil)).Once()

	req := httptest.NewRequest("GET", "/api/orders", nil)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleGetOrderById(w, req)
	res := w.Result()

	require.Equal(t, 403, res.StatusCode, "expected status code 403")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 403, errorResponse.Status, "expected error status to be 403")
	require.Contains(t, errorResponse.Message, "forbidden", "expected error message to contain 'forbidden'")
}

func Test_handlers_OrdersHandler_HandleAddOrderItem(t *testing.T) {
	mockOrderService := &mockservice.OrderService{}
	handler := NewOrdersHandler(mockOrderService)
	require.NotNil(t, handler, "expected NewOrdersHandler to return a non-nil handler")

	mockOrderService.On("AddOrderItem", mock.Anything, 1, domain.OrderItem{
		MenuItemID: 1,
		Quantity:   2,
	}).Return(nil).Once()
	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, dtos.AddOrderItemRequest{
		MenuItemID: 1,
		Quantity:   2,
	})
	require.NoError(t, err, "expected no error while encoding request")
	req := httptest.NewRequest("POST", "/api/orders/1/items", buf)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleAddOrderItem(w, req)
	res := w.Result()
	require.Equal(t, 200, res.StatusCode, "expected status code 200")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	addItemResp, err := decodeResponse[dtos.AddOrderItemResponse](res)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 1, addItemResp.ID, "expected order ID to be 1 as it's not set in response")
}

func Test_handlers_OrdersHandler_HandleAddOrderItem_InvalidId(t *testing.T) {
	mockOrderService := &mockservice.OrderService{}
	handler := NewOrdersHandler(mockOrderService)
	require.NotNil(t, handler, "expected NewOrdersHandler to return a non-nil handler")

	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, dtos.AddOrderItemRequest{
		MenuItemID: 1,
		Quantity:   2,
	})
	require.NoError(t, err, "expected no error while encoding request")
	req := httptest.NewRequest("POST", "/api/orders/abc/items", buf)
	req.SetPathValue("id", "abc")

	w := httptest.NewRecorder()
	handler.HandleAddOrderItem(w, req)
	res := w.Result()

	require.Equal(t, 400, res.StatusCode, "expected status code 400")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 400, errorResponse.Status, "expected error status to be 400")
	require.Contains(t, errorResponse.Message, "invalid order id", "expected error message to contain 'invalid order id'")
}

func Test_handlers_OrdersHandler_HandleAddOrderItem_InvalidRequest(t *testing.T) {
	mockOrderService := &mockservice.OrderService{}
	handler := NewOrdersHandler(mockOrderService)
	require.NotNil(t, handler, "expected NewOrdersHandler to return a non-nil handler")

	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, struct {
		MenuItemID string `json:"menu_item_id"`
		Quantity   int    `json:"quantity"`
	}{
		MenuItemID: "one",
		Quantity:   2,
	})
	require.NoError(t, err, "expected no error while encoding request")
	req := httptest.NewRequest("POST", "/api/orders/1/items", buf)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleAddOrderItem(w, req)
	res := w.Result()

	require.Equal(t, 400, res.StatusCode, "expected status code 400")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 400, errorResponse.Status, "expected error status to be 400")
	require.Contains(t, errorResponse.Message, "invalid", "expected error message to contain 'invalid'")
}

func Test_handlers_OrdersHandler_HandleAddOrderItem_InternalServerError(t *testing.T) {
	mockOrderService := &mockservice.OrderService{}
	handler := NewOrdersHandler(mockOrderService)
	require.NotNil(t, handler, "expected NewOrdersHandler to return a non-nil handler")

	mockOrderService.On("AddOrderItem", mock.Anything, 1, domain.OrderItem{
		MenuItemID: 1,
		Quantity:   2,
	}).Return(apperr.NewAppError(apperr.ErrInternal, "failed to add order item", nil)).Once()
	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, dtos.AddOrderItemRequest{
		MenuItemID: 1,
		Quantity:   2,
	})
	require.NoError(t, err, "expected no error while encoding request")
	req := httptest.NewRequest("POST", "/api/orders/1/items", buf)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleAddOrderItem(w, req)
	res := w.Result()
	require.Equal(t, 500, res.StatusCode, "expected status code 500")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 500, errorResponse.Status, "expected error status to be 500")
	require.Contains(t, errorResponse.Message, "internal server error", "expected error message to contain 'internal server error'")
}

func Test_handlers_OrdersHandler_HandleAddOrderItem_BadRequestFromService(t *testing.T) {
	mockOrderService := &mockservice.OrderService{}
	handler := NewOrdersHandler(mockOrderService)
	require.NotNil(t, handler, "expected NewOrdersHandler to return a non-nil handler")

	mockOrderService.On("AddOrderItem", mock.Anything, 1, domain.OrderItem{
		MenuItemID: 1,
		Quantity:   2,
	}).Return(apperr.NewAppError(apperr.ErrInvalid, "invalid order item data", nil)).Once()
	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, dtos.AddOrderItemRequest{
		MenuItemID: 1,
		Quantity:   2,
	})
	require.NoError(t, err, "expected no error while encoding request")
	req := httptest.NewRequest("POST", "/api/orders/1/items", buf)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleAddOrderItem(w, req)
	res := w.Result()
	require.Equal(t, 400, res.StatusCode, "expected status code 400")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 400, errorResponse.Status, "expected error status to be 400")
	require.Contains(t, errorResponse.Message, "invalid order item data", "expected error message to contain 'invalid order item data'")
}

func Test_handlers_OrdersHandler_HandleAddOrderItem_UnauthorizedFromService(t *testing.T) {
	mockOrderService := &mockservice.OrderService{}
	handler := NewOrdersHandler(mockOrderService)
	require.NotNil(t, handler, "expected NewOrdersHandler to return a non-nil handler")

	mockOrderService.On("AddOrderItem", mock.Anything, 1, domain.OrderItem{
		MenuItemID: 1,
		Quantity:   2,
	}).Return(apperr.NewAppError(apperr.ErrUnauthorized, "unauthorized", nil)).Once()
	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, dtos.AddOrderItemRequest{
		MenuItemID: 1,
		Quantity:   2,
	})
	require.NoError(t, err, "expected no error while encoding request")
	req := httptest.NewRequest("POST", "/api/orders/1/items", buf)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleAddOrderItem(w, req)
	res := w.Result()
	require.Equal(t, 401, res.StatusCode, "expected status code 401")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 401, errorResponse.Status, "expected error status to be 401")
	require.Contains(t, errorResponse.Message, "unauthorized", "expected error message to contain 'unauthorized'")
}

func Test_handlers_OrdersHandler_HandleAddOrderItem_ForbiddenFromService(t *testing.T) {
	mockOrderService := &mockservice.OrderService{}
	handler := NewOrdersHandler(mockOrderService)
	require.NotNil(t, handler, "expected NewOrdersHandler to return a non-nil handler")

	mockOrderService.On("AddOrderItem", mock.Anything, 1, domain.OrderItem{
		MenuItemID: 1,
		Quantity:   2,
	}).Return(apperr.NewAppError(apperr.ErrForbidden, "forbidden", nil)).Once()
	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, dtos.AddOrderItemRequest{
		MenuItemID: 1,
		Quantity:   2,
	})
	require.NoError(t, err, "expected no error while encoding request")
	req := httptest.NewRequest("POST", "/api/orders/1/items", buf)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleAddOrderItem(w, req)
	res := w.Result()
	require.Equal(t, 403, res.StatusCode, "expected status code 403")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 403, errorResponse.Status, "expected error status to be 403")
	require.Contains(t, errorResponse.Message, "forbidden", "expected error message to contain 'forbidden'")
}

func Test_handlers_OrdersHandler_HandleAddOrderItem_NotFoundFromService(t *testing.T) {
	mockOrderService := &mockservice.OrderService{}
	handler := NewOrdersHandler(mockOrderService)
	require.NotNil(t, handler, "expected NewOrdersHandler to return a non-nil handler")

	mockOrderService.On("AddOrderItem", mock.Anything, 1, domain.OrderItem{
		MenuItemID: 1,
		Quantity:   2,
	}).Return(apperr.NewAppError(apperr.ErrNotFound, "order not found", nil)).Once()
	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, dtos.AddOrderItemRequest{
		MenuItemID: 1,
		Quantity:   2,
	})
	require.NoError(t, err, "expected no error while encoding request")
	req := httptest.NewRequest("POST", "/api/orders/1/items", buf)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleAddOrderItem(w, req)
	res := w.Result()
	require.Equal(t, 404, res.StatusCode, "expected status code 404")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	errorResponse, err := decodeJson[dtos.BaseResponse](res.Body)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 404, errorResponse.Status, "expected error status to be 404")
	require.Contains(t, errorResponse.Message, "order not found", "expected error message to contain 'order not found'")
}
