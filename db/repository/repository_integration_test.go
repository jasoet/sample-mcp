//go:build integration

package repository_test

import (
	"gorm.io/gorm"
	"sample-mcp/db"
	pkgDB "sample-mcp/pkg/db"
	"testing"
	"time"
)

var TestDB *gorm.DB

func TestMain(m *testing.M) {
	cfg := pkgDB.ConnectionConfig{
		DbType:       pkgDB.Postgresql,
		Host:         "localhost",
		Port:         5432,
		Username:     "jasoet",
		Password:     "localhost",
		DbName:       "mcp_db",
		Timeout:      10 * time.Second,
		MaxIdleConns: 5,
		MaxOpenConns: 10,
	}

	pool, err := cfg.Pool()
	if err != nil {
		panic(err)
	}

	sqlDb, err := pool.DB()
	err = db.RunMigrations(sqlDb)
	if err != nil {
		panic(err)
	}

	TestDB = pool
	m.Run()
}
