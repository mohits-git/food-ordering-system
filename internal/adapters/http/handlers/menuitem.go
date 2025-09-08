package handlers

import (
	"log"
	"net/http"

	"github.com/mohits-git/food-ordering-system/internal/adapters/http/dtos"
	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/mohits-git/food-ordering-system/internal/ports"
	"github.com/mohits-git/food-ordering-system/internal/utils/apperr"
)

type MenuItemHandler struct {
	menuItemsService ports.MenuItemService
}

func NewMenuItemHandler(menuItemService ports.MenuItemService) *MenuItemHandler {
	return &MenuItemHandler{menuItemService}
}

func (h *MenuItemHandler) HandleAddMenuItemToRestaurant(w http.ResponseWriter, r *http.Request) {
	restaurantId := getIdFromPath(r, "id")
	if restaurantId <= 0 {
		writeError(w, http.StatusBadRequest, "invalid restaurant id")
		return
	}
	addRequest, err := decodeRequest[dtos.AddMenuItemRequest](r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	menuItem := domain.NewMenuItem(0, addRequest.Name, addRequest.Price, addRequest.Available, restaurantId)
	menuItemId, err := h.menuItemsService.CreateMenuItemForRestaurant(r.Context(), menuItem)
	if err != nil {
		log.Println("Error creating menu item:", err)
		if apperr.IsInvalidError(err) {
			writeError(w, http.StatusBadRequest, "invalid menu item data")
		} else if apperr.IsUnauthorizedError(err) {
			writeError(w, http.StatusUnauthorized, "unauthenticated user")
		} else if apperr.IsForbiddenError(err) {
			writeError(w, http.StatusForbidden, "only restaurant owners can add menu items")
		} else {
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	writeResponse(w, http.StatusCreated, "menu item added successfully", dtos.AddMenuItemResponse{ID: menuItemId})
}

func (h *MenuItemHandler) HandleGetRestaurantMenuItems(w http.ResponseWriter, r *http.Request) {
	restaurantId := getIdFromPath(r, "id")
	if restaurantId <= 0 {
		writeError(w, http.StatusBadRequest, "invalid restaurant id")
		return
	}

	menuItems, err := h.menuItemsService.GetAllMenuItemsByRestaurantId(r.Context(), restaurantId)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	// build response
	var itemsResponse []dtos.MenuItemResponse
	for _, item := range menuItems {
		itemsResponse = append(itemsResponse, dtos.NewMenuItemResponse(item))
	}
	response := dtos.GetMenuItemsResponse{
		RestaurantID: restaurantId,
		Items:        itemsResponse,
	}

	writeResponse(w, http.StatusOK, "menu items fetched successfully", response)
}

func (h *MenuItemHandler) HandleUpdateAvailability(w http.ResponseWriter, r *http.Request) {
	menuItemId := getIdFromPath(r, "id")
	if menuItemId <= 0 {
		writeError(w, http.StatusBadRequest, "invalid menu item id")
		return
	}

	updateRequest, err := decodeRequest[dtos.UpdateMenuItemAvailabilityRequest](r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err = h.menuItemsService.UpdateAvailability(r.Context(), menuItemId, updateRequest.Available)
	if err != nil {
		if apperr.IsInvalidError(err) {
			writeError(w, http.StatusBadRequest, "invalid menu item id")
		} else if apperr.IsUnauthorizedError(err) {
			writeError(w, http.StatusUnauthorized, "unauthenticated user")
		} else if apperr.IsForbiddenError(err) {
			writeError(w, http.StatusForbidden, "only restaurant owners can update menu items")
		} else if apperr.IsNotFoundError(err) {
			writeError(w, http.StatusNotFound, "menu item not found")
		} else {
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeResponse(w, http.StatusOK, "menu item availability updated successfully", dtos.UpdateMenuItemResponse{})
}
