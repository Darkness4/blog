package db_test

import (
	"os"
	"runtime/debug"
	"testing"

	"github.com/Darkness4/blog/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func TestCreateOrIncrementPageViewsOnUniqueIP(t *testing.T) {
	_ = godotenv.Load(".env")
	_ = godotenv.Load(".env.local")

	ctx := t.Context()

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		t.Log("skip test")
		return
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Log(string(debug.Stack()))
		t.Fatal(err)
	}

	q := db.New(pool)

	// Purge
	_ = q.DeletePageViewsIPs(ctx)
	_ = q.DeletePageViews(ctx)

	// Test
	err = q.CreateOrIncrementPageViewsOnUniqueIP(ctx, pool, "page-id", "ip")
	if err != nil {
		t.Log(string(debug.Stack()))
		t.Fatal(err)
	}

	// Check exists
	pv, err := q.FindPageViews(ctx, "page-id")
	if err != nil {
		t.Log(string(debug.Stack()))
		t.Fatal(err)
	}
	if pv.Views != 1 {
		t.Log(string(debug.Stack()))
		t.Fatalf("expected 1, got %d", pv.Views)
	}

	// Check conflict
	err = q.CreateOrIncrementPageViewsOnUniqueIP(ctx, pool, "page-id", "ip")
	if err != nil {
		t.Log(string(debug.Stack()))
		t.Fatal(err)
	}

	// Check exists
	pv, err = q.FindPageViews(ctx, "page-id")
	if err != nil {
		t.Log(string(debug.Stack()))
		t.Fatal(err)
	}
	if pv.Views != 1 {
		t.Log(string(debug.Stack()))
		t.Fatalf("expected 1, got %d", pv.Views)
	}
}
