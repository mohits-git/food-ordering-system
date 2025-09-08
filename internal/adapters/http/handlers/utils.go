package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/mohits-git/food-ordering-system/internal/adapters/http/dtos"
	"github.com/mohits-git/food-ordering-system/internal/utils/apperr"
)

func encodeJson[T any](w io.Writer, data T) error {
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return apperr.NewAppError(apperr.ErrInternal, "failed to encode json response", err)
	}
	return nil
}

func decodeJson[T any](body io.ReadCloser) (T, error) {
	var data T
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return data, apperr.NewAppError(apperr.ErrInvalid, "failed to decode json", err)
	}
	return data, nil
}

func decodeRequest[T any](r *http.Request) (T, error) {
	return decodeJson[T](r.Body)
}

func decodeResponse[T any](r *http.Response) (T, error) {
	type Response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    T      `json:"data"`
	}

	response, err := decodeJson[Response](r.Body)
	if err != nil {
		var data T
		return data, err
	}

	return response.Data, nil
}

func writeResponse[T any](w http.ResponseWriter, statusCode int, msg string, data T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encodeJson(w, dtos.NewResponse(statusCode, msg, data))
}

func writeError(w http.ResponseWriter, statusCode int, msg string) {
	writeResponse(w, statusCode, msg, struct{}{})
}

func getBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}

func getIdFromPath(r *http.Request, key string) int {
	idParam := r.PathValue(key)
	if idParam == "" {
		return 0
	}
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		return 0
	}
	return id
}
