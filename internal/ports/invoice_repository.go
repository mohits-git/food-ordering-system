package ports

import (
	"context"
	"github.com/mohits-git/food-ordering-system/internal/domain"
)

type InvoiceRepository interface {
	SaveInvoice(cxt context.Context, invoice domain.Invoice) (int, error)
	FindInvoiceById(cxt context.Context, id int) (domain.Invoice, error)
	ChangeInvoiceStatus(cxt context.Context, invoiceId int, status domain.PaymentStatus) error
}
