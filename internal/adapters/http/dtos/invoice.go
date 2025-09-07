package dtos

import "github.com/mohits-git/food-ordering-system/internal/domain"

type InvoiceResponse struct {
	ID            int     `json:"id"`
	OrderID       int     `json:"order_id"`
	Total         float64 `json:"total"`
	Tax           float64 `json:"tax"`
	ToPay         float64 `json:"to_pay"`
	PaymentStatus string  `json:"payment_status"`
}

func NewInvoiceResponse(invoice domain.Invoice) InvoiceResponse {
	return InvoiceResponse{
		ID:            invoice.ID,
		OrderID:       invoice.OrderID,
		Total:         invoice.Total,
		Tax:           invoice.Tax,
		ToPay:         invoice.BillWithTax(),
		PaymentStatus: string(invoice.PaymentStatus),
	}
}

type PaymentRequest struct {
	Amount float64 `json:"amount"`
}
