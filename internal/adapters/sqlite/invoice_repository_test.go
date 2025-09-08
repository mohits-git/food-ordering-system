package sqlite

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mattn/go-sqlite3"
	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/stretchr/testify/require"
)

func Test_sqlite_InvoiceRepository_NewInvoiceRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Expected no error when creating sqlmock")
	defer db.Close()
	_ = mock

	repo := NewInvoiceRepository(db)
	require.NotNil(t, repo, "Expected NewInvoiceRepository to return a non-nil repository")
}

func Test_sqlite_InvoiceRepository_SaveInvoice(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Expected no error when creating sqlmock")
	defer db.Close()

	repo := NewInvoiceRepository(db)
	require.NotNil(t, repo, "Expected NewInvoiceRepository to return a non-nil repository")

	// Define test cases
	tests := []struct {
		name             string
		invoice          domain.Invoice
		mockSetup        func()
		expectedID       int
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			name: "Successful insert",
			invoice: domain.Invoice{
				OrderID:       1,
				Total:         100.00,
				Tax:           10.00,
				PaymentStatus: domain.Unpaid,
			},
			mockSetup: func() {
				mock.ExpectQuery("INSERT INTO invoices").
					WithArgs(1, 10000, 1000, domain.Unpaid).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			expectedID:    1,
			expectedError: false,
		},
		{
			name: "Database error",
			invoice: domain.Invoice{
				OrderID:       1,
				Total:         100.00,
				Tax:           10.00,
				PaymentStatus: domain.Unpaid,
			},
			mockSetup: func() {
				mock.ExpectQuery("INSERT INTO invoices").
					WithArgs(1, 10000, 1000, domain.Unpaid).
					WillReturnError(sqlite3.Error{Code: sqlite3.ErrConstraint, ExtendedCode: sqlite3.ErrConstraintUnique})
			},
			expectedID:       0,
			expectedError:    true,
			expectedErrorMsg: "unique constraint violation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			id, err := repo.SaveInvoice(t.Context(), tt.invoice)
			if tt.expectedError {
				require.Error(t, err)
				if tt.expectedErrorMsg != "" {
					require.Contains(t, err.Error(), tt.expectedErrorMsg)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedID, id)
			}

			err = mock.ExpectationsWereMet()
			require.NoError(t, err, "There were unfulfilled expectations")
		})
	}
}

func Test_sqlite_InvoiceRepository_FindInvoiceById(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Expected no error when creating sqlmock")
	defer db.Close()

	repo := NewInvoiceRepository(db)
	require.NotNil(t, repo, "Expected NewInvoiceRepository to return a non-nil repository")

	// Define test cases
	tests := []struct {
		name             string
		invoiceID        int
		mockSetup        func()
		expectedInvoice  domain.Invoice
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			name:      "Successful fetch",
			invoiceID: 1,
			mockSetup: func() {
				mock.ExpectQuery("SELECT id, order_id, total, tax, payment_status FROM invoices WHERE id = ?").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "total", "tax", "payment_status"}).
						AddRow(1, 1, 10000, 1000, domain.Unpaid))
			},
			expectedInvoice: domain.Invoice{
				ID:            1,
				OrderID:       1,
				Total:         100.00,
				Tax:           10.00,
				PaymentStatus: domain.Unpaid,
			},
			expectedError: false,
		},
		{
			name:      "Invoice not found",
			invoiceID: 2,
			mockSetup: func() {
				mock.ExpectQuery("SELECT id, order_id, total, tax, payment_status FROM invoices WHERE id = ?").
					WithArgs(2).
					WillReturnError(sql.ErrNoRows)
			},
			expectedInvoice: domain.Invoice{},
			expectedError:   false,
		},
		{
			name:      "Database error",
			invoiceID: 3,
			mockSetup: func() {
				mock.ExpectQuery("SELECT id, order_id, total, tax, payment_status FROM invoices WHERE id = ?").
					WithArgs(3).
					WillReturnError(sqlite3.Error{Code: sqlite3.ErrConstraint, ExtendedCode: sqlite3.ErrConstraintUnique})
			},
			expectedInvoice:  domain.Invoice{},
			expectedError:    true,
			expectedErrorMsg: "unique constraint violation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			invoice, err := repo.FindInvoiceById(t.Context(), tt.invoiceID)
			if tt.expectedError {
				require.Error(t, err)
				if tt.expectedErrorMsg != "" {
					require.Contains(t, err.Error(), tt.expectedErrorMsg)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedInvoice, invoice)
			}

			err = mock.ExpectationsWereMet()
			require.NoError(t, err, "There were unfulfilled expectations")
		})
	}
}

func Test_sqlite_InvoiceRepository_ChangeInvoiceStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Expected no error when creating sqlmock")
	defer db.Close()

	repo := NewInvoiceRepository(db)
	require.NotNil(t, repo, "Expected NewInvoiceRepository to return a non-nil repository")

	// Define test cases
	tests := []struct {
		name             string
		invoiceID        int
		newStatus        domain.PaymentStatus
		mockSetup        func()
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			name:      "Successful update",
			invoiceID: 1,
			newStatus: domain.Paid,
			mockSetup: func() {
				mock.ExpectExec("UPDATE invoices SET payment_status").
					WithArgs(domain.Paid, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name:      "Database error",
			invoiceID: 2,
			newStatus: domain.Paid,
			mockSetup: func() {
				mock.ExpectExec("UPDATE invoices SET payment_status").
					WithArgs(domain.Paid, 2).
					WillReturnError(sqlite3.Error{Code: sqlite3.ErrConstraint, ExtendedCode: sqlite3.ErrConstraintUnique})
			},
			expectedError:    true,
			expectedErrorMsg: "unique constraint violation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			err := repo.ChangeInvoiceStatus(t.Context(), tt.invoiceID, tt.newStatus)
			if tt.expectedError {
				require.Error(t, err)
				if tt.expectedErrorMsg != "" {
					require.Contains(t, err.Error(), tt.expectedErrorMsg)
				}
			} else {
				require.NoError(t, err)
			}

			err = mock.ExpectationsWereMet()
			require.NoError(t, err, "There were unfulfilled expectations")
		})
	}
}

func Test_sqlite_InvoiceRepository_FindInvoicesByOrderId(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Expected no error when creating sqlmock")
	defer db.Close()

	repo := NewInvoiceRepository(db)
	require.NotNil(t, repo, "Expected NewInvoiceRepository to return a non-nil repository")

	// Define test cases
	tests := []struct {
		name             string
		orderID          int
		mockSetup        func()
		expectedInvoices []domain.Invoice
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			name:    "Successful fetch",
			orderID: 1,
			mockSetup: func() {
				mock.ExpectQuery("SELECT id, order_id, total, tax, payment_status FROM invoices WHERE order_id").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "total", "tax", "payment_status"}).
						AddRow(1, 1, 10000, 1000, domain.Unpaid).
						AddRow(2, 1, 20000, 2000, domain.Paid))
			},
			expectedInvoices: []domain.Invoice{
				{
					ID:            1,
					OrderID:       1,
					Total:         100.00,
					Tax:           10.00,
					PaymentStatus: domain.Unpaid,
				},
				{
					ID:            2,
					OrderID:       1,
					Total:         200.00,
					Tax:           20.00,
					PaymentStatus: domain.Paid,
				},
			},
			expectedError: false,
		},
		{
			name:    "No invoices found",
			orderID: 2,
			mockSetup: func() {
				mock.ExpectQuery("SELECT id, order_id, total, tax, payment_status FROM invoices WHERE order_id").
					WithArgs(2).
					WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "total", "tax", "payment_status"}))
			},
			expectedInvoices: []domain.Invoice{},
			expectedError:    false,
		},
		{
			name:    "Database error",
			orderID: 3,
			mockSetup: func() {
				mock.ExpectQuery("SELECT id, order_id, total, tax, payment_status FROM invoices WHERE order_id").
					WithArgs(3).
					WillReturnError(sqlite3.Error{Code: sqlite3.ErrConstraint, ExtendedCode: sqlite3.ErrConstraintUnique})
			},
			expectedInvoices: nil,
			expectedError:    true,
			expectedErrorMsg: "unique constraint violation",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			invoices, err := repo.FindInvoicesByOrderId(t.Context(), tt.orderID)
			if tt.expectedError {
				require.Error(t, err)
				if tt.expectedErrorMsg != "" {
					require.Contains(t, err.Error(), tt.expectedErrorMsg)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedInvoices, invoices)
			}

			err = mock.ExpectationsWereMet()
			require.NoError(t, err, "There were unfulfilled expectations")
		})
	}
}
