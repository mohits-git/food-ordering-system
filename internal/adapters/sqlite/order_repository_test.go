package sqlite

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mohits-git/food-ordering-system/internal/domain"
)

func TestOrderRepository_SaveOrder(t *testing.T) {
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
