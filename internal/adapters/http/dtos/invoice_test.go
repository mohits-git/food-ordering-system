package dtos

import (
	"testing"

	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/stretchr/testify/require"
)

func Test_dtos_NewInvoiceResponse(t *testing.T) {
	response := NewInvoiceResponse(domain.Invoice{
		ID:            1,
		OrderID:       2,
		Total:         100,
		Tax:           10,
		PaymentStatus: domain.Paid,
	})

	expected := InvoiceResponse{
		ID:            1,
		OrderID:       2,
		Total:         100,
		Tax:           10,
		ToPay:         110,
		PaymentStatus: string(domain.Paid),
	}

	require.Equal(t, expected, response, "expected NewInvoiceResponse() to return expected value")
}
