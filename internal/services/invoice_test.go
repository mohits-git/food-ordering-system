package services

import (
	"testing"

	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/mohits-git/food-ordering-system/internal/utils/apperr"
	"github.com/mohits-git/food-ordering-system/internal/utils/authctx"
	mockrepository "github.com/mohits-git/food-ordering-system/tests/mock_repository"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_services_InvoiceService_NewInvoiceService(t *testing.T) {
	mockInvoiceRepo := mockrepository.InvoiceRepository{}
	mockOrderRepo := mockrepository.OrderRepository{}
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	service := NewInvoiceService(&mockInvoiceRepo, &mockOrderRepo, &mockMenuItemRepo)
	require.NotNil(t, service)
}

func Test_services_InvoiceService_GenerateInvoice(t *testing.T) {
	mockInvoiceRepo := mockrepository.InvoiceRepository{}
	mockOrderRepo := mockrepository.OrderRepository{}
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	service := NewInvoiceService(&mockInvoiceRepo, &mockOrderRepo, &mockMenuItemRepo)

	userCtx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.CUSTOMER,
	})

	order := domain.Order{
		ID:           1,
		CustomerID:   1,
		RestaurantID: 1,
		OrderItems: []domain.OrderItem{
			{MenuItemID: 1, Quantity: 2},
			{MenuItemID: 2, Quantity: 1},
		},
	}

	mockOrderRepo.On("FindOrderById", mock.Anything, order.ID).
		Return(order, nil)
	mockMenuItemRepo.On("FindMenuItemsByRestaurantId", mock.Anything, order.RestaurantID).
		Return([]domain.MenuItem{
			{ID: 1, Name: "Item 1", Price: 100.0, Available: true},
			{ID: 2, Name: "Item 2", Price: 200.0, Available: true},
		}, nil)
	mockInvoiceRepo.On("FindInvoicesByOrderId", mock.Anything, order.ID).
		Return([]domain.Invoice{}, nil)
	mockInvoiceRepo.On("SaveInvoice", mock.Anything, mock.MatchedBy(func(inv domain.Invoice) bool {
		return inv.OrderID == order.ID && inv.Total == 400.0 && inv.Tax == 40.0 && inv.Total+inv.Tax == 440.0
	})).
		Return(1, nil)
	mockInvoiceRepo.On("FindInvoiceById", mock.Anything, 1).
		Return(domain.Invoice{ID: 1, OrderID: order.ID, Total: 400.0, Tax: 40.0, PaymentStatus: domain.Unpaid}, nil)

	invoice, err := service.GenerateInvoice(userCtx, order.ID)
	require.NoError(t, err)
	require.Equal(t, 1, invoice.ID)
	require.Equal(t, order.ID, invoice.OrderID)
	require.Equal(t, 400.0, invoice.Total)
	require.Equal(t, 40.0, invoice.Tax)
	require.Equal(t, domain.Unpaid, invoice.PaymentStatus)

	mockOrderRepo.AssertExpectations(t)
	mockMenuItemRepo.AssertExpectations(t)
	mockInvoiceRepo.AssertExpectations(t)
}

func Test_services_InvoiceService_GenerateInvoice_Unauthenticated(t *testing.T) {
	mockInvoiceRepo := mockrepository.InvoiceRepository{}
	mockOrderRepo := mockrepository.OrderRepository{}
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	service := NewInvoiceService(&mockInvoiceRepo, &mockOrderRepo, &mockMenuItemRepo)

	orderId := 1

	mockOrderRepo.On("FindOrderById", mock.Anything, orderId).
		Return(domain.Order{}, nil)

	_, err := service.GenerateInvoice(t.Context(), orderId)
	appErr, ok := err.(*apperr.AppError)
	require.Error(t, err)
	require.True(t, ok)
	require.Equal(t, apperr.ErrUnauthorized, appErr.Code)

	mockMenuItemRepo.AssertNotCalled(t, "FindMenuItemsByRestaurantId", mock.Anything, mock.Anything)
	mockInvoiceRepo.AssertNotCalled(t, "FindInvoicesByOrderId", mock.Anything, orderId)
	mockInvoiceRepo.AssertNotCalled(t, "SaveInvoice", mock.Anything, mock.Anything)
	mockOrderRepo.AssertExpectations(t)
}

func Test_services_InvoiceService_GenerateInvoice_Forbidden(t *testing.T) {
	mockInvoiceRepo := mockrepository.InvoiceRepository{}
	mockOrderRepo := mockrepository.OrderRepository{}
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	service := NewInvoiceService(&mockInvoiceRepo, &mockOrderRepo, &mockMenuItemRepo)

	userCtx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 2,
		Role:   domain.CUSTOMER,
	})

	order := domain.Order{
		ID:           1,
		CustomerID:   1,
		RestaurantID: 1,
		OrderItems: []domain.OrderItem{
			{MenuItemID: 1, Quantity: 2},
			{MenuItemID: 2, Quantity: 1},
		},
	}

	mockOrderRepo.On("FindOrderById", mock.Anything, order.ID).
		Return(order, nil)

	invoice, err := service.GenerateInvoice(userCtx, order.ID)
	appErr, ok := err.(*apperr.AppError)
	require.Error(t, err)
	require.True(t, ok)
	require.Equal(t, apperr.ErrForbidden, appErr.Code)
	require.Equal(t, domain.Invoice{}, invoice)

	mockOrderRepo.AssertExpectations(t)
	mockMenuItemRepo.AssertNotCalled(t, "FindMenuItemsByRestaurantId", mock.Anything, mock.Anything)
	mockInvoiceRepo.AssertNotCalled(t, "FindInvoicesByOrderId", mock.Anything, order.ID)
	mockInvoiceRepo.AssertNotCalled(t, "SaveInvoice", mock.Anything, mock.Anything)
}

func Test_services_InvoiceService_GetInvoiceById(t *testing.T) {
	mockInvoiceRepo := mockrepository.InvoiceRepository{}
	mockOrderRepo := mockrepository.OrderRepository{}
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	service := NewInvoiceService(&mockInvoiceRepo, &mockOrderRepo, &mockMenuItemRepo)

	userCtx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.CUSTOMER,
	})

	invoice := domain.Invoice{
		ID:            1,
		OrderID:       1,
		Total:         400.0,
		Tax:           40.0,
		PaymentStatus: domain.Unpaid,
	}

	mockInvoiceRepo.On("FindInvoiceById", mock.Anything, invoice.ID).
		Return(invoice, nil)
	mockOrderRepo.On("FindOrderById", mock.Anything, invoice.OrderID).
		Return(domain.Order{ID: invoice.OrderID, CustomerID: 1}, nil)

	result, err := service.GetInvoiceById(userCtx, invoice.ID)
	require.NoError(t, err)
	require.Equal(t, invoice, result)

	mockInvoiceRepo.AssertExpectations(t)
	mockOrderRepo.AssertExpectations(t)
}

func Test_services_InvoiceService_GetInvoiceById_Unauthenticated(t *testing.T) {
	mockInvoiceRepo := mockrepository.InvoiceRepository{}
	mockOrderRepo := mockrepository.OrderRepository{}
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	service := NewInvoiceService(&mockInvoiceRepo, &mockOrderRepo, &mockMenuItemRepo)

	invoiceId := 1

	_, err := service.GetInvoiceById(t.Context(), invoiceId)
	appErr, ok := err.(*apperr.AppError)
	require.Error(t, err)
	require.True(t, ok)
	require.Equal(t, apperr.ErrUnauthorized, appErr.Code)

	mockInvoiceRepo.AssertNotCalled(t, "FindInvoiceById", mock.Anything, invoiceId)
	mockOrderRepo.AssertNotCalled(t, "FindOrderById", mock.Anything, mock.Anything)
}

func Test_services_InvoiceService_GetInvoiceById_Forbidden(t *testing.T) {
	mockInvoiceRepo := mockrepository.InvoiceRepository{}
	mockOrderRepo := mockrepository.OrderRepository{}
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	service := NewInvoiceService(&mockInvoiceRepo, &mockOrderRepo, &mockMenuItemRepo)

	userCtx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 2,
		Role:   domain.CUSTOMER,
	})

	invoice := domain.Invoice{
		ID:            1,
		OrderID:       1,
		Total:         400.0,
		Tax:           40.0,
		PaymentStatus: domain.Unpaid,
	}

	mockInvoiceRepo.On("FindInvoiceById", mock.Anything, invoice.ID).
		Return(invoice, nil)
	mockOrderRepo.On("FindOrderById", mock.Anything, invoice.OrderID).
		Return(domain.Order{ID: invoice.OrderID, CustomerID: 1}, nil)

	result, err := service.GetInvoiceById(userCtx, invoice.ID)
	appErr, ok := err.(*apperr.AppError)
	require.Error(t, err)
	require.True(t, ok)
	require.Equal(t, apperr.ErrForbidden, appErr.Code)
	require.Equal(t, domain.Invoice{}, result)

	mockInvoiceRepo.AssertExpectations(t)
	mockOrderRepo.AssertExpectations(t)
}

func Test_services_InvoiceService_DoInvoicePayment(t *testing.T) {
	mockInvoiceRepo := mockrepository.InvoiceRepository{}
	mockOrderRepo := mockrepository.OrderRepository{}
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	service := NewInvoiceService(&mockInvoiceRepo, &mockOrderRepo, &mockMenuItemRepo)

	userCtx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.CUSTOMER,
	})

	invoice := domain.Invoice{
		ID:            1,
		OrderID:       1,
		Total:         400.0,
		Tax:           40.0,
		PaymentStatus: domain.Unpaid,
	}

	mockInvoiceRepo.On("FindInvoiceById", mock.Anything, invoice.ID).
		Return(invoice, nil)
	mockOrderRepo.On("FindOrderById", mock.Anything, invoice.OrderID).
		Return(domain.Order{ID: invoice.OrderID, CustomerID: 1}, nil)
	mockInvoiceRepo.On("ChangeInvoiceStatus", mock.Anything, invoice.ID, domain.Paid).
		Return(nil)

	err := service.DoInvoicePayment(userCtx, invoice.ID, 440.0)
	require.NoError(t, err)

	mockInvoiceRepo.AssertExpectations(t)
	mockOrderRepo.AssertExpectations(t)
}

func Test_services_InvoiceService_DoInvoicePayment_Unauthenticated(t *testing.T) {
	mockInvoiceRepo := mockrepository.InvoiceRepository{}
	mockOrderRepo := mockrepository.OrderRepository{}
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	service := NewInvoiceService(&mockInvoiceRepo, &mockOrderRepo, &mockMenuItemRepo)

	invoiceId := 1
	payment := 440.0

	err := service.DoInvoicePayment(t.Context(), invoiceId, payment)
	appErr, ok := err.(*apperr.AppError)
	require.Error(t, err)
	require.True(t, ok)
	require.Equal(t, apperr.ErrUnauthorized, appErr.Code)

	mockInvoiceRepo.AssertNotCalled(t, "FindInvoiceById", mock.Anything, invoiceId)
	mockOrderRepo.AssertNotCalled(t, "FindOrderById", mock.Anything, mock.Anything)
}

func Test_services_InvoiceService_DoInvoicePayment_Forbidden(t *testing.T) {
	mockInvoiceRepo := mockrepository.InvoiceRepository{}
	mockOrderRepo := mockrepository.OrderRepository{}
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	service := NewInvoiceService(&mockInvoiceRepo, &mockOrderRepo, &mockMenuItemRepo)

	userCtx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 2,
		Role:   domain.CUSTOMER,
	})

	invoice := domain.Invoice{
		ID:            1,
		OrderID:       1,
		Total:         400.0,
		Tax:           40.0,
		PaymentStatus: domain.Unpaid,
	}

	mockInvoiceRepo.On("FindInvoiceById", mock.Anything, invoice.ID).
		Return(invoice, nil)
	mockOrderRepo.On("FindOrderById", mock.Anything, invoice.OrderID).
		Return(domain.Order{ID: invoice.OrderID, CustomerID: 1}, nil)

	err := service.DoInvoicePayment(userCtx, invoice.ID, 440.0)
	appErr, ok := err.(*apperr.AppError)
	require.Error(t, err)
	require.True(t, ok)
	require.Equal(t, apperr.ErrForbidden, appErr.Code)

	mockInvoiceRepo.AssertExpectations(t)
	mockOrderRepo.AssertExpectations(t)
}

func Test_services_InvoiceService_DoInvoicePayment_InsufficientPayment(t *testing.T) {
	mockInvoiceRepo := mockrepository.InvoiceRepository{}
	mockOrderRepo := mockrepository.OrderRepository{}
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	service := NewInvoiceService(&mockInvoiceRepo, &mockOrderRepo, &mockMenuItemRepo)

	userCtx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.CUSTOMER,
	})

	invoice := domain.Invoice{
		ID:            1,
		OrderID:       1,
		Total:         400.0,
		Tax:           40.0,
		PaymentStatus: domain.Unpaid,
	}

	mockInvoiceRepo.On("FindInvoiceById", mock.Anything, invoice.ID).
		Return(invoice, nil)
	mockOrderRepo.On("FindOrderById", mock.Anything, invoice.OrderID).
		Return(domain.Order{ID: invoice.OrderID, CustomerID: 1}, nil)

	err := service.DoInvoicePayment(userCtx, invoice.ID, 300.0)
	appErr, ok := err.(*apperr.AppError)
	require.Error(t, err)
	require.True(t, ok)
	require.Equal(t, apperr.ErrInvalid, appErr.Code)
	require.Equal(t, "insufficient payment amount", appErr.Message)

	mockInvoiceRepo.AssertExpectations(t)
	mockOrderRepo.AssertExpectations(t)
}

func Test_services_InvoiceService_DoInvoicePayment_AlreadyPaid(t *testing.T) {
	mockInvoiceRepo := mockrepository.InvoiceRepository{}
	mockOrderRepo := mockrepository.OrderRepository{}
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	service := NewInvoiceService(&mockInvoiceRepo, &mockOrderRepo, &mockMenuItemRepo)

	userCtx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.CUSTOMER,
	})

	invoice := domain.Invoice{
		ID:            1,
		OrderID:       1,
		Total:         400.0,
		Tax:           40.0,
		PaymentStatus: domain.Paid,
	}

	mockInvoiceRepo.On("FindInvoiceById", mock.Anything, invoice.ID).
		Return(invoice, nil)
	mockOrderRepo.On("FindOrderById", mock.Anything, invoice.OrderID).
		Return(domain.Order{ID: invoice.OrderID, CustomerID: 1}, nil)

	err := service.DoInvoicePayment(userCtx, invoice.ID, 440.0)
	appErr, ok := err.(*apperr.AppError)
	require.Error(t, err)
	require.True(t, ok)
	require.Equal(t, apperr.ErrInvalid, appErr.Code)

	mockInvoiceRepo.AssertExpectations(t)
	mockOrderRepo.AssertExpectations(t)
}

func Test_services_InvoiceService_DoInvoicePayment_InvalidInvoiceId(t *testing.T) {
	mockInvoiceRepo := mockrepository.InvoiceRepository{}
	mockOrderRepo := mockrepository.OrderRepository{}
	mockMenuItemRepo := mockrepository.MenuItemRepository{}
	service := NewInvoiceService(&mockInvoiceRepo, &mockOrderRepo, &mockMenuItemRepo)

	userCtx := authctx.WithUserClaims(t.Context(), &authctx.UserClaims{
		UserID: 1,
		Role:   domain.CUSTOMER,
	})

	err := service.DoInvoicePayment(userCtx, 0, 440.0)
	appErr, ok := err.(*apperr.AppError)
	require.Error(t, err)
	require.True(t, ok)
	require.Equal(t, apperr.ErrInvalid, appErr.Code)
	require.Equal(t, "invalid invoice id", appErr.Message)

	mockInvoiceRepo.AssertNotCalled(t, "FindInvoiceById", mock.Anything, mock.Anything)
	mockOrderRepo.AssertNotCalled(t, "FindOrderById", mock.Anything, mock.Anything)
}
