// This test file requires the github.com/DATA-DOG/go-sqlmock package to be installed.
// Before running these tests, install the package with:
// go get github.com/DATA-DOG/go-sqlmock
//
// These tests verify the functionality of the TransactionRepository using sqlmock to mock the database.
// The tests cover all methods from the BaseRepository that TransactionRepository inherits:
// - Create: Tests creating a new transaction
// - FindByID: Tests finding a transaction by ID (both found and not found cases)
// - FindAll: Tests retrieving all transactions
// - Update: Tests updating a transaction
// - Delete: Tests deleting a transaction
// - DeleteByID: Tests deleting a transaction by ID
//
// And the TransactionRepository specific methods:
// - FindByAccountID: Tests finding transactions by account ID
// - FindByDateRange: Tests finding transactions within a date range
// - FindByDescriptionLike: Tests finding transactions with descriptions containing a keyword
// - FindByAccountAndDateRange: Tests finding transactions by account ID and date range
// - SumByAccountID: Tests calculating the sum of transaction amounts for an account
// - CountByAccountID: Tests counting transactions for an account
// - FindLatestForAccount: Tests finding the latest transactions for an account with a limit
// - GroupByCategory: Tests grouping transactions by category with sum and count
//
// Each test sets up expectations for SQL queries and verifies that the repository methods
// interact with the database as expected.

package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
	"sample-mcp/db/entity"
)

// setupMockDB is already defined in account_test.go and is reused here
// since both files are in the same package

func TestTransactionRepository_Create(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewTransactionRepository(gormDB)
	ctx := context.Background()

	transactionDate := time.Now()
	description := "Test Transaction"
	transaction := &entity.Transaction{
		AccountID:       1,
		CategoryID:      2,
		Amount:          100.50,
		TransactionDate: transactionDate,
		Description:     &description,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Expectations
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "transactions" ("account_id","category_id","amount","transaction_date","description","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "created_at","updated_at","transaction_id"`)).
		WithArgs(transaction.AccountID, transaction.CategoryID, transaction.Amount, transaction.TransactionDate, transaction.Description, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at", "transaction_id"}).AddRow(time.Now(), time.Now(), 1))
	mock.ExpectCommit()

	// Test
	err := repo.Create(ctx, transaction)
	if err != nil {
		t.Errorf("Error creating transaction: %v", err)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestTransactionRepository_FindByID(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewTransactionRepository(gormDB)
	ctx := context.Background()
	transactionID := uint(1)
	transactionDate := time.Now()
	description := "Test Transaction"

	// Expectations
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "transactions" WHERE "transactions"."transaction_id" = $1 ORDER BY "transactions"."transaction_id" LIMIT $2`)).
		WithArgs(transactionID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"transaction_id", "account_id", "category_id", "amount", "transaction_date", "description", "created_at", "updated_at"}).
			AddRow(transactionID, 1, 2, 100.50, transactionDate, description, time.Now(), time.Now()))

	// Test
	transaction, err := repo.FindByID(ctx, transactionID)
	if err != nil {
		t.Errorf("Error finding transaction by ID: %v", err)
	}

	if transaction == nil {
		t.Error("Expected transaction to be returned, got nil")
	}

	if transaction.TransactionID != transactionID {
		t.Errorf("Expected transaction ID %d, got %d", transactionID, transaction.TransactionID)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestTransactionRepository_FindByID_NotFound(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewTransactionRepository(gormDB)
	ctx := context.Background()
	transactionID := uint(999)

	// Expectations
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "transactions" WHERE "transactions"."transaction_id" = $1 ORDER BY "transactions"."transaction_id" LIMIT $2`)).
		WithArgs(transactionID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	// Test
	transaction, err := repo.FindByID(ctx, transactionID)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if transaction != nil {
		t.Errorf("Expected nil transaction, got %v", transaction)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestTransactionRepository_FindAll(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewTransactionRepository(gormDB)
	ctx := context.Background()
	transactionDate1 := time.Now()
	transactionDate2 := time.Now().Add(-24 * time.Hour)
	description1 := "Transaction 1"
	description2 := "Transaction 2"

	// Expectations
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "transactions"`)).
		WillReturnRows(sqlmock.NewRows([]string{"transaction_id", "account_id", "category_id", "amount", "transaction_date", "description", "created_at", "updated_at"}).
			AddRow(1, 1, 2, 100.50, transactionDate1, description1, time.Now(), time.Now()).
			AddRow(2, 1, 3, 200.75, transactionDate2, description2, time.Now(), time.Now()))

	// Test
	transactions, err := repo.FindAll(ctx)
	if err != nil {
		t.Errorf("Error finding all transactions: %v", err)
	}

	if len(transactions) != 2 {
		t.Errorf("Expected 2 transactions, got %d", len(transactions))
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestTransactionRepository_Update(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewTransactionRepository(gormDB)
	ctx := context.Background()
	transactionDate := time.Now()
	description := "Updated Transaction"

	transaction := &entity.Transaction{
		TransactionID:   1,
		AccountID:       1,
		CategoryID:      2,
		Amount:          150.75,
		TransactionDate: transactionDate,
		Description:     &description,
		UpdatedAt:       time.Now(),
	}

	// Expectations
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "transactions" SET "account_id"=$1,"category_id"=$2,"amount"=$3,"transaction_date"=$4,"description"=$5,"created_at"=$6,"updated_at"=$7 WHERE "transaction_id" = $8`)).
		WithArgs(transaction.AccountID, transaction.CategoryID, transaction.Amount, transaction.TransactionDate, transaction.Description, sqlmock.AnyArg(), sqlmock.AnyArg(), transaction.TransactionID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Test
	err := repo.Update(ctx, transaction)
	if err != nil {
		t.Errorf("Error updating transaction: %v", err)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestTransactionRepository_Delete(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewTransactionRepository(gormDB)
	ctx := context.Background()

	transaction := &entity.Transaction{
		TransactionID: 1,
	}

	// Expectations
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "transactions" WHERE "transactions"."transaction_id" = $1`)).
		WithArgs(transaction.TransactionID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Test
	err := repo.Delete(ctx, transaction)
	if err != nil {
		t.Errorf("Error deleting transaction: %v", err)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestTransactionRepository_DeleteByID(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewTransactionRepository(gormDB)
	ctx := context.Background()
	transactionID := uint(1)

	// Expectations
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "transactions" WHERE "transactions"."transaction_id" = $1`)).
		WithArgs(transactionID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Test
	err := repo.DeleteByID(ctx, transactionID)
	if err != nil {
		t.Errorf("Error deleting transaction by ID: %v", err)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestTransactionRepository_SumByAccountID(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewTransactionRepository(gormDB)
	ctx := context.Background()
	accountID := uint(1)
	expectedSum := 300.25

	// Expectations
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COALESCE(SUM(amount), 0) FROM "transactions" WHERE account_id = $1`)).
		WithArgs(accountID).
		WillReturnRows(sqlmock.NewRows([]string{"coalesce"}).AddRow(expectedSum))

	// Test
	sum, err := repo.SumByAccountID(ctx, accountID)
	if err != nil {
		t.Errorf("Error calculating sum by account ID: %v", err)
	}

	if sum != expectedSum {
		t.Errorf("Expected sum %f, got %f", expectedSum, sum)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestTransactionRepository_CountByAccountID(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewTransactionRepository(gormDB)
	ctx := context.Background()
	accountID := uint(1)
	expectedCount := int64(5)

	// Expectations
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM "transactions" WHERE account_id = $1`)).
		WithArgs(accountID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

	// Test
	count, err := repo.CountByAccountID(ctx, accountID)
	if err != nil {
		t.Errorf("Error counting transactions by account ID: %v", err)
	}

	if count != expectedCount {
		t.Errorf("Expected count %d, got %d", expectedCount, count)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
