package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
	"sample-mcp/db/entity"
)

func TestCategoryRepository_Create(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewCategoryRepository(gormDB)
	ctx := context.Background()

	category := &entity.Category{
		Name:         "Test Category",
		CategoryType: "Expense",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Expectations
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "categories" ("name","category_type","created_at","updated_at") VALUES ($1,$2,$3,$4) RETURNING "created_at","updated_at","category_id"`)).
		WithArgs(category.Name, category.CategoryType, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at", "category_id"}).AddRow(time.Now(), time.Now(), 1))
	mock.ExpectCommit()

	// Test
	err := repo.Create(ctx, category)
	if err != nil {
		t.Errorf("Error creating category: %v", err)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCategoryRepository_FindByID(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewCategoryRepository(gormDB)
	ctx := context.Background()
	categoryID := uint(1)

	// Expectations
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories" WHERE "categories"."category_id" = $1 ORDER BY "categories"."category_id" LIMIT $2`)).
		WithArgs(categoryID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"category_id", "name", "category_type", "created_at", "updated_at"}).
			AddRow(categoryID, "Test Category", "Expense", time.Now(), time.Now()))

	// Test
	category, err := repo.FindByID(ctx, categoryID)
	if err != nil {
		t.Errorf("Error finding category by ID: %v", err)
	}

	if category == nil {
		t.Error("Expected category to be returned, got nil")
	}

	if category.CategoryID != categoryID {
		t.Errorf("Expected category ID %d, got %d", categoryID, category.CategoryID)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCategoryRepository_FindByID_NotFound(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewCategoryRepository(gormDB)
	ctx := context.Background()
	categoryID := uint(999)

	// Expectations
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories" WHERE "categories"."category_id" = $1 ORDER BY "categories"."category_id" LIMIT $2`)).
		WithArgs(categoryID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	// Test
	category, err := repo.FindByID(ctx, categoryID)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if category != nil {
		t.Errorf("Expected nil category, got %v", category)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCategoryRepository_FindAll(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewCategoryRepository(gormDB)
	ctx := context.Background()

	// Expectations
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories"`)).
		WillReturnRows(sqlmock.NewRows([]string{"category_id", "name", "category_type", "created_at", "updated_at"}).
			AddRow(1, "Category 1", "Expense", time.Now(), time.Now()).
			AddRow(2, "Category 2", "Income", time.Now(), time.Now()))

	// Test
	categories, err := repo.FindAll(ctx)
	if err != nil {
		t.Errorf("Error finding all categories: %v", err)
	}

	if len(categories) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(categories))
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCategoryRepository_Update(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewCategoryRepository(gormDB)
	ctx := context.Background()

	category := &entity.Category{
		CategoryID:   1,
		Name:         "Updated Category",
		CategoryType: "Income",
		UpdatedAt:    time.Now(),
	}

	// Expectations
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "categories" SET "name"=$1,"category_type"=$2,"created_at"=$3,"updated_at"=$4 WHERE "category_id" = $5`)).
		WithArgs(category.Name, category.CategoryType, sqlmock.AnyArg(), sqlmock.AnyArg(), category.CategoryID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Test
	err := repo.Update(ctx, category)
	if err != nil {
		t.Errorf("Error updating category: %v", err)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCategoryRepository_Delete(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewCategoryRepository(gormDB)
	ctx := context.Background()

	category := &entity.Category{
		CategoryID: 1,
	}

	// Expectations
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "categories" WHERE "categories"."category_id" = $1`)).
		WithArgs(category.CategoryID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Test
	err := repo.Delete(ctx, category)
	if err != nil {
		t.Errorf("Error deleting category: %v", err)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCategoryRepository_DeleteByID(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewCategoryRepository(gormDB)
	ctx := context.Background()
	categoryID := uint(1)

	// Expectations
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "categories" WHERE "categories"."category_id" = $1`)).
		WithArgs(categoryID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Test
	err := repo.DeleteByID(ctx, categoryID)
	if err != nil {
		t.Errorf("Error deleting category by ID: %v", err)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCategoryRepository_FindByType(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewCategoryRepository(gormDB)
	ctx := context.Background()
	categoryType := "Expense"

	// Expectations
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories" WHERE category_type = $1`)).
		WithArgs(categoryType).
		WillReturnRows(sqlmock.NewRows([]string{"category_id", "name", "category_type", "created_at", "updated_at"}).
			AddRow(1, "Food", categoryType, time.Now(), time.Now()).
			AddRow(2, "Transportation", categoryType, time.Now(), time.Now()))

	// Test
	categories, err := repo.FindByType(ctx, categoryType)
	if err != nil {
		t.Errorf("Error finding categories by type: %v", err)
	}

	if len(categories) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(categories))
	}

	for _, category := range categories {
		if category.CategoryType != categoryType {
			t.Errorf("Expected category type %s, got %s", categoryType, category.CategoryType)
		}
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCategoryRepository_FindByType_NotFound(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewCategoryRepository(gormDB)
	ctx := context.Background()
	categoryType := "NonexistentType"

	// Expectations
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories" WHERE category_type = $1`)).
		WithArgs(categoryType).
		WillReturnRows(sqlmock.NewRows([]string{"category_id", "name", "category_type", "created_at", "updated_at"}))

	// Test
	categories, err := repo.FindByType(ctx, categoryType)
	if err != nil {
		t.Errorf("Error finding categories by type: %v", err)
	}

	if len(categories) != 0 {
		t.Errorf("Expected 0 categories, got %d", len(categories))
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCategoryRepository_FindByNameLike(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewCategoryRepository(gormDB)
	ctx := context.Background()
	keyword := "Food"

	// Expectations
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories" WHERE name ILIKE $1`)).
		WithArgs("%" + keyword + "%").
		WillReturnRows(sqlmock.NewRows([]string{"category_id", "name", "category_type", "created_at", "updated_at"}).
			AddRow(1, "Food", "Expense", time.Now(), time.Now()).
			AddRow(2, "Fast Food", "Expense", time.Now(), time.Now()))

	// Test
	categories, err := repo.FindByNameLike(ctx, keyword)
	if err != nil {
		t.Errorf("Error finding categories by name like: %v", err)
	}

	if len(categories) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(categories))
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCategoryRepository_FindByNameLike_NotFound(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewCategoryRepository(gormDB)
	ctx := context.Background()
	keyword := "Nonexistent"

	// Expectations
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories" WHERE name ILIKE $1`)).
		WithArgs("%" + keyword + "%").
		WillReturnRows(sqlmock.NewRows([]string{"category_id", "name", "category_type", "created_at", "updated_at"}))

	// Test
	categories, err := repo.FindByNameLike(ctx, keyword)
	if err != nil {
		t.Errorf("Error finding categories by name like: %v", err)
	}

	if len(categories) != 0 {
		t.Errorf("Expected 0 categories, got %d", len(categories))
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCategoryRepository_FindByNameLike_Error(t *testing.T) {
	// Setup
	_, mock, gormDB, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewCategoryRepository(gormDB)
	ctx := context.Background()
	keyword := "Error"

	// Expectations
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "categories" WHERE name ILIKE $1`)).
		WithArgs("%" + keyword + "%").
		WillReturnError(errors.New("database error"))

	// Test
	categories, err := repo.FindByNameLike(ctx, keyword)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if categories != nil {
		t.Errorf("Expected nil categories, got %v", categories)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
