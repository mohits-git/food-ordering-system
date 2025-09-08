package apiclient

import (
	"encoding/json"
	"io"

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

func decodeResponse[T any](body io.ReadCloser) (T, error) {

	type Response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    T      `json:"data"`
	}

	response, err := decodeJson[Response](body)
	if err != nil {
		var data T
		return data, err
	}

	return response.Data, nil

}

type ErrorResponse struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Data    struct{} `json:"data"`
}

func decodeError(body io.ReadCloser) (ErrorResponse, error) {
	response, err := decodeJson[ErrorResponse](body)
	if err != nil {
		return ErrorResponse{}, err
	}

	return response, nil
}
