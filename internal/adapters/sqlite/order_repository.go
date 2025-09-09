package sqlite

import (
	"context"
	"database/sql"
	"log"

	"github.com/mohits-git/food-ordering-system/internal/domain"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (o *OrderRepository) SaveOrder(ctx context.Context, order domain.Order) (int, error) {
	tx, err := o.db.BeginTx(ctx, nil)

	query := "INSERT INTO orders (user_id, restaurant_id) VALUES (?, ?) RETURNING id"
	var id int
	err = tx.QueryRowContext(ctx, query, order.CustomerID, order.RestaurantID).Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, HandleSQLiteError(err)
	}

	// save order items
	itemQuery := "INSERT INTO orderitems (order_id, menuitem_id, quantity) VALUES (?, ?, ?)"
	for _, item := range order.OrderItems {
		_, err := tx.ExecContext(ctx, itemQuery, id, item.MenuItemID, item.Quantity)
		if err != nil {
			tx.Rollback()
			return 0, HandleSQLiteError(err)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, HandleSQLiteError(err)
	}
	return id, nil
}

func (o *OrderRepository) FindOrderById(ctx context.Context, id int) (domain.Order, error) {
	var order domain.Order
	query := "SELECT id, user_id, restaurant_id FROM orders WHERE id = ?"
	err := o.db.QueryRowContext(ctx, query, id).Scan(&order.ID, &order.CustomerID, &order.RestaurantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Order{}, nil // or return a custom NotFound error
		}
		return domain.Order{}, HandleSQLiteError(err)
	}

	// fetch order items
	itemQuery := "SELECT menuitem_id, quantity FROM orderitems WHERE order_id = ?"
	rows, err := o.db.QueryContext(ctx, itemQuery, id)
	if err != nil {
		return domain.Order{}, HandleSQLiteError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.OrderItem
		if err := rows.Scan(&item.MenuItemID, &item.Quantity); err != nil {
			return domain.Order{}, HandleSQLiteError(err)
		}
		order.OrderItems = append(order.OrderItems, item)
	}
	if err := rows.Err(); err != nil {
		return domain.Order{}, HandleSQLiteError(err)
	}

	return order, nil
}

func (o *OrderRepository) UpdateOrder(ctx context.Context, order domain.Order) error {
	tx, err := o.db.BeginTx(ctx, nil)
	defer func() {
		if r := recover(); r != nil {
			if err := tx.Rollback(); err != nil {
				log.Println("update order tx rollback error: ", err)
			}
		}
	}()

	query := "UPDATE orders SET user_id = ?, restaurant_id = ? WHERE id = ?"
	_, err = tx.ExecContext(ctx, query, order.CustomerID, order.RestaurantID, order.ID)
	if err != nil {
		tx.Rollback()
		return HandleSQLiteError(err)
	}

	// For simplicity, delete existing items and re-insert
	delQuery := "DELETE FROM orderitems WHERE order_id = ?"
	_, err = tx.ExecContext(ctx, delQuery, order.ID)
	if err != nil {
		tx.Rollback()
		return HandleSQLiteError(err)
	}

	itemQuery := "INSERT INTO orderitems (order_id, menuitem_id, quantity) VALUES (?, ?, ?)"
	for _, item := range order.OrderItems {
		_, err := tx.ExecContext(ctx, itemQuery, order.ID, item.MenuItemID, item.Quantity)
		if err != nil {
			tx.Rollback()
			return HandleSQLiteError(err)
		}
	}

	if err := tx.Commit(); err != nil {
		return HandleSQLiteError(err)
	}
	return nil
}
