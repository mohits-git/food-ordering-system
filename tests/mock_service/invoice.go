package mockservice

import (
	"context"

	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/stretchr/testify/mock"
)

type InvoiceService struct {
	mock.Mock
}

func (s *InvoiceService) GenerateInvoice(ctx context.Context, orderId int) (domain.Invoice, error) {
	args := s.Called(ctx, orderId)
	return args.Get(0).(domain.Invoice), args.Error(1)
}

func (s *InvoiceService) GetInvoiceById(ctx context.Context, id int) (domain.Invoice, error) {
	args := s.Called(ctx, id)
	return args.Get(0).(domain.Invoice), args.Error(1)
}

func (s *InvoiceService) DoInvoicePayment(ctx context.Context, invoiceId int, payment float64) error {
	args := s.Called(ctx, invoiceId, payment)
	return args.Error(0)
}
