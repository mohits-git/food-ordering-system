package mockrepository

import (
	"context"

	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/stretchr/testify/mock"
)

type InvoiceRepository struct {
	mock.Mock
}

func (i *InvoiceRepository) SaveInvoice(cxt context.Context, invoice domain.Invoice) (int, error) {
	args := i.Called(cxt, invoice)
	return args.Int(0), args.Error(1)
}

func (i *InvoiceRepository) FindInvoiceById(cxt context.Context, id int) (domain.Invoice, error) {
	args := i.Called(cxt, id)
	return args.Get(0).(domain.Invoice), args.Error(1)
}

func (i *InvoiceRepository) FindInvoicesByOrderId(ctx context.Context, orderId int) ([]domain.Invoice, error) {
	args := i.Called(ctx, orderId)
	return args.Get(0).([]domain.Invoice), args.Error(1)
}

func (i *InvoiceRepository) ChangeInvoiceStatus(cxt context.Context, invoiceId int, status domain.PaymentStatus) error {
	args := i.Called(cxt, invoiceId, status)
	return args.Error(0)
}
