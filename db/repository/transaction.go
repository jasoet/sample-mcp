package repository

import (
	"context"
	"gorm.io/gorm"
	"sample-mcp/db/repository/plain"
	"time"

	"sample-mcp/db/entity"
)

type TransactionRepository struct {
	*BaseRepository[entity.Transaction]
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{
		BaseRepository: &BaseRepository[entity.Transaction]{DB: db},
	}
}

func (r *TransactionRepository) FindByAccountID(ctx context.Context, accountID uint) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	if err := r.DB.WithContext(ctx).
		Preload("Account").
		Preload("Category").
		Where("account_id = ?", accountID).
		Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) FindByDateRange(ctx context.Context, start, end time.Time) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	if err := r.DB.WithContext(ctx).
		Preload("Account").
		Preload("Category").
		Where("transaction_date BETWEEN ? AND ?", start, end).
		Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) FindByDescriptionLike(ctx context.Context, keyword string) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	if err := r.DB.WithContext(ctx).
		Preload("Account").
		Preload("Category").
		Where("description IS NOT NULL AND description ILIKE ?", "%"+keyword+"%").
		Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) FindByAccountAndDateRange(
	ctx context.Context,
	accountID uint,
	start, end time.Time,
) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	if err := r.DB.WithContext(ctx).
		Preload("Account").
		Preload("Category").
		Where("account_id = ? AND transaction_date BETWEEN ? AND ?", accountID, start, end).
		Order("transaction_date DESC").
		Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) SumByAccountID(ctx context.Context, accountID uint) (float64, error) {
	var sum float64
	err := r.DB.WithContext(ctx).
		Model(&entity.Transaction{}).
		Where("account_id = ?", accountID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&sum).Error
	return sum, err
}

func (r *TransactionRepository) CountByAccountID(ctx context.Context, accountID uint) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).
		Model(&entity.Transaction{}).
		Where("account_id = ?", accountID).
		Select("COUNT(*)").
		Scan(&count).Error
	return count, err
}

func (r *TransactionRepository) FindLatestForAccount(ctx context.Context, accountID uint, limit int) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	if err := r.DB.WithContext(ctx).
		Preload("Account").
		Preload("Category").
		Where("account_id = ?", accountID).
		Order("transaction_date DESC").
		Limit(limit).
		Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) GroupByCategory(ctx context.Context, accountID uint) ([]plain.TransactionSummary, error) {
	var result []plain.TransactionSummary
	err := r.DB.WithContext(ctx).
		Table("transactions").
		Select("c.name as category_name, SUM(t.amount) as total_amount, COUNT(t.transaction_id) as count").
		Joins("JOIN categories c ON t.category_id = c.category_id").
		Where("t.account_id = ?", accountID).
		Group("c.name").
		Scan(&result).Error
	return result, err
}
