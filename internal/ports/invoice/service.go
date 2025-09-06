package ports

import (
	"context"

	"github.com/mohits-git/food-ordering-system/internal/domain"
)

type InvoiceService interface {
  GenerateInvoice(cxt context.Context, order domain.Order) (domain.Invoice, error)
  GetInvoiceById(cxt context.Context, id string) (domain.Invoice, error)
  UpdateInvoice(cxt context.Context, invoice domain.Invoice) error
}
