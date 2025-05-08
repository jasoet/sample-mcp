package repository

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sample-mcp/db/entity"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *gorm.DB, func()) {
	// Create a new SQL mock
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}

	// Create a new GORM DB instance with the mock
	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open gorm connection: %v", err)
	}

	// Return cleanup function
	cleanup := func() {
		mockDB.Close()
	}

	return mockDB, mock, gormDB, cleanup
}

func TestAccountRepository_Create(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewAccountRepository(gormDB)
	ctx := context.Background()

	account := &entity.Account{
		Name:        "Test Account",
		AccountType: "Savings",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Expectations
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "accounts" ("name","account_type","created_at","updated_at") VALUES ($1,$2,$3,$4) RETURNING "created_at","updated_at","account_id"`)).
		WithArgs(account.Name, account.AccountType, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at", "account_id"}).AddRow(time.Now(), time.Now(), 1))
	mock.ExpectCommit()

	// Test
	err := repo.Create(ctx, account)
	if err != nil {
		t.Errorf("Error creating account: %v", err)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestAccountRepository_FindByID(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewAccountRepository(gormDB)
	ctx := context.Background()
	accountID := uint(1)

	// Expectations
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "accounts" WHERE "accounts"."account_id" = $1 ORDER BY "accounts"."account_id" LIMIT $2`)).
		WithArgs(accountID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"account_id", "name", "account_type", "created_at", "updated_at"}).
			AddRow(accountID, "Test Account", "Savings", time.Now(), time.Now()))

	// Test
	account, err := repo.FindByID(ctx, accountID)
	if err != nil {
		t.Errorf("Error finding account by ID: %v", err)
	}

	if account == nil {
		t.Error("Expected account to be returned, got nil")
	}

	if account.AccountID != accountID {
		t.Errorf("Expected account ID %d, got %d", accountID, account.AccountID)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestAccountRepository_FindByID_NotFound(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewAccountRepository(gormDB)
	ctx := context.Background()
	accountID := uint(999)

	// Expectations
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "accounts" WHERE "accounts"."account_id" = $1 ORDER BY "accounts"."account_id" LIMIT $2`)).
		WithArgs(accountID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	// Test
	account, err := repo.FindByID(ctx, accountID)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if account != nil {
		t.Errorf("Expected nil account, got %v", account)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestAccountRepository_FindAll(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewAccountRepository(gormDB)
	ctx := context.Background()

	// Expectations
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "accounts"`)).
		WillReturnRows(sqlmock.NewRows([]string{"account_id", "name", "account_type", "created_at", "updated_at"}).
			AddRow(1, "Account 1", "Savings", time.Now(), time.Now()).
			AddRow(2, "Account 2", "Checking", time.Now(), time.Now()))

	// Test
	accounts, err := repo.FindAll(ctx)
	if err != nil {
		t.Errorf("Error finding all accounts: %v", err)
	}

	if len(accounts) != 2 {
		t.Errorf("Expected 2 accounts, got %d", len(accounts))
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestAccountRepository_Update(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewAccountRepository(gormDB)
	ctx := context.Background()

	account := &entity.Account{
		AccountID:   1,
		Name:        "Updated Account",
		AccountType: "Checking",
		UpdatedAt:   time.Now(),
	}

	// Expectations
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "accounts" SET "name"=$1,"account_type"=$2,"created_at"=$3,"updated_at"=$4 WHERE "account_id" = $5`)).
		WithArgs(account.Name, account.AccountType, sqlmock.AnyArg(), sqlmock.AnyArg(), account.AccountID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Test
	err := repo.Update(ctx, account)
	if err != nil {
		t.Errorf("Error updating account: %v", err)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestAccountRepository_Delete(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewAccountRepository(gormDB)
	ctx := context.Background()

	account := &entity.Account{
		AccountID: 1,
	}

	// Expectations
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "accounts" WHERE "accounts"."account_id" = $1`)).
		WithArgs(account.AccountID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Test
	err := repo.Delete(ctx, account)
	if err != nil {
		t.Errorf("Error deleting account: %v", err)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestAccountRepository_DeleteByID(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewAccountRepository(gormDB)
	ctx := context.Background()
	accountID := uint(1)

	// Expectations
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "accounts" WHERE "accounts"."account_id" = $1`)).
		WithArgs(accountID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Test
	err := repo.DeleteByID(ctx, accountID)
	if err != nil {
		t.Errorf("Error deleting account by ID: %v", err)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestAccountRepository_FindByName(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewAccountRepository(gormDB)
	ctx := context.Background()
	name := "Test Account"

	// Expectations
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "accounts" WHERE name = $1 ORDER BY "accounts"."account_id" LIMIT $2`)).
		WithArgs(name, 1).
		WillReturnRows(sqlmock.NewRows([]string{"account_id", "name", "account_type", "created_at", "updated_at"}).
			AddRow(1, name, "Savings", time.Now(), time.Now()))

	// Test
	account, err := repo.FindByName(ctx, name)
	if err != nil {
		t.Errorf("Error finding account by name: %v", err)
	}

	if account == nil {
		t.Error("Expected account to be returned, got nil")
	}

	if account.Name != name {
		t.Errorf("Expected account name %s, got %s", name, account.Name)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestAccountRepository_FindByName_NotFound(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewAccountRepository(gormDB)
	ctx := context.Background()
	name := "Nonexistent Account"

	// Expectations
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "accounts" WHERE name = $1 ORDER BY "accounts"."account_id" LIMIT $2`)).
		WithArgs(name, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	// Test
	account, err := repo.FindByName(ctx, name)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if account != nil {
		t.Errorf("Expected nil account, got %v", account)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestAccountRepository_FindByNameLike(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewAccountRepository(gormDB)
	ctx := context.Background()
	keyword := "Test"

	// Expectations
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "accounts" WHERE name ILIKE $1`)).
		WithArgs("%" + keyword + "%").
		WillReturnRows(sqlmock.NewRows([]string{"account_id", "name", "account_type", "created_at", "updated_at"}).
			AddRow(1, "Test Account 1", "Savings", time.Now(), time.Now()).
			AddRow(2, "Test Account 2", "Checking", time.Now(), time.Now()))

	// Test
	accounts, err := repo.FindByNameLike(ctx, keyword)
	if err != nil {
		t.Errorf("Error finding accounts by name like: %v", err)
	}

	if len(accounts) != 2 {
		t.Errorf("Expected 2 accounts, got %d", len(accounts))
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestAccountRepository_FindByNameLike_NotFound(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewAccountRepository(gormDB)
	ctx := context.Background()
	keyword := "Nonexistent"

	// Expectations
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "accounts" WHERE name ILIKE $1`)).
		WithArgs("%" + keyword + "%").
		WillReturnRows(sqlmock.NewRows([]string{"account_id", "name", "account_type", "created_at", "updated_at"}))

	// Test
	accounts, err := repo.FindByNameLike(ctx, keyword)
	if err != nil {
		t.Errorf("Error finding accounts by name like: %v", err)
	}

	if len(accounts) != 0 {
		t.Errorf("Expected 0 accounts, got %d", len(accounts))
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestAccountRepository_FindByNameLike_Error(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewAccountRepository(gormDB)
	ctx := context.Background()
	keyword := "Error"

	// Expectations
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "accounts" WHERE name ILIKE $1`)).
		WithArgs("%" + keyword + "%").
		WillReturnError(errors.New("database error"))

	// Test
	accounts, err := repo.FindByNameLike(ctx, keyword)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if accounts != nil {
		t.Errorf("Expected nil accounts, got %v", accounts)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
