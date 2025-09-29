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
	query := `INSERT INTO restaurants (name, owner_id, image_url) VALUES (?, ?, ?) RETURNING id`
	var id int
	err := r.db.QueryRowContext(ctx, query, restaurant.Name, restaurant.OwnerID, restaurant.ImageURL).Scan(&id)
	if err != nil {
		return 0, HandleSQLiteError(err)
	}
	return id, nil
}

func (r *RestaurantRepository) FindAllRestaurants(ctx context.Context) ([]domain.Restaurant, error) {
	query := `SELECT id, name, owner_id, image_url FROM restaurants`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, HandleSQLiteError(err)
	}
	defer rows.Close()

	var restaurants []domain.Restaurant
	for rows.Next() {
		var restaurant domain.Restaurant
		var imageUrl sql.NullString
		if err := rows.Scan(&restaurant.ID, &restaurant.Name, &restaurant.OwnerID, imageUrl); err != nil {
			return nil, HandleSQLiteError(err)
		}
		if imageUrl.Valid {
			restaurant.ImageURL = imageUrl.String
		}
		restaurants = append(restaurants, restaurant)
	}
	if err := rows.Err(); err != nil {
		return nil, HandleSQLiteError(err)
	}
	return restaurants, nil
}

func (r *RestaurantRepository) FindRestaurantById(ctx context.Context, id int) (domain.Restaurant, error) {
	query := `SELECT id, name, owner_id, image_url FROM restaurants WHERE id = ?`
	var restaurant domain.Restaurant
	var imageUrl sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(&restaurant.ID, &restaurant.Name, &restaurant.OwnerID, imageUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Restaurant{}, nil
		}
		return domain.Restaurant{}, HandleSQLiteError(err)
	}
	if imageUrl.Valid {
		restaurant.ImageURL = imageUrl.String
	}
	return restaurant, nil
}
