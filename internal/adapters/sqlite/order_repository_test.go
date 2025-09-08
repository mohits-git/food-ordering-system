package sqlite

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_sqlite_OrderRepository_SaveOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewOrderRepository(db)

	ctx := context.Background()
	order := domain.Order{
		CustomerID:   1,
		RestaurantID: 2,
		OrderItems: []domain.OrderItem{
			{MenuItemID: 1, Quantity: 2},
			{MenuItemID: 2, Quantity: 1},
		},
	}

	// Mock the transaction begin
	mock.ExpectBegin()

	// Mock the insert into orders table
	mock.ExpectQuery("INSERT INTO orders").
		WithArgs(order.CustomerID, order.RestaurantID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Mock the insert into orderitems table
	for _, item := range order.OrderItems {
		mock.ExpectExec("INSERT INTO orderitems").
			WithArgs(1, item.MenuItemID, item.Quantity).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	// Mock the transaction commit
	mock.ExpectCommit()

	id, err := repo.SaveOrder(ctx, order)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if id != 1 {
		t.Errorf("expected order ID to be 1, got %d", id)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_sqlite_FindOrderById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewOrderRepository(db)

	ctx := context.Background()
	orderID := 1

	// Mock the select from orders table
	mock.ExpectQuery("SELECT id, user_id, restaurant_id FROM orders WHERE id = ?").
		WithArgs(orderID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "restaurant_id"}).
			AddRow(1, 1, 2))

	// Mock the select from orderitems table
	mock.ExpectQuery("SELECT menuitem_id, quantity FROM orderitems WHERE order_id = ?").
		WithArgs(orderID).
		WillReturnRows(sqlmock.NewRows([]string{"menuitem_id", "quantity"}).
			AddRow(1, 2).
			AddRow(2, 1))

	order, err := repo.FindOrderById(ctx, orderID)
	require.NoError(t, err, "unexpected error while fetching order")
	assert.Equal(t, orderID, order.ID, "expected order ID to match")
	assert.Equal(t, 2, len(order.OrderItems), "expected two order items")

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoErrorf(t, err, "there were unfulfilled expectations: %s", err)
}

func Test_sqlite_UpdateOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewOrderRepository(db)

	ctx := context.Background()
	order := domain.Order{
		ID:           1,
		CustomerID:   1,
		RestaurantID: 2,
		OrderItems: []domain.OrderItem{
			{MenuItemID: 1, Quantity: 3},
			{MenuItemID: 2, Quantity: 2},
		},
	}

	// Mock the transaction begin
	mock.ExpectBegin()

	// Mock the update to orders table
	mock.ExpectExec("UPDATE orders").
		WithArgs(order.CustomerID, order.RestaurantID, order.ID).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(nil)

	// Mock the delete from orderitems table
	mock.ExpectExec("DELETE FROM orderitems").
		WithArgs(order.ID).
		WillReturnResult(sqlmock.NewResult(1, 2)).
		WillReturnError(nil)

	// Mock the insert into orderitems table
	for _, item := range order.OrderItems {
		mock.ExpectExec("INSERT INTO orderitems").
			WithArgs(order.ID, item.MenuItemID, item.Quantity).
			WillReturnResult(sqlmock.NewResult(1, 1)).
			WillReturnError(nil)
	}

	// Mock the transaction commit
	mock.ExpectCommit()

	err = repo.UpdateOrder(ctx, order)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
