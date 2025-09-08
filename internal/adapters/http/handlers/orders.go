package handlers

import (
	"net/http"

	"github.com/mohits-git/food-ordering-system/internal/adapters/http/dtos"
	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/mohits-git/food-ordering-system/internal/ports"
	"github.com/mohits-git/food-ordering-system/internal/utils/apperr"
	"github.com/mohits-git/food-ordering-system/internal/utils/authctx"
)

type OrdersHandler struct {
	orderService ports.OrderService
}

func NewOrdersHandler(orderService ports.OrderService) *OrdersHandler {
	return &OrdersHandler{orderService}
}

func (h *OrdersHandler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	orderRequest, err := decodeRequest[dtos.CreateOrderRequest](r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	user, ok := authctx.UserClaimsFromCtx(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	order := domain.Order{
		CustomerID:   user.UserID,
		RestaurantID: orderRequest.RestaurantID,
		OrderItems:   []domain.OrderItem{},
	}
	for _, item := range orderRequest.OrderItems {
		order.OrderItems = append(order.OrderItems, item.ToDomain())
	}

	id, err := h.orderService.CreateOrder(r.Context(), order)
	if err != nil {
		if apperr.IsUnauthorizedError(err) {
			writeError(w, http.StatusUnauthorized, "unauthorized")
		} else if apperr.IsForbiddenError(err) {
			writeError(w, http.StatusForbidden, "forbidden")
		} else if apperr.IsInvalidError(err) {
			writeError(w, http.StatusBadRequest, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeResponse(w, http.StatusCreated, "order created successfully", dtos.CreateOrderResponse{ID: id})
}

func (h *OrdersHandler) HandleGetOrderById(w http.ResponseWriter, r *http.Request) {
	id := getIdFromPath(r, "id")
	if id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid order id")
		return
	}

	order, err := h.orderService.GetOrderById(r.Context(), id)
	if err != nil {
		if apperr.IsNotFoundError(err) {
			writeError(w, http.StatusNotFound, "order not found")
		} else if apperr.IsUnauthorizedError(err) {
			writeError(w, http.StatusUnauthorized, "unauthorized")
		} else if apperr.IsForbiddenError(err) {
			writeError(w, http.StatusForbidden, "forbidden")
		} else {
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	// building response
	orderItemsDTO := []dtos.OrderItemsDTO{}
	for _, item := range order.OrderItems {
		orderItemsDTO = append(orderItemsDTO, dtos.OrderItemsDTO{
			MenuItemID: item.MenuItemID,
			Quantity:   item.Quantity,
		})
	}
	resp := dtos.GetOrderByIdResponse{
		ID:           order.ID,
		CustomerID:   order.CustomerID,
		RestaurantID: order.RestaurantID,
		OrderItems:   orderItemsDTO,
	}
	writeResponse(w, http.StatusOK, "order fetched successfully", resp)
}

func (h *OrdersHandler) HandleAddOrderItem(w http.ResponseWriter, r *http.Request) {
	orderID := getIdFromPath(r, "id")
	if orderID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid order id")
		return
	}

	addItemRequest, err := decodeRequest[dtos.AddOrderItemRequest](r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	orderItem := domain.OrderItem{
		MenuItemID: addItemRequest.MenuItemID,
		Quantity:   addItemRequest.Quantity,
	}

	err = h.orderService.AddOrderItem(r.Context(), orderID, orderItem)
	if err != nil {
		if apperr.IsNotFoundError(err) {
			writeError(w, http.StatusNotFound, "order not found")
		} else if apperr.IsUnauthorizedError(err) {
			writeError(w, http.StatusUnauthorized, "unauthorized")
		} else if apperr.IsForbiddenError(err) {
			writeError(w, http.StatusForbidden, "forbidden")
		} else if apperr.IsInvalidError(err) {
			writeError(w, http.StatusBadRequest, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeResponse(w, http.StatusOK, "item added to order successfully", dtos.AddOrderItemResponse{ID: orderID})
}
