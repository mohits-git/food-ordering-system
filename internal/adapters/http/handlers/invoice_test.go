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

func Test_handlers_InvoiceHandler_NewInvoiceHandler(t *testing.T) {
	mockInvoiceService := &mockservice.InvoiceService{}
	handler := NewInvoiceHandler(mockInvoiceService)
	require.NotNil(t, handler, "expected NewInvoiceHandler to return a non-nil handler")
}

func Test_handlers_InvoiceHandler_HandlerCreateInvoice(t *testing.T) {
	mockInvoiceService := &mockservice.InvoiceService{}
	handler := NewInvoiceHandler(mockInvoiceService)
	require.NotNil(t, handler, "expected NewInvoiceHandler to return a non-nil handler")

	mockInvoiceService.On("GenerateInvoice", mock.Anything, 1).Return(domain.Invoice{
		ID:            1,
		OrderID:       1,
		Total:         100.0,
		Tax:           10.0,
		PaymentStatus: "PAID",
	}, nil).Once()

	req := httptest.NewRequest("POST", "/api/orders/1/invoices", nil)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleCreateInvoice(w, req)
	res := w.Result()
	require.Equal(t, 201, res.StatusCode, "expected status code 201")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	invoice, err := decodeResponse[dtos.InvoiceResponse](res)
	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 1, invoice.ID, "expected invoice ID to be 1")
	mockInvoiceService.AssertExpectations(t)
}

func Test_handlers_InvoiceHandler_HandleCreateInvoice_BadRequest(t *testing.T) {
	mockInvoiceService := &mockservice.InvoiceService{}
	handler := NewInvoiceHandler(mockInvoiceService)
	require.NotNil(t, handler, "expected NewInvoiceHandler to return a non-nil handler")

	req := httptest.NewRequest("POST", "/api/orders/abc/invoices", nil)
	req.SetPathValue("id", "abc")

	w := httptest.NewRecorder()
	handler.HandleCreateInvoice(w, req)
	res := w.Result()
	require.Equal(t, 400, res.StatusCode, "expected status code 400")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	_, err := decodeJson[dtos.BaseResponse](res.Body)
	require.NoError(t, err, "expected error while decoding response")
	mockInvoiceService.AssertExpectations(t)
}

func Test_handlers_InvoiceHandler_HandleCreateInvoice_ServiceError(t *testing.T) {
	mockInvoiceService := &mockservice.InvoiceService{}
	handler := NewInvoiceHandler(mockInvoiceService)
	require.NotNil(t, handler, "expected NewInvoiceHandler to return a non-nil handler")

	mockInvoiceService.On("GenerateInvoice", mock.Anything, 1).Return(
		domain.Invoice{}, apperr.NewAppError(apperr.ErrInternal, "failed to generate invoice", nil)).Once()

	req := httptest.NewRequest("POST", "/api/orders/1/invoices", nil)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleCreateInvoice(w, req)
	res := w.Result()
	require.Equal(t, 500, res.StatusCode, "expected status code 500")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	_, err := decodeJson[dtos.BaseResponse](res.Body)
	require.NoError(t, err, "expected error while decoding response")
	mockInvoiceService.AssertExpectations(t)
}

func Test_handlers_InvoiceHandler_HandleCreateInvoice_NotFound(t *testing.T) {
	mockInvoiceService := &mockservice.InvoiceService{}
	handler := NewInvoiceHandler(mockInvoiceService)
	require.NotNil(t, handler, "expected NewInvoiceHandler to return a non-nil handler")

	mockInvoiceService.On("GenerateInvoice", mock.Anything, 1).Return(
		domain.Invoice{}, apperr.NewAppError(apperr.ErrNotFound, "order not found", nil)).Once()

	req := httptest.NewRequest("POST", "/api/orders/1/invoices", nil)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleCreateInvoice(w, req)
	res := w.Result()
	require.Equal(t, 404, res.StatusCode, "expected status code 404")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	_, err := decodeJson[dtos.BaseResponse](res.Body)
	require.NoError(t, err, "expected error while decoding response")
	mockInvoiceService.AssertExpectations(t)
}

func Test_handlers_InvoiceHandler_HandleCreateInvoice_Conflict(t *testing.T) {
	mockInvoiceService := &mockservice.InvoiceService{}
	handler := NewInvoiceHandler(mockInvoiceService)
	require.NotNil(t, handler, "expected NewInvoiceHandler to return a non-nil handler")

	mockInvoiceService.On("GenerateInvoice", mock.Anything, 1).Return(
		domain.Invoice{}, apperr.NewAppError(apperr.ErrConflict, "invoice already exists", nil)).Once()

	req := httptest.NewRequest("POST", "/api/orders/1/invoices", nil)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleCreateInvoice(w, req)
	res := w.Result()
	require.Equal(t, 409, res.StatusCode, "expected status code 409")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	_, err := decodeJson[dtos.BaseResponse](res.Body)
	require.NoError(t, err, "expected error while decoding response")
	mockInvoiceService.AssertExpectations(t)
}

func Test_handlers_InvoiceHandler_HandleCreateInvoice_Forbidden(t *testing.T) {
	mockInvoiceService := &mockservice.InvoiceService{}
	handler := NewInvoiceHandler(mockInvoiceService)
	require.NotNil(t, handler, "expected NewInvoiceHandler to return a non-nil handler")

	mockInvoiceService.On("GenerateInvoice", mock.Anything, 1).Return(
		domain.Invoice{}, apperr.NewAppError(apperr.ErrForbidden, "forbidden", nil)).Once()

	req := httptest.NewRequest("POST", "/api/orders/1/invoices", nil)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleCreateInvoice(w, req)
	res := w.Result()
	require.Equal(t, 403, res.StatusCode, "expected status code 403")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	_, err := decodeJson[dtos.BaseResponse](res.Body)
	require.NoError(t, err, "expected error while decoding response")
	mockInvoiceService.AssertExpectations(t)
}

func Test_handlers_InvoiceHandler_HandleCreateInvoice_Unauthorized(t *testing.T) {
	mockInvoiceService := &mockservice.InvoiceService{}
	handler := NewInvoiceHandler(mockInvoiceService)
	require.NotNil(t, handler, "expected NewInvoiceHandler to return a non-nil handler")

	mockInvoiceService.On("GenerateInvoice", mock.Anything, 1).Return(
		domain.Invoice{}, apperr.NewAppError(apperr.ErrUnauthorized, "unauthorized", nil)).Once()

	req := httptest.NewRequest("POST", "/api/orders/1/invoices", nil)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleCreateInvoice(w, req)
	res := w.Result()
	require.Equal(t, 401, res.StatusCode, "expected status code 401")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	_, err := decodeJson[dtos.BaseResponse](res.Body)
	require.NoError(t, err, "expected error while decoding response")
	mockInvoiceService.AssertExpectations(t)
}

func Test_handlers_InvoiceHandler_HandleGetInvoice(t *testing.T) {
	mockInvoiceService := &mockservice.InvoiceService{}
	handler := NewInvoiceHandler(mockInvoiceService)
	require.NotNil(t, handler, "expected NewInvoiceHandler to return a non-nil handler")

	mockInvoiceService.On("GetInvoiceById", mock.Anything, 1).Return(domain.Invoice{
		ID:            1,
		OrderID:       1,
		Total:         100.0,
		Tax:           10.0,
		PaymentStatus: "PAID",
	}, nil).Once()

	req := httptest.NewRequest("GET", "/api/invoices/1", nil)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleGetInvoice(w, req)
	res := w.Result()
	require.Equal(t, 200, res.StatusCode, "expected status code 200")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	invoice, err := decodeResponse[dtos.InvoiceResponse](res)
	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 1, invoice.ID, "expected invoice ID to be 1")
	mockInvoiceService.AssertExpectations(t)
}

func Test_handlers_InvoiceHandler_HandleGetInvoice_BadRequest(t *testing.T) {
	mockInvoiceService := &mockservice.InvoiceService{}
	handler := NewInvoiceHandler(mockInvoiceService)
	require.NotNil(t, handler, "expected NewInvoiceHandler to return a non-nil handler")

	req := httptest.NewRequest("GET", "/api/invoices/abc", nil)
	req.SetPathValue("id", "abc")

	w := httptest.NewRecorder()
	handler.HandleGetInvoice(w, req)
	res := w.Result()
	require.Equal(t, 400, res.StatusCode, "expected status code 400")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	_, err := decodeJson[dtos.BaseResponse](res.Body)
	require.NoError(t, err, "expected error while decoding response")
	mockInvoiceService.AssertExpectations(t)
}

func Test_handlers_InvoiceHandler_HandleGetInvoice_ServiceError(t *testing.T) {
	mockInvoiceService := &mockservice.InvoiceService{}
	handler := NewInvoiceHandler(mockInvoiceService)
	require.NotNil(t, handler, "expected NewInvoiceHandler to return a non-nil handler")

	mockInvoiceService.On("GetInvoiceById", mock.Anything, 1).Return(
		domain.Invoice{}, apperr.NewAppError(apperr.ErrInternal, "failed to get invoice", nil)).Once()

	req := httptest.NewRequest("GET", "/api/invoices/1", nil)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleGetInvoice(w, req)
	res := w.Result()
	require.Equal(t, 500, res.StatusCode, "expected status code 500")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	_, err := decodeJson[dtos.BaseResponse](res.Body)
	require.NoError(t, err, "expected error while decoding response")
	mockInvoiceService.AssertExpectations(t)
}

func Test_handlers_InvoiceHandler_HandleGetInvoice_NotFound(t *testing.T) {
	mockInvoiceService := &mockservice.InvoiceService{}
	handler := NewInvoiceHandler(mockInvoiceService)
	require.NotNil(t, handler, "expected NewInvoiceHandler to return a non-nil handler")

	mockInvoiceService.On("GetInvoiceById", mock.Anything, 1).Return(
		domain.Invoice{}, apperr.NewAppError(apperr.ErrNotFound, "invoice not found", nil)).Once()

	req := httptest.NewRequest("GET", "/api/invoices/1", nil)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleGetInvoice(w, req)
	res := w.Result()
	require.Equal(t, 404, res.StatusCode, "expected status code 404")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	_, err := decodeJson[dtos.BaseResponse](res.Body)
	require.NoError(t, err, "expected error while decoding response")
	mockInvoiceService.AssertExpectations(t)
}

func Test_handlers_InvoiceHandler_HandleGetInvoice_Forbidden(t *testing.T) {
	mockInvoiceService := &mockservice.InvoiceService{}
	handler := NewInvoiceHandler(mockInvoiceService)
	require.NotNil(t, handler, "expected NewInvoiceHandler to return a non-nil handler")

	mockInvoiceService.On("GetInvoiceById", mock.Anything, 1).Return(
		domain.Invoice{}, apperr.NewAppError(apperr.ErrForbidden, "forbidden", nil)).Once()

	req := httptest.NewRequest("GET", "/api/invoices/1", nil)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleGetInvoice(w, req)
	res := w.Result()
	require.Equal(t, 403, res.StatusCode, "expected status code 403")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	_, err := decodeJson[dtos.BaseResponse](res.Body)
	require.NoError(t, err, "expected error while decoding response")
	mockInvoiceService.AssertExpectations(t)
}

func Test_handlers_InvoiceHandler_HandleGetInvoice_Unauthorized(t *testing.T) {
	mockInvoiceService := &mockservice.InvoiceService{}
	handler := NewInvoiceHandler(mockInvoiceService)
	require.NotNil(t, handler, "expected NewInvoiceHandler to return a non-nil handler")

	mockInvoiceService.On("GetInvoiceById", mock.Anything, 1).Return(
		domain.Invoice{}, apperr.NewAppError(apperr.ErrUnauthorized, "unauthorized", nil)).Once()

	req := httptest.NewRequest("GET", "/api/invoices/1", nil)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleGetInvoice(w, req)
	res := w.Result()
	require.Equal(t, 401, res.StatusCode, "expected status code 401")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	_, err := decodeJson[dtos.BaseResponse](res.Body)
	require.NoError(t, err, "expected error while decoding response")
	mockInvoiceService.AssertExpectations(t)
}

func Test_handlers_InvoiceHandler_HandleInvoicePayment(t *testing.T) {
	mockInvoiceService := &mockservice.InvoiceService{}
	handler := NewInvoiceHandler(mockInvoiceService)
	require.NotNil(t, handler, "expected NewInvoiceHandler to return a non-nil handler")

	mockInvoiceService.On("DoInvoicePayment", mock.Anything, 1, 110.1).Return(nil).Once()

	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, dtos.PaymentRequest{
		Amount: 110.1,
	})

	req := httptest.NewRequest("POST", "/api/invoices/1/pay", buf)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleInvoicePayment(w, req)

	res := w.Result()
	require.Equal(t, 200, res.StatusCode, "expected status code 200")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	response, err := decodeJson[dtos.BaseResponse](res.Body)

	require.NoError(t, err, "expected no error while decoding response")
	require.Equal(t, 200, response.Status, "expected response status to be 200")
	require.Equal(t, "invoice payment successful", response.Message, "expected response message to be 'invoice payment successful'")
	mockInvoiceService.AssertExpectations(t)
}

func Test_handlers_InvoiceHandler_HandleInvoicePayment_BadRequest(t *testing.T) {
	mockInvoiceService := &mockservice.InvoiceService{}
	handler := NewInvoiceHandler(mockInvoiceService)
	require.NotNil(t, handler, "expected NewInvoiceHandler to return a non-nil handler")

	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, dtos.PaymentRequest{
		Amount: 110.1,
	})
	require.NoError(t, err, "expected no error while encoding request body")

	req := httptest.NewRequest("POST", "/api/invoices/abc/pay", buf)
	req.SetPathValue("id", "abc")

	w := httptest.NewRecorder()
	handler.HandleInvoicePayment(w, req)

	res := w.Result()
	require.Equal(t, 400, res.StatusCode, "expected status code 400")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	_, err = decodeJson[dtos.BaseResponse](res.Body)
	require.NoError(t, err, "expected error while decoding response")
	mockInvoiceService.AssertExpectations(t)
}

func Test_handlers_InvoiceHandler_HandleInvoicePayment_InvalidRequest(t *testing.T) {
	mockInvoiceService := &mockservice.InvoiceService{}
	handler := NewInvoiceHandler(mockInvoiceService)
	require.NotNil(t, handler, "expected NewInvoiceHandler to return a non-nil handler")

	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, struct {
		Amount string
	}{
		Amount: "invalid",
	})
	require.NoError(t, err, "expected no error while encoding request body")

	req := httptest.NewRequest("POST", "/api/invoices/1/pay", buf)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleInvoicePayment(w, req)

	res := w.Result()
	require.Equal(t, 400, res.StatusCode, "expected status code 400")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	_, err = decodeJson[dtos.BaseResponse](res.Body)
	require.NoError(t, err, "expected error while decoding response")
	mockInvoiceService.AssertExpectations(t)
}

func Test_handlers_InvoiceHandler_HandleInvoicePayment_ServiceError(t *testing.T) {
	mockInvoiceService := &mockservice.InvoiceService{}
	handler := NewInvoiceHandler(mockInvoiceService)
	require.NotNil(t, handler, "expected NewInvoiceHandler to return a non-nil handler")

	mockInvoiceService.On("DoInvoicePayment", mock.Anything, 1, 110.1).Return(
		apperr.NewAppError(apperr.ErrInternal, "failed to process payment", nil)).Once()

	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, dtos.PaymentRequest{
		Amount: 110.1,
	})
	require.NoError(t, err, "expected no error while encoding request body")

	req := httptest.NewRequest("POST", "/api/invoices/1/pay", buf)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleInvoicePayment(w, req)

	res := w.Result()
	require.Equal(t, 500, res.StatusCode, "expected status code 500")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	_, err = decodeJson[dtos.BaseResponse](res.Body)
	require.NoError(t, err, "expected error while decoding response")
	mockInvoiceService.AssertExpectations(t)
}

func Test_handlers_InvoiceHandler_HandleInvoicePayment_NotFound(t *testing.T) {
	mockInvoiceService := &mockservice.InvoiceService{}
	handler := NewInvoiceHandler(mockInvoiceService)
	require.NotNil(t, handler, "expected NewInvoiceHandler to return a non-nil handler")

	mockInvoiceService.On("DoInvoicePayment", mock.Anything, 1, 110.1).Return(
		apperr.NewAppError(apperr.ErrNotFound, "invoice not found", nil)).Once()

	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, dtos.PaymentRequest{
		Amount: 110.1,
	})
	require.NoError(t, err, "expected no error while encoding request body")

	req := httptest.NewRequest("POST", "/api/invoices/1/pay", buf)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleInvoicePayment(w, req)

	res := w.Result()
	require.Equal(t, 404, res.StatusCode, "expected status code 404")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	_, err = decodeJson[dtos.BaseResponse](res.Body)
	require.NoError(t, err, "expected error while decoding response")
	mockInvoiceService.AssertExpectations(t)
}

func Test_handlers_InvoiceHandler_HandleInvoicePayment_Forbidden(t *testing.T) {
	mockInvoiceService := &mockservice.InvoiceService{}
	handler := NewInvoiceHandler(mockInvoiceService)
	require.NotNil(t, handler, "expected NewInvoiceHandler to return a non-nil handler")

	mockInvoiceService.On("DoInvoicePayment", mock.Anything, 1, 110.1).Return(
		apperr.NewAppError(apperr.ErrForbidden, "forbidden", nil)).Once()

	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, dtos.PaymentRequest{
		Amount: 110.1,
	})
	require.NoError(t, err, "expected no error while encoding request body")

	req := httptest.NewRequest("POST", "/api/invoices/1/pay", buf)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleInvoicePayment(w, req)

	res := w.Result()
	require.Equal(t, 403, res.StatusCode, "expected status code 403")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	_, err = decodeJson[dtos.BaseResponse](res.Body)
	require.NoError(t, err, "expected error while decoding response")
	mockInvoiceService.AssertExpectations(t)
}

func Test_handlers_InvoiceHandler_HandleInvoicePayment_Unauthorized(t *testing.T) {
	mockInvoiceService := &mockservice.InvoiceService{}
	handler := NewInvoiceHandler(mockInvoiceService)
	require.NotNil(t, handler, "expected NewInvoiceHandler to return a non-nil handler")

	mockInvoiceService.On("DoInvoicePayment", mock.Anything, 1, 110.1).Return(
		apperr.NewAppError(apperr.ErrUnauthorized, "unauthorized", nil)).Once()

	buf := bytes.NewBuffer(nil)
	err := encodeJson(buf, dtos.PaymentRequest{
		Amount: 110.1,
	})
	require.NoError(t, err, "expected no error while encoding request body")

	req := httptest.NewRequest("POST", "/api/invoices/1/pay", buf)
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	handler.HandleInvoicePayment(w, req)

	res := w.Result()
	require.Equal(t, 401, res.StatusCode, "expected status code 401")
	require.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected content type application/json")

	defer res.Body.Close()
	_, err = decodeJson[dtos.BaseResponse](res.Body)
	require.NoError(t, err, "expected error while decoding response")
	mockInvoiceService.AssertExpectations(t)
}
