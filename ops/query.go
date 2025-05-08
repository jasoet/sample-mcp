package ops

import (
	"context"
	"sample-mcp/db/entity"
	"sample-mcp/db/repository"
	"sample-mcp/db/repository/plain"
	"time"
)

// QueryOps provides operations for querying data from repositories
type QueryOps struct {
	accountRepo     *repository.AccountRepository
	categoryRepo    *repository.CategoryRepository
	transactionRepo *repository.TransactionRepository
}

// NewQueryOps creates a new QueryOps instance
func NewQueryOps(
	accountRepo *repository.AccountRepository,
	categoryRepo *repository.CategoryRepository,
	transactionRepo *repository.TransactionRepository,
) *QueryOps {
	return &QueryOps{
		accountRepo:     accountRepo,
		categoryRepo:    categoryRepo,
		transactionRepo: transactionRepo,
	}
}

// GetAccountByID retrieves an account by its ID
func (q *QueryOps) GetAccountByID(ctx context.Context, accountID uint) (*entity.Account, error) {
	return q.accountRepo.FindByID(ctx, accountID)
}

// GetAccountByName retrieves an account by its name
func (q *QueryOps) GetAccountByName(ctx context.Context, name string) (*entity.Account, error) {
	return q.accountRepo.FindByName(ctx, name)
}

// SearchAccounts searches for accounts with names containing the keyword
func (q *QueryOps) SearchAccounts(ctx context.Context, keyword string) ([]entity.Account, error) {
	return q.accountRepo.FindByNameLike(ctx, keyword)
}

// GetAllAccounts retrieves all accounts
func (q *QueryOps) GetAllAccounts(ctx context.Context) ([]entity.Account, error) {
	return q.accountRepo.FindAll(ctx)
}

// GetCategoryByID retrieves a category by its ID
func (q *QueryOps) GetCategoryByID(ctx context.Context, categoryID uint) (*entity.Category, error) {
	return q.categoryRepo.FindByID(ctx, categoryID)
}

// GetCategoriesByType retrieves categories by their type
func (q *QueryOps) GetCategoriesByType(ctx context.Context, categoryType string) ([]entity.Category, error) {
	return q.categoryRepo.FindByType(ctx, categoryType)
}

// SearchCategories searches for categories with names containing the keyword
func (q *QueryOps) SearchCategories(ctx context.Context, keyword string) ([]entity.Category, error) {
	return q.categoryRepo.FindByNameLike(ctx, keyword)
}

// GetAllCategories retrieves all categories
func (q *QueryOps) GetAllCategories(ctx context.Context) ([]entity.Category, error) {
	return q.categoryRepo.FindAll(ctx)
}

// GetTransactionByID retrieves a transaction by its ID
func (q *QueryOps) GetTransactionByID(ctx context.Context, transactionID uint) (*entity.Transaction, error) {
	return q.transactionRepo.FindByID(ctx, transactionID)
}

// GetTransactionsByAccountID retrieves all transactions for an account
func (q *QueryOps) GetTransactionsByAccountID(ctx context.Context, accountID uint) ([]entity.Transaction, error) {
	return q.transactionRepo.FindByAccountID(ctx, accountID)
}

// GetTransactionsByDateRange retrieves transactions within a date range
func (q *QueryOps) GetTransactionsByDateRange(ctx context.Context, start, end time.Time) ([]entity.Transaction, error) {
	return q.transactionRepo.FindByDateRange(ctx, start, end)
}

// GetTransactionsByAccountAndDateRange retrieves transactions for an account within a date range
func (q *QueryOps) GetTransactionsByAccountAndDateRange(
	ctx context.Context,
	accountID uint,
	start, end time.Time,
) ([]entity.Transaction, error) {
	return q.transactionRepo.FindByAccountAndDateRange(ctx, accountID, start, end)
}

// SearchTransactionsByDescription searches for transactions with descriptions containing the keyword
func (q *QueryOps) SearchTransactionsByDescription(ctx context.Context, keyword string) ([]entity.Transaction, error) {
	return q.transactionRepo.FindByDescriptionLike(ctx, keyword)
}

// GetAccountBalance calculates the balance for an account
func (q *QueryOps) GetAccountBalance(ctx context.Context, accountID uint) (float64, error) {
	return q.transactionRepo.SumByAccountID(ctx, accountID)
}

// GetTransactionCount gets the number of transactions for an account
func (q *QueryOps) GetTransactionCount(ctx context.Context, accountID uint) (int64, error) {
	return q.transactionRepo.CountByAccountID(ctx, accountID)
}

// GetLatestTransactions gets the latest transactions for an account
func (q *QueryOps) GetLatestTransactions(ctx context.Context, accountID uint, limit int) ([]entity.Transaction, error) {
	return q.transactionRepo.FindLatestForAccount(ctx, accountID, limit)
}

// GetTransactionSummaryByCategory gets transaction summaries grouped by category for an account
func (q *QueryOps) GetTransactionSummaryByCategory(ctx context.Context, accountID uint) ([]plain.TransactionSummary, error) {
	return q.transactionRepo.GroupByCategory(ctx, accountID)
}

// GetAllTransactions retrieves all transactions
func (q *QueryOps) GetAllTransactions(ctx context.Context) ([]entity.Transaction, error) {
	return q.transactionRepo.FindAll(ctx)
}
