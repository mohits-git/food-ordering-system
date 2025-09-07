package ports

import (
	"context"

	"github.com/mohits-git/food-ordering-system/internal/domain"
)

type InvoiceService interface {
	GenerateInvoice(cxt context.Context, orderId int) (domain.Invoice, error)
	GetInvoiceById(cxt context.Context, id int) (domain.Invoice, error)
	DoInvoicePayment(cxt context.Context, invoiceId int, payment float64) error
}
