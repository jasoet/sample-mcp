//go:build integration

package db

import (
	"database/sql"
	"testing"
	"time"

	pkgdb "sample-mcp/pkg/db"
)

func TestRunMigrations(t *testing.T) {
	config := &pkgdb.ConnectionConfig{
		DbType:       pkgdb.Postgresql,
		Host:         "localhost",
		Port:         5432,
		Username:     "jasoet",
		Password:     "localhost",
		DbName:       "mcp_db",
		Timeout:      10 * time.Second,
		MaxIdleConns: 5,
		MaxOpenConns: 10,
	}

	db, err := config.SqlDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = RunMigrations(db)
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	tablesExist, err := verifyTablesExist(db)
	if err != nil {
		t.Fatalf("Failed to verify tables: %v", err)
	}

	if !tablesExist {
		t.Errorf("Expected tables to exist after migrations, but they don't")
	}
}

func verifyTablesExist(db *sql.DB) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'accounts'
		) AND EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'categories'
		) AND EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'transactions'
		);
	`

	var exists bool
	err := db.QueryRow(query).Scan(&exists)
	return exists, err
}
