package ports

import (
	"context"
	"github.com/mohits-git/food-ordering-system/internal/domain"
)

type InvoiceRepository interface {
	SaveInvoice(cxt context.Context, invoice domain.Invoice) error
	UpdateInvoice(cxt context.Context, invoice domain.Invoice) error
	FindInvoiceById(cxt context.Context, id int) (domain.Invoice, error)
}
