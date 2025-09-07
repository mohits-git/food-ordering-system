package sqlite

import (
	"context"
	"database/sql"

	"github.com/mohits-git/food-ordering-system/internal/domain"
)

type RestaurantRepository struct {
  db *sql.DB
}

func NewRestaurantRepository(db *sql.DB) *RestaurantRepository {
  return &RestaurantRepository{
    db: db,
  }
}

func (r *RestaurantRepository) SaveRestaurant(ctx context.Context, restaurant domain.Restaurant) (int, error) {
  query := `INSERT INTO restaurants (name, owner_id) VALUES (?, ?) RETURNING id`
  var id int
  err := r.db.QueryRowContext(ctx, query, restaurant.Name, restaurant.OwnerID).Scan(&id)
  if err != nil {
    return 0, HandleSQLiteError(err)
  }
  return id, nil
}

func (r *RestaurantRepository) FindAllRestaurants(ctx context.Context) ([]domain.Restaurant, error) { 
  query := `SELECT id, name, owner_id FROM restaurants`
  rows, err := r.db.QueryContext(ctx, query)
  if err != nil {
    return nil, HandleSQLiteError(err)
  }
  defer rows.Close()

  var restaurants []domain.Restaurant
  for rows.Next() {
    var restaurant domain.Restaurant
    if err := rows.Scan(&restaurant.ID, &restaurant.Name, &restaurant.OwnerID); err != nil {
      return nil, HandleSQLiteError(err)
    }
    restaurants = append(restaurants, restaurant)
  }
  if err := rows.Err(); err != nil {
    return nil, HandleSQLiteError(err)
  }
  return restaurants, nil
}
