//go:build integration

package repository_test

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"sample-mcp/db/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sample-mcp/db/entity"
)

func getTestRepo() *repository.AccountRepository {
	return &repository.AccountRepository{
		BaseRepository: &repository.BaseRepository[entity.Account]{DB: TestDB},
	}
}

func createTestAccounts(t *testing.T, accounts ...*entity.Account) *entity.Account {
	if len(accounts) == 0 {
		account := &entity.Account{
			Name:        "Savings Account twothreefourfive",
			AccountType: "SAVINGS",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		accounts = append(accounts, account)
	}

	for _, account := range accounts {
		err := TestDB.Create(account).Error
		require.NoError(t, err)
	}

	return accounts[0]
}

func TestAccountRepository_Create(t *testing.T) {
	repo := getTestRepo()

	account := &entity.Account{
		Name:        "Salary Account",
		AccountType: "CHECKING",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.Create(context.Background(), account)
	assert.NoError(t, err)
	assert.NotZero(t, account.AccountID)
}

func TestAccountRepository_FindByID(t *testing.T) {
	repo := getTestRepo()

	created := createTestAccounts(t)

	found, err := repo.FindByID(context.Background(), created.AccountID)
	assert.NoError(t, err)
	assert.Equal(t, created.Name, found.Name)
	assert.Equal(t, created.AccountType, found.AccountType)
}

func TestAccountRepository_FindAll(t *testing.T) {
	repo := getTestRepo()

	createTestAccounts(t)
	createTestAccounts(t)

	list, err := repo.FindAll(context.Background())
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 2)
}

func TestAccountRepository_Update(t *testing.T) {
	repo := getTestRepo()

	account := createTestAccounts(t)

	account.Name = "Updated Name"
	account.UpdatedAt = time.Now()

	err := repo.Update(context.Background(), account)
	assert.NoError(t, err)

	updated, err := repo.FindByID(context.Background(), account.AccountID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", updated.Name)
}

func TestAccountRepository_Delete(t *testing.T) {
	repo := getTestRepo()

	account := createTestAccounts(t)

	err := repo.Delete(context.Background(), account)
	assert.NoError(t, err)

	_, err = repo.FindByID(context.Background(), account.AccountID)
	assert.Error(t, err)
	assert.EqualError(t, err, gorm.ErrRecordNotFound.Error())
}

func TestAccountRepository_DeleteByID(t *testing.T) {
	repo := getTestRepo()

	account := createTestAccounts(t)

	err := repo.DeleteByID(context.Background(), account.AccountID)
	assert.NoError(t, err)

	_, err = repo.FindByID(context.Background(), account.AccountID)
	assert.Error(t, err)
	assert.EqualError(t, err, gorm.ErrRecordNotFound.Error())
}

func TestAccountRepository_FindByName(t *testing.T) {
	repo := getTestRepo()
	specificAccount := &entity.Account{
		Name:        fmt.Sprintf("Salary Account - %s", time.Now().Format("20060102150405")),
		AccountType: "CHECKING",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	account := createTestAccounts(t, specificAccount)

	found, err := repo.FindByName(context.Background(), account.Name)
	assert.NoError(t, err)
	assert.Equal(t, account.AccountID, found.AccountID)
}

func TestAccountRepository_FindByNameLike(t *testing.T) {
	repo := getTestRepo()

	names := []string{"Groceries Savings", "Vacation Savings", "Random Spending"}
	for _, name := range names {
		repo.Create(context.Background(), &entity.Account{
			Name:        name,
			AccountType: "SAVINGS",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
	}

	results, err := repo.FindByNameLike(context.Background(), "Savings")
	assert.NoError(t, err)
	assert.True(t, len(results) > 3)
}
