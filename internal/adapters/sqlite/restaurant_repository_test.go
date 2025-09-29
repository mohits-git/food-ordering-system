package sqlite

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mattn/go-sqlite3"
	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_sqlite_RestaurantRepository_NewRestaurantRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Expected no error when creating sqlmock")
	defer db.Close()
	_ = mock

	repo := NewRestaurantRepository(db)
	require.NotNil(t, repo, "Expected NewRestaurantRepository to return a non-nil repository")
}

func Test_sqlite_RestaurantRepository_SaveRestaurant(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Expected no error when creating sqlmock")
	defer db.Close()

	repo := NewRestaurantRepository(db)
	require.NotNil(t, repo, "Expected NewRestaurantRepository to return a non-nil repository")

	// Define test cases
	tests := []struct {
		name             string
		restaurant       domain.Restaurant
		mockSetup        func()
		expectedID       int
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			name: "Successful insert",
			restaurant: domain.Restaurant{
				Name:     "Test Restaurant",
				OwnerID:  1,
				ImageURL: "file.com",
			},
			mockSetup: func() {
				mock.ExpectQuery("INSERT INTO restaurants").
					WithArgs("Test Restaurant", 1, "file.com").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			expectedID:    1,
			expectedError: false,
		},
		{
			name: "Database error",
			restaurant: domain.Restaurant{
				Name:     "Test Restaurant",
				OwnerID:  1,
				ImageURL: "file.com",
			},
			mockSetup: func() {
				mock.ExpectQuery("INSERT INTO restaurants").
					WithArgs("Test Restaurant", 1, "file.com").
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

			id, err := repo.SaveRestaurant(t.Context(), tt.restaurant)
			if tt.expectedError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrorMsg)
				assert.Equal(t, tt.expectedID, id)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			}

			err = mock.ExpectationsWereMet()
			require.NoError(t, err, "There were unfulfilled expectations")
		})
	}
}

func Test_sqlite_RestaurantRepository_FindAllRestaurants(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Expected no error when creating sqlmock")
	defer db.Close()

	repo := NewRestaurantRepository(db)
	require.NotNil(t, repo, "Expected NewRestaurantRepository to return a non-nil repository")

	// Define test cases
	tests := []struct {
		name             string
		mockSetup        func()
		expectedResults  []domain.Restaurant
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			name: "Successful fetch",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "owner_id", "image_url"}).
					AddRow(1, "Restaurant 1", 1, "file.com").
					AddRow(2, "Restaurant 2", 2, "file.com")
				mock.ExpectQuery("SELECT id, name, owner_id, image_url FROM restaurants").WillReturnRows(rows)
			},
			expectedResults: []domain.Restaurant{
				{ID: 1, Name: "Restaurant 1", OwnerID: 1, ImageURL: "file.com"},
				{ID: 2, Name: "Restaurant 2", OwnerID: 2, ImageURL: "file.com"},
			},
			expectedError: false,
		},
		{
			name: "Database error",
			mockSetup: func() {
				mock.ExpectQuery("SELECT id, name, owner_id, image_url FROM restaurants").
					WillReturnError(sql.ErrConnDone)
			},
			expectedResults:  nil,
			expectedError:    true,
			expectedErrorMsg: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			results, err := repo.FindAllRestaurants(t.Context())
			if tt.expectedError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrorMsg)
				assert.Nil(t, results)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedResults, results)
			}

			err = mock.ExpectationsWereMet()
			require.NoError(t, err, "There were unfulfilled expectations")
		})
	}
}

func Test_sqlite_RestaurantRepository_FindRestaurantById(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Expected no error when creating sqlmock")
	defer db.Close()

	repo := NewRestaurantRepository(db)
	require.NotNil(t, repo, "Expected NewRestaurantRepository to return a non-nil repository")

	// Define test cases
	tests := []struct {
		name             string
		restaurantID     int
		mockSetup        func()
		expectedResult   domain.Restaurant
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			name:         "Successful fetch",
			restaurantID: 1,
			mockSetup: func() {
				row := sqlmock.NewRows([]string{"id", "name", "owner_id", "image_url"}).
					AddRow(1, "Restaurant 1", 1, "file.com")
				mock.ExpectQuery("SELECT id, name, owner_id, image_url FROM restaurants WHERE id = ?").
					WithArgs(1).
					WillReturnRows(row)
			},
			expectedResult: domain.Restaurant{ID: 1, Name: "Restaurant 1", OwnerID: 1, ImageURL: "file.com"},
			expectedError:  false,
		},
		{
			name:         "Restaurant not found",
			restaurantID: 2,
			mockSetup: func() {
				mock.ExpectQuery("SELECT id, name, owner_id, image_url FROM restaurants WHERE id = ?").
					WithArgs(2).
					WillReturnError(sql.ErrNoRows)
			},
			expectedResult: domain.Restaurant{},
			expectedError:  false,
		},
		{
			name:         "Database error",
			restaurantID: 3,
			mockSetup: func() {
				mock.ExpectQuery("SELECT id, name, owner_id, image_url FROM restaurants WHERE id = ?").
					WithArgs(3).
					WillReturnError(sql.ErrConnDone)
			},
			expectedResult:   domain.Restaurant{},
			expectedError:    true,
			expectedErrorMsg: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			result, err := repo.FindRestaurantById(t.Context(), tt.restaurantID)
			if tt.expectedError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrorMsg)
				assert.Equal(t, tt.expectedResult, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
			err = mock.ExpectationsWereMet()
			require.NoError(t, err, "There were unfulfilled expectations")
		})
	}
}
