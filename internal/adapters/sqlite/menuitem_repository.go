package sqlite

import (
	"context"
	"database/sql"

	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/mohits-git/food-ordering-system/internal/utils/apperr"
)

type MenuItemRepository struct {
	db *sql.DB
}

func NewMenuItemRepository(db *sql.DB) *MenuItemRepository {
	return &MenuItemRepository{db: db}
}

func (m *MenuItemRepository) SaveMenuItem(cxt context.Context, item domain.MenuItem) (int, error) {
	query := `INSERT INTO menuitems (name, price, available, restaurant_id, image_url) VALUES (?, ?, ?, ?, ?) RETURNING id`
	var id int
	err := m.db.QueryRowContext(cxt, query, item.Name, item.Price, item.Available, item.RestaurantID, item.ImageURL).Scan(&id)
	if err != nil {
		return 0, HandleSQLiteError(err)
	}
	return id, nil
}

func (m *MenuItemRepository) UpdateMenuItemAvailability(cxt context.Context, id int, available bool) error {
	query := `UPDATE menuitems SET available = ? WHERE id = ?`
	_, err := m.db.ExecContext(cxt, query, available, id)
	if err != nil {
		return HandleSQLiteError(err)
	}
	return nil
}

func (m *MenuItemRepository) FindMenuItemsByRestaurantId(cxt context.Context, restaurantId int) ([]domain.MenuItem, error) {
	query := `SELECT id, name, price, available, restaurant_id, image_url FROM menuitems WHERE restaurant_id = ?`
	rows, err := m.db.QueryContext(cxt, query, restaurantId)
	if err != nil {
		return nil, HandleSQLiteError(err)
	}
	defer rows.Close()

	menuItems := []domain.MenuItem{}
	for rows.Next() {
		var item domain.MenuItem
		var imageUrl sql.NullString
		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Available, &item.RestaurantID, &imageUrl); err != nil {
			return nil, HandleSQLiteError(err)
		}
		if imageUrl.Valid {
			item.ImageURL = imageUrl.String
		}
		menuItems = append(menuItems, item)
	}

	if err := rows.Err(); err != nil {
		return nil, HandleSQLiteError(err)
	}
	return menuItems, nil
}

func (m *MenuItemRepository) FindMenuItemById(cxt context.Context, id int) (domain.MenuItem, error) {
	query := `SELECT id, name, price, available, restaurant_id, image_url FROM menuitems WHERE id = ?`
	var item domain.MenuItem
	var imageUrl sql.NullString
	err := m.db.QueryRowContext(cxt, query, id).Scan(&item.ID, &item.Name, &item.Price, &item.Available, &item.RestaurantID, &imageUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.MenuItem{}, apperr.NewAppError(apperr.ErrNotFound, "menu item not found", nil)
		}
		return domain.MenuItem{}, HandleSQLiteError(err)
	}
	if imageUrl.Valid {
		item.ImageURL = imageUrl.String
	}
	return item, nil
}
