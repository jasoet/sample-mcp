//go:build integration

package repository_test

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sample-mcp/db/entity"
	"sample-mcp/db/repository"
)

// Helper to get a fresh CategoryRepository
func getCategoryRepo() *repository.CategoryRepository {
	return repository.NewCategoryRepository(TestDB)
}

// Helper to create a test category in DB
func createTestCategory(t *testing.T, name, categoryType string) *entity.Category {
	category := &entity.Category{
		Name:         fmt.Sprintf("%s - %s", name, time.Now().Format("20060102150405")),
		CategoryType: categoryType,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err := TestDB.Create(category).Error
	require.NoError(t, err)
	return category
}

// Test Create method
func TestCategoryRepository_Create(t *testing.T) {
	repo := getCategoryRepo()

	category := &entity.Category{
		Name:         fmt.Sprintf("Groceries - %s", time.Now().Format("20060102150405")),
		CategoryType: "EXPENSE",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := repo.Create(context.Background(), category)
	assert.NoError(t, err)
	assert.NotZero(t, category.CategoryID)
}

// Test FindByID method
func TestCategoryRepository_FindByID(t *testing.T) {
	repo := getCategoryRepo()

	created := createTestCategory(t, "Entertainment", "EXPENSE")

	found, err := repo.FindByID(context.Background(), created.CategoryID)
	assert.NoError(t, err)
	assert.Equal(t, created.Name, found.Name)
	assert.Equal(t, created.CategoryType, found.CategoryType)
}

// Test FindAll method
func TestCategoryRepository_FindAll(t *testing.T) {
	repo := getCategoryRepo()

	createTestCategory(t, "Rent", "EXPENSE")
	createTestCategory(t, "Salary", "INCOME")

	list, err := repo.FindAll(context.Background())
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 2)
}

// Test Update method
func TestCategoryRepository_Update(t *testing.T) {
	repo := getCategoryRepo()

	category := createTestCategory(t, "Old Name", "INCOME")

	updateName := fmt.Sprintf("Updated Name - %s", time.Now().Format("20060102150405"))
	category.Name = updateName
	category.UpdatedAt = time.Now()

	err := repo.Update(context.Background(), category)
	assert.NoError(t, err)

	updated, err := repo.FindByID(context.Background(), category.CategoryID)
	assert.NoError(t, err)
	assert.Equal(t, updateName, updated.Name)
}

// Test Delete and DeleteByID methods
func TestCategoryRepository_Delete(t *testing.T) {
	repo := getCategoryRepo()

	category := createTestCategory(t, "To Delete", "EXPENSE")

	err := repo.Delete(context.Background(), category)
	assert.NoError(t, err)

	_, err = repo.FindByID(context.Background(), category.CategoryID)
	assert.Error(t, err)
	assert.EqualError(t, err, gorm.ErrRecordNotFound.Error())
}

func TestCategoryRepository_DeleteByID(t *testing.T) {
	repo := getCategoryRepo()

	category := createTestCategory(t, "To Delete by ID", "INCOME")

	err := repo.DeleteByID(context.Background(), category.CategoryID)
	assert.NoError(t, err)

	_, err = repo.FindByID(context.Background(), category.CategoryID)
	assert.Error(t, err)
	assert.EqualError(t, err, gorm.ErrRecordNotFound.Error())
}

// Test FindByType method
func TestCategoryRepository_FindByType(t *testing.T) {
	repo := getCategoryRepo()

	createTestCategory(t, "Food", "EXPENSE")
	createTestCategory(t, "Transport", "EXPENSE")
	createTestCategory(t, "Investment", "INCOME")

	results, err := repo.FindByType(context.Background(), "EXPENSE")
	assert.NoError(t, err)
	assert.True(t, len(results) > 2)
}

// Test FindByNameLike method
func TestCategoryRepository_FindByNameLike(t *testing.T) {
	repo := getCategoryRepo()

	names := []string{"Groceries", "Gas Bill", "Gym Membership", "Savings"}
	for _, name := range names {
		repo.Create(context.Background(), &entity.Category{
			Name:         name,
			CategoryType: "EXPENSE",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		})
	}

	results, err := repo.FindByNameLike(context.Background(), "Gy")
	assert.NoError(t, err)
	assert.True(t, len(results) > 0)
}
