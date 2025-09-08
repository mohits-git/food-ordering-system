package sqlite

import (
	"context"
	"database/sql"

	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/mohits-git/food-ordering-system/internal/ports"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) ports.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindUserById(ctx context.Context, id int) (domain.User, error) {
	var user domain.User
	query := "SELECT id, name, email, role, password FROM users WHERE id = ?"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.Password)
	if err != nil {
		err = HandleSQLiteError(err)
		return domain.User{}, err
	}
	return user, nil
}

func (r *UserRepository) FindUserByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	query := "SELECT id, name, email, role, password FROM users WHERE email = ?"
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.Password)
	if err != nil {
		return domain.User{}, HandleSQLiteError(err)
	}
	return user, nil
}

func (r *UserRepository) SaveUser(ctx context.Context, user domain.User) (int, error) {
	query := "INSERT INTO users (name, email, role, password) VALUES (?, ?, ?, ?)"
	res, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.Role, user.Password)
	if err != nil {
		return 0, HandleSQLiteError(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, HandleSQLiteError(err)
	}
	return int(id), nil
}
