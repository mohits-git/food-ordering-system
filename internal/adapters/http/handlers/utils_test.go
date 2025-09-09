package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mohits-git/food-ordering-system/internal/adapters/http/dtos"
	"github.com/stretchr/testify/require"
)

func Test_handlers_encodeJson(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, dtos.BaseResponse{
		Status:  200,
		Message: "success",
	})
	require.NoError(t, err, "expected no error while encoding json")

	var response dtos.BaseResponse
	err = json.NewDecoder(buf).Decode(&response)
	require.NoError(t, err, "expected no error while decoding json")
	require.Equal(t, 200, response.Status, "expected status to be 200")
	require.Equal(t, "success", response.Message, "expected message to be 'success'")
}

func Test_handlers_encodeJson_Error(t *testing.T) {
	w := bytes.NewBuffer(nil)
	err := encodeJson(w, func() {})
	require.Error(t, err, "expected error while encoding json with error writer")
}

func Test_handlers_decodeJson(t *testing.T) {
	data := `{"status":200,"message":"success"}`
	buf := bytes.NewBufferString(data)
	bufCloser := io.NopCloser(buf)

	response, err := decodeJson[dtos.BaseResponse](bufCloser)
	require.NoError(t, err, "expected no error while decoding json")
	require.Equal(t, 200, response.Status, "expected status to be 200")
	require.Equal(t, "success", response.Message, "expected message to be 'success'")
}

func Test_handlers_decodeJson_Invalid(t *testing.T) {
	data := `{"status":200,"message":"success"`
	buf := bytes.NewBufferString(data)
	bufCloser := io.NopCloser(buf)

	_, err := decodeJson[dtos.BaseResponse](bufCloser)
	require.Error(t, err, "expected error while decoding invalid json")
}

func Test_handlers_decodeResponse(t *testing.T) {
	data := `{"status":200,"message":"success","data":{"status":200,"message":"inner success"}}`
	buf := bytes.NewBufferString(data)
	bufCloser := io.NopCloser(buf)

	httpResponse := &http.Response{
		StatusCode: 200,
		Body:       bufCloser,
	}

	response, err := decodeResponse[dtos.BaseResponse](httpResponse)
	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 200, response.Status, "expected status to be 200")
	require.Equal(t, "inner success", response.Message, "expected message to be 'inner success'")
}

func Test_handlers_decodeResponse_Invalid(t *testing.T) {
	data := `{"status":"200","message":"success","data":{"status":200,"message":"inner success"`
	buf := bytes.NewBufferString(data)
	bufCloser := io.NopCloser(buf)

	httpResponse := &http.Response{
		StatusCode: 200,
		Body:       bufCloser,
	}

	_, err := decodeResponse[dtos.BaseResponse](httpResponse)
	require.Error(t, err, "expected error while decoding invalid response")
}

func Test_handlers_writeResponse(t *testing.T) {
	rw := httptest.NewRecorder()

	writeResponse(rw, 200, "success", struct{}{})

	res := rw.Result()

	require.Equal(t, 200, res.StatusCode, "expected status code to be 200")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type to be application/json")

	var response dtos.BaseResponse
	err := json.NewDecoder(res.Body).Decode(&response)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 200, response.Status, "expected status to be 200")
	require.Equal(t, "success", response.Message, "expected message to be 'success'")
}

func Test_handlers_writeError(t *testing.T) {
	rw := httptest.NewRecorder()

	writeError(rw, 500, "internal server error")

	res := rw.Result()

	require.Equal(t, 500, res.StatusCode, "expected status code to be 500")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type to be application/json")

	var response dtos.BaseResponse
	err := json.NewDecoder(res.Body).Decode(&response)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 500, response.Status, "expected status to be 500")
	require.Equal(t, "internal server error", response.Message, "expected message to be 'internal server error'")
}

func Test_handlers_getBearerToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("Authorization", "Bearer testtoken")
	token := getBearerToken(req)
	require.Equal(t, "testtoken", token, "expected token to be 'testtoken'")
}

func Test_handlers_getBearerToken_Invalid(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("Authorization ", "InvalidToken")
	token := getBearerToken(req)
	require.Equal(t, "", token, "expected token to be empty")
}

func Test_getIdFromPath(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/test/123", nil)
	req.SetPathValue("id", "123")
	id := getIdFromPath(req, "id")
	require.Equal(t, 123, id, "expected id to be 123")
}

func Test_getIdFromPath_Invalid(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/test/abc", nil)
	req.SetPathValue("id", "abc")
	id := getIdFromPath(req, "id")
	require.Equal(t, 0, id, "expected id to be 0 for invalid id")

	req = httptest.NewRequest("GET", "/api/test/-1", nil)
	req.SetPathValue("id", "-1")
	id = getIdFromPath(req, "id")
	require.Equal(t, 0, id, "expected id to be 0 for negative id")

	req = httptest.NewRequest("GET", "/api/test/", nil)
	id = getIdFromPath(req, "id")
	require.Equal(t, 0, id, "expected id to be 0 for missing id")
}
