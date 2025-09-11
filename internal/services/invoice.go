package services

import (
	"context"

	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/mohits-git/food-ordering-system/internal/ports"
	"github.com/mohits-git/food-ordering-system/internal/utils/apperr"
	"github.com/mohits-git/food-ordering-system/internal/utils/authctx"
)

type InvoiceService struct {
	invoiceRepo  ports.InvoiceRepository
	orderRepo    ports.OrderRepository
	menuItemRepo ports.MenuItemRepository
}

func NewInvoiceService(
	invoiceRepo ports.InvoiceRepository,
	orderRepo ports.OrderRepository,
	menuItemRepo ports.MenuItemRepository,
) *InvoiceService {
	return &InvoiceService{
		invoiceRepo:  invoiceRepo,
		orderRepo:    orderRepo,
		menuItemRepo: menuItemRepo,
	}
}

func (s *InvoiceService) getRestaurantItemsMap(ctx context.Context, restaurantId int) (map[int]domain.MenuItem, error) {
	restaurantItems, err := s.menuItemRepo.FindMenuItemsByRestaurantId(ctx, restaurantId)
	if err != nil {
		return nil, err
	}
	if len(restaurantItems) == 0 {
		return nil, apperr.NewAppError(apperr.ErrInternal, "internal server error", nil)
	}
	restaurantItemMap := make(map[int]domain.MenuItem)
	for _, item := range restaurantItems {
		restaurantItemMap[item.ID] = item
	}
	return restaurantItemMap, nil
}

func (s *InvoiceService) getItemsAvailabilityMap(restaurantItemsMap map[int]domain.MenuItem) map[int]bool {
	availabilityMap := make(map[int]bool)
	for id, item := range restaurantItemsMap {
		availabilityMap[id] = item.Available
	}
	return availabilityMap
}

func (s *InvoiceService) calculateTax(amount float64) float64 {
	return amount * 0.10 // 10% tax
}

func (s *InvoiceService) cancelInvoices(ctx context.Context, orderId int) error {
	allInvoices, err := s.invoiceRepo.FindInvoicesByOrderId(ctx, orderId)
	if err != nil {
		return err
	}
	for _, inv := range allInvoices {
		if inv.PaymentStatus == domain.Unpaid {
			err := s.invoiceRepo.ChangeInvoiceStatus(ctx, inv.ID, domain.Cancelled)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *InvoiceService) getOrderById(cxt context.Context, orderId int) (domain.Order, error) {
	if orderId <= 0 {
		return domain.Order{}, apperr.NewAppError(apperr.ErrInvalid, "invalid order id", nil)
	}

	order, err := s.orderRepo.FindOrderById(cxt, orderId)
	if err != nil {
		return domain.Order{}, err
	}
	if order.ID == 0 {
		return domain.Order{}, apperr.NewAppError(apperr.ErrNotFound, "order not found", nil)
	}

	return order, nil
}

func (s *InvoiceService) getTotalPrice(order domain.Order, menuItems map[int]domain.MenuItem) float64 {
	total := 0.0
	for _, item := range order.OrderItems {
		if menuItem, exists := menuItems[item.MenuItemID]; exists && menuItem.IsAvailable() {
			total += menuItem.Price * float64(item.Quantity)
		}
	}
	return total
}

func (s *InvoiceService) GenerateInvoice(ctx context.Context, orderId int) (domain.Invoice, error) {
	order, err := s.getOrderById(ctx, orderId)

	user, ok := authctx.UserClaimsFromCtx(ctx)
	if !ok {
		return domain.Invoice{}, apperr.NewAppError(apperr.ErrUnauthorized, "user not authenticated", nil)
	}
	if user.Role != domain.CUSTOMER || user.UserID != order.CustomerID {
		return domain.Invoice{}, apperr.NewAppError(apperr.ErrForbidden, "only customers can generate invoices", nil)
	}

	restaurantItemsMap, err := s.getRestaurantItemsMap(ctx, order.RestaurantID)
	if err != nil {
		return domain.Invoice{}, err
	}
	restaurantItemsAvailableMap := s.getItemsAvailabilityMap(restaurantItemsMap)

	if !order.Validate(restaurantItemsAvailableMap) {
		return domain.Invoice{}, apperr.NewAppError(apperr.ErrInvalid, "invalid order data, or item not available", nil)
	}

	// update other invoices for this order to be cancelled
	s.cancelInvoices(ctx, orderId)

	// Create an invoice based on the order details
	total := s.getTotalPrice(order, restaurantItemsMap)
	tax := s.calculateTax(total)
	invoice := domain.Invoice{
		OrderID:       order.ID,
		Total:         total,
		Tax:           tax,
		PaymentStatus: domain.Unpaid,
	}

	id, err := s.invoiceRepo.SaveInvoice(ctx, invoice)
	if err != nil {
		return domain.Invoice{}, err
	}

	savedInvoice, err := s.invoiceRepo.FindInvoiceById(ctx, id)
	if err != nil {
		return domain.Invoice{}, err
	}

	return savedInvoice, nil
}

func (s *InvoiceService) GetInvoiceById(cxt context.Context, id int) (domain.Invoice, error) {
	if id <= 0 {
		return domain.Invoice{}, apperr.NewAppError(apperr.ErrInvalid, "invalid invoice id", nil)
	}

	user, ok := authctx.UserClaimsFromCtx(cxt)
	if !ok {
		return domain.Invoice{}, apperr.NewAppError(apperr.ErrUnauthorized, "user not authenticated", nil)
	}
	if user.Role != domain.CUSTOMER {
		return domain.Invoice{}, apperr.NewAppError(apperr.ErrForbidden, "only customers can access invoices", nil)
	}

	invoice, err := s.invoiceRepo.FindInvoiceById(cxt, id)
	if err != nil {
		return domain.Invoice{}, err
	}
	if invoice.ID == 0 {
		return domain.Invoice{}, apperr.NewAppError(apperr.ErrNotFound, "invoice not found", nil)
	}

	order, err := s.orderRepo.FindOrderById(cxt, invoice.OrderID)
	if err != nil {
		return domain.Invoice{}, err
	}
	if order.CustomerID != user.UserID {
		return domain.Invoice{}, apperr.NewAppError(apperr.ErrForbidden, "access to the invoice is forbidden", nil)
	}

	return invoice, nil
}

func (s *InvoiceService) DoInvoicePayment(cxt context.Context, invoiceId int, payment float64) error {
	if invoiceId <= 0 {
		return apperr.NewAppError(apperr.ErrInvalid, "invalid invoice id", nil)
	}
	user, ok := authctx.UserClaimsFromCtx(cxt)
	if !ok {
		return apperr.NewAppError(apperr.ErrUnauthorized, "user not authenticated", nil)
	}
	if user.Role != domain.CUSTOMER {
		return apperr.NewAppError(apperr.ErrForbidden, "only customers can update invoice status", nil)
	}

	invoice, err := s.invoiceRepo.FindInvoiceById(cxt, invoiceId)
	if err != nil {
		return err
	}
	if invoice.ID == 0 {
		return apperr.NewAppError(apperr.ErrNotFound, "invoice not found", nil)
	}

	order, err := s.orderRepo.FindOrderById(cxt, invoice.OrderID)
	if err != nil {
		return err
	}
	if order.CustomerID != user.UserID {
		return apperr.NewAppError(apperr.ErrForbidden, "access to the invoice is forbidden", nil)
	}

	if invoice.PaymentStatus == domain.Paid {
		return apperr.NewAppError(apperr.ErrInvalid, "invalid request", nil)
	}

	if invoice.Tax+invoice.Total > payment {
		return apperr.NewAppError(apperr.ErrInvalid, "insufficient payment amount", nil)
	}

	err = s.invoiceRepo.ChangeInvoiceStatus(cxt, invoiceId, domain.Paid)
	if err != nil {
		return err
	}

	return nil
}
