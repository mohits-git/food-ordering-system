package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/mohits-git/food-ordering-system/internal/domain"
)

type InvoiceRepository struct {
	db *sql.DB
}

func NewInvoiceRepository(db *sql.DB) *InvoiceRepository {
	return &InvoiceRepository{db: db}
}

func (r *InvoiceRepository) SaveInvoice(cxt context.Context, invoice domain.Invoice) (int, error) {
	query := `INSERT INTO invoices (order_id, total, tax, payment_status) VALUES (?, ?, ?, ?) RETURNING id`
	var id int
	total := int(invoice.Total * 100)
	tax := int(invoice.Tax * 100)
	err := r.db.QueryRowContext(cxt, query, invoice.OrderID, total, tax, invoice.PaymentStatus).Scan(&id)
	if err != nil {
		return 0, HandleSQLiteError(err)
	}
	return id, nil
}

func (r *InvoiceRepository) FindInvoiceById(cxt context.Context, id int) (domain.Invoice, error) {
	query := `SELECT id, order_id, total, tax, payment_status FROM invoices WHERE id = ?`
	var invoice domain.Invoice
	var total int
	var tax int
	err := r.db.QueryRowContext(cxt, query, id).Scan(&invoice.ID, &invoice.OrderID, &total, &tax, &invoice.PaymentStatus)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Invoice{}, nil
		}
		return domain.Invoice{}, HandleSQLiteError(err)
	}
	invoice.Total = float64(total) / 100
	invoice.Tax = float64(tax) / 100
	return invoice, nil
}

func (r *InvoiceRepository) ChangeInvoiceStatus(cxt context.Context, invoiceId int, status domain.PaymentStatus) error {
	query := `UPDATE invoices SET payment_status = ? WHERE id = ?`
	_, err := r.db.ExecContext(cxt, query, status, invoiceId)
	if err != nil {
		return HandleSQLiteError(err)
	}
	return nil
}

func (r *InvoiceRepository) FindInvoicesByOrderId(ctx context.Context, orderId int) ([]domain.Invoice, error) {
	query := `SELECT id, order_id, total, tax, payment_status FROM invoices WHERE order_id = ?`
	rows, err := r.db.QueryContext(ctx, query, orderId)
	if err != nil {
		return nil, HandleSQLiteError(err)
	}
	defer rows.Close()

	invoices := []domain.Invoice{}
	for rows.Next() {
		var invoice domain.Invoice
		var total int
		var tax int
		err := rows.Scan(&invoice.ID, &invoice.OrderID, &total, &tax, &invoice.PaymentStatus)
		if err != nil {
			return nil, HandleSQLiteError(err)
		}
		invoice.Total = float64(total) / 100
		invoice.Tax = float64(tax) / 100
		invoices = append(invoices, invoice)
	}
	if err = rows.Err(); err != nil {
		return nil, HandleSQLiteError(err)
	}
	return invoices, nil
}
