package handlers

import (
	"log"
	"net/http"

	"github.com/mohits-git/food-ordering-system/internal/adapters/http/dtos"
	"github.com/mohits-git/food-ordering-system/internal/ports"
	"github.com/mohits-git/food-ordering-system/internal/utils/apperr"
)

type RestaurantHandler struct {
	restaurantService ports.RestaurantService
}

func NewRestaurantHandler(restaurantService ports.RestaurantService) *RestaurantHandler {
	return &RestaurantHandler{
		restaurantService: restaurantService,
	}
}

func (h *RestaurantHandler) HandleGetRestaurants(w http.ResponseWriter, r *http.Request) {
	restaurants, err := h.restaurantService.GetAllRestaurants(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch restaurants")
		return
	}

	restaurantDTOs := make([]dtos.RestaurantDTO, 0)
	for _, restaurant := range restaurants {
		restaurantDTOs = append(restaurantDTOs, dtos.NewRestaurantDTO(restaurant))
	}
	resp := dtos.GetRestaurantsResponse{Restaurants: restaurantDTOs}
	writeResponse(w, http.StatusOK, "restaurants fetched successfully", resp)
}

func (h *RestaurantHandler) HandleCreateRestaurant(w http.ResponseWriter, r *http.Request) {
	createReq, err := decodeJson[dtos.CreateRestaurantRequest](r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	id, err := h.restaurantService.CreateRestaurant(r.Context(), createReq.Name)
	if err != nil {
		log.Println("error creating restaurant:", err)
		if apperr.IsForbiddenError(err) {
			writeError(w, http.StatusForbidden, "forbidden, only owners can create restaurants")
		} else if apperr.IsUnauthorizedError(err) {
			writeError(w, http.StatusUnauthorized, "unauthorized, please login")
		} else if apperr.IsInvalidError(err) {
			writeError(w, http.StatusBadRequest, "invalid empty restaurant name")
		} else {
			writeError(w, http.StatusInternalServerError, "failed to create restaurant")
		}
		return
	}

	writeResponse(w, http.StatusCreated, "restaurant created successfully", dtos.CreateRestaurantResponse{ID: id})
}
