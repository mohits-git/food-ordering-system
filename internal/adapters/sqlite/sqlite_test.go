package sqlite

import (
	"context"
	"database/sql"
	"testing"
)

func TestConnect(t *testing.T) {
	ctx := context.Background()
	dsn := ":memory:"

	db, err := Connect(ctx, dsn)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	if db == nil {
		t.Fatal("expected db to be non-nil")
	}

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}
}

func TestConnect_Concurrent(t *testing.T) {
  ctx := context.Background()
  dsn := ":memory:"

  const goroutines = 10
  errCh := make(chan error, goroutines)
  dbCh := make(chan *sql.DB, goroutines)

  singletonDB, err := Connect(ctx, dsn)
  if err != nil {
    t.Fatalf("failed to connect to database: %v", err)
  }
  if singletonDB == nil {
    t.Fatal("expected Connect() to return non-nil *sql.DB")
  }


  for range goroutines {
    go func() {
      db, err := Connect(ctx, dsn)
      if err != nil {
        errCh <- err
        return
      }
      dbCh <- db
    }()
  }

  for range goroutines {
    select {
    case err := <-errCh:
      t.Fatalf("failed to connect to database: %v", err)
    case db := <-dbCh:
      if db == nil {
        t.Fatal("expected db to be non-nil")
      }
      if db != singletonDB {
        t.Fatal("expected all calls to Connect() to return the same *sql.DB instance")
      }
      if err := db.PingContext(ctx); err != nil {
        t.Fatalf("failed to ping database: %v", err)
      }
    }
  }
}
