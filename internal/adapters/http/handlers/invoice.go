package handlers

import (
	"log"
	"net/http"

	"github.com/mohits-git/food-ordering-system/internal/adapters/http/dtos"
	"github.com/mohits-git/food-ordering-system/internal/ports"
	"github.com/mohits-git/food-ordering-system/internal/utils/apperr"
)

type InvoiceHandler struct {
	invoiceService ports.InvoiceService
}

func NewInvoiceHandler(invoiceService ports.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{invoiceService}
}

func (h *InvoiceHandler) HandleCreateInvoice(w http.ResponseWriter, r *http.Request) {
	orderId := getIdFromPath(r, "id")
	if orderId <= 0 {
		writeError(w, http.StatusBadRequest, "invalid order id")
		return
	}
	invoice, err := h.invoiceService.GenerateInvoice(r.Context(), orderId)
	if err != nil {
		if apperr.IsNotFoundError(err) {
			writeError(w, http.StatusNotFound, "order not found")
		} else if apperr.IsForbiddenError(err) {
			writeError(w, http.StatusForbidden, "cannot create invoice for this order")
		} else if apperr.IsConflictError(err) {
			writeError(w, http.StatusConflict, "invoice already exists for this order")
		} else {
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
	}
	resp := dtos.NewInvoiceResponse(invoice)
	writeResponse(w, http.StatusCreated, "invoice created successfully", resp)
}

func (h *InvoiceHandler) HandleGetInvoice(w http.ResponseWriter, r *http.Request) {
	invoiceId := getIdFromPath(r, "id")
	if invoiceId <= 0 {
		writeError(w, http.StatusBadRequest, "invalid invoice id")
		return
	}

	invoice, err := h.invoiceService.GetInvoiceById(r.Context(), invoiceId)
	if err != nil {
		if apperr.IsNotFoundError(err) {
			writeError(w, http.StatusNotFound, "invoice not found")
		} else if apperr.IsForbiddenError(err) {
			writeError(w, http.StatusForbidden, "cannot access this invoice")
		} else {
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	resp := dtos.NewInvoiceResponse(invoice)
	writeResponse(w, http.StatusOK, "invoice fetched successfully", resp)
}

func (h *InvoiceHandler) HandleInvoicePayment(w http.ResponseWriter, r *http.Request) {
	invoiceId := getIdFromPath(r, "id")
	if invoiceId <= 0 {
		writeError(w, http.StatusBadRequest, "invalid invoice id")
		return
	}

	paymentReq, err := decodeJson[dtos.PaymentRequest](r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err = h.invoiceService.DoInvoicePayment(r.Context(), invoiceId, paymentReq.Amount)
	if err != nil {
		log.Println("error in payment:", err)
		if apperr.IsNotFoundError(err) {
			writeError(w, http.StatusNotFound, "invoice not found")
		} else if apperr.IsForbiddenError(err) {
			writeError(w, http.StatusForbidden, "cannot update this invoice")
		} else if apperr.IsConflictError(err) {
			writeError(w, http.StatusConflict, "invoice already paid")
		} else {
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeResponse(w, http.StatusOK, "invoice payment successful", struct{}{})
}
