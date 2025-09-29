package sqlite

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mohits-git/food-ordering-system/internal/domain"
	"github.com/mohits-git/food-ordering-system/internal/utils/apperr"
	"github.com/stretchr/testify/require"
)

func Test_sqlite_MenuItemRepository_NewMenuItemRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Expected no error when creating sqlmock")
	defer db.Close()
	_ = mock

	repo := NewMenuItemRepository(db)
	require.NotNil(t, repo, "Expected NewMenuItemRepository to return a non-nil repository")
}

func Test_sqlite_MenuItemRepository_SaveMenuItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Expected no error when creating sqlmock")
	defer db.Close()

	repo := NewMenuItemRepository(db)
	require.NotNil(t, repo, "Expected NewMenuItemRepository to return a non-nil repository")

	// Define test cases
	tests := []struct {
		name          string
		menuItem      domain.MenuItem
		mockSetup     func()
		expectedID    int
		expectedError bool
	}{
		{
			name: "Successful insert",
			menuItem: domain.MenuItem{
				Name:         "Test Item",
				Price:        9.99,
				Available:    true,
				RestaurantID: 1,
				ImageURL:     "file.com",
			},
			mockSetup: func() {
				mock.ExpectQuery("INSERT INTO menuitems").
					WithArgs("Test Item", 9.99, true, 1, "file.com").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			expectedID:    1,
			expectedError: false,
		},
		{
			name: "Database error",
			menuItem: domain.MenuItem{
				Name:         "Test Item",
				Price:        9.99,
				Available:    true,
				RestaurantID: 1,
				ImageURL:     "file.com",
			},
			mockSetup: func() {
				mock.ExpectQuery("INSERT INTO menuitems").
					WithArgs("Test Item", 9.99, true, 1, "file.com").
					WillReturnError(sqlmock.ErrCancelled)
			},
			expectedID:    0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			id, err := repo.SaveMenuItem(t.Context(), tt.menuItem)
			if tt.expectedError {
				require.Error(t, err)
				require.Equal(t, tt.expectedID, id)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedID, id)
			}

			err = mock.ExpectationsWereMet()
			require.NoError(t, err, "There were unfulfilled expectations")
		})
	}
}

func Test_sqlite_MenuItemRepository_UpdateMenuItemAvailability(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Expected no error when creating sqlmock")
	defer db.Close()

	repo := NewMenuItemRepository(db)
	require.NotNil(t, repo, "Expected NewMenuItemRepository to return a non-nil repository")

	// Define test cases
	tests := []struct {
		name          string
		menuItemID    int
		available     bool
		mockSetup     func()
		expectedError bool
	}{
		{
			name:       "Successful update",
			menuItemID: 1,
			available:  true,
			mockSetup: func() {
				mock.ExpectExec("UPDATE menuitems SET available").
					WithArgs(true, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: false,
		},
		{
			name:       "Database error",
			menuItemID: 1,
			available:  true,
			mockSetup: func() {
				mock.ExpectExec("UPDATE menuitems SET available").
					WithArgs(true, 1).
					WillReturnError(sqlmock.ErrCancelled)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			err := repo.UpdateMenuItemAvailability(t.Context(), tt.menuItemID, tt.available)
			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			err = mock.ExpectationsWereMet()
			require.NoError(t, err, "There were unfulfilled expectations")
		})
	}
}

func Test_sqlite_MenuItemRepository_FindMenuItemsByRestaurantId(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Expected no error when creating sqlmock")
	defer db.Close()

	repo := NewMenuItemRepository(db)
	require.NotNil(t, repo, "Expected NewMenuItemRepository to return a non-nil repository")

	// Define test cases
	tests := []struct {
		name            string
		restaurantID    int
		mockSetup       func()
		expectedResults []domain.MenuItem
		expectedError   bool
	}{
		{
			name:         "Successful fetch",
			restaurantID: 1,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "price", "available", "restaurant_id", "image_url"}).
					AddRow(1, "Item 1", 9.99, true, 1, "file.com").
					AddRow(2, "Item 2", 19.99, false, 1, "file.com")
				mock.ExpectQuery("SELECT id, name, price, available, restaurant_id, image_url FROM menuitems WHERE restaurant_id = \\?").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedResults: []domain.MenuItem{
				{ID: 1, Name: "Item 1", Price: 9.99, Available: true, RestaurantID: 1, ImageURL: "file.com"},
				{ID: 2, Name: "Item 2", Price: 19.99, Available: false, RestaurantID: 1, ImageURL: "file.com"},
			},
			expectedError: false,
		},
		{
			name:         "No items found",
			restaurantID: 2,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "price", "available", "restaurant_id", "image_url"})
				mock.ExpectQuery("SELECT id, name, price, available, restaurant_id, image_url FROM menuitems WHERE restaurant_id = \\?").
					WithArgs(2).
					WillReturnRows(rows)
			},
			expectedResults: []domain.MenuItem{},
			expectedError:   false,
		},
		{
			name:         "Database error",
			restaurantID: 1,
			mockSetup: func() {
				mock.ExpectQuery("SELECT id, name, price, available, restaurant_id, image_url FROM menuitems WHERE restaurant_id = \\?").
					WithArgs(1).
					WillReturnError(sqlmock.ErrCancelled)
			},
			expectedResults: nil,
			expectedError:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			results, err := repo.FindMenuItemsByRestaurantId(t.Context(), tt.restaurantID)
			if tt.expectedError {
				require.Error(t, err)
				require.Equal(t, tt.expectedResults, results)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedResults, results)
			}

			err = mock.ExpectationsWereMet()
			require.NoError(t, err, "There were unfulfilled expectations")
		})
	}
}

func Test_sqlite_MenuItemRepository_FindMenuItemById(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Expected no error when creating sqlmock")
	defer db.Close()

	repo := NewMenuItemRepository(db)
	require.NotNil(t, repo, "Expected NewMenuItemRepository to return a non-nil repository")

	// Define test cases
	tests := []struct {
		name            string
		menuItemID      int
		mockSetup       func()
		expectedResult  domain.MenuItem
		expectedError   bool
		expectedErrCode apperr.AppErrorCode
	}{
		{
			name:       "Successful fetch",
			menuItemID: 1,
			mockSetup: func() {
				row := sqlmock.NewRows([]string{"id", "name", "price", "available", "restaurant_id", "image_url"}).
					AddRow(1, "Item 1", 9.99, true, 1, "file.com")
				mock.ExpectQuery("SELECT id, name, price, available, restaurant_id, image_url FROM menuitems WHERE id = \\?").
					WithArgs(1).
					WillReturnRows(row)
			},
			expectedResult: domain.MenuItem{ID: 1, Name: "Item 1", Price: 9.99, Available: true, RestaurantID: 1, ImageURL: "file.com"},
			expectedError:  false,
		},
		{
			name:       "Menu item not found",
			menuItemID: 2,
			mockSetup: func() {
				mock.ExpectQuery("SELECT id, name, price, available, restaurant_id, image_url FROM menuitems WHERE id = \\?").
					WithArgs(2).
					WillReturnError(sql.ErrNoRows)
			},
			expectedResult:  domain.MenuItem{},
			expectedError:   true,
			expectedErrCode: apperr.ErrNotFound,
		},
		{
			name:       "Database error",
			menuItemID: 3,
			mockSetup: func() {
				mock.ExpectQuery("SELECT id, name, price, available, restaurant_id, image_url FROM menuitems WHERE id = \\?").
					WithArgs(3).
					WillReturnError(sqlmock.ErrCancelled)
			},
			expectedResult:  domain.MenuItem{},
			expectedError:   true,
			expectedErrCode: apperr.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			result, err := repo.FindMenuItemById(t.Context(), tt.menuItemID)
			if tt.expectedError {
				require.Error(t, err)
				require.Equal(t, tt.expectedResult, result)
				if appErr, ok := err.(*apperr.AppError); ok {
					require.Equal(t, tt.expectedErrCode, appErr.Code)
				} else {
					t.Errorf("Expected error to be of type AppError")
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, result)
			}

			err = mock.ExpectationsWereMet()
			require.NoError(t, err, "There were unfulfilled expectations")
		})
	}
}
