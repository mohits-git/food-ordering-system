package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/mohits-git/food-ordering-system/internal/adapters/http/dtos"
	"github.com/mohits-git/food-ordering-system/internal/utils/apperr"
)

func encodeJson[T any](w io.Writer, data T) error {
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return apperr.NewAppError(apperr.ErrInternal, "failed to encode json response", err)
	}
	return nil
}

func decodeJson[T any](r *http.Request) (T, error) {
	var data T
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return data, apperr.NewAppError(apperr.ErrInvalid, "failed to decode json request", err)
	}
	return data, nil
}

func writeResponse[T any](w http.ResponseWriter, statusCode int, msg string, data T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encodeJson(w, dtos.NewResponse(statusCode, msg, data))
}

func writeError(w http.ResponseWriter, statusCode int, msg string) {
	writeResponse(w, statusCode, msg, struct{}{})
}
