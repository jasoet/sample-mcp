package repository

import (
	"context"
	"gorm.io/gorm"
	"sample-mcp/db/entity"
)

type AccountRepository struct {
	*BaseRepository[entity.Account]
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{
		BaseRepository: &BaseRepository[entity.Account]{DB: db},
	}
}

func (r *AccountRepository) FindByName(ctx context.Context, name string) (*entity.Account, error) {
	var account entity.Account
	if err := r.DB.WithContext(ctx).Where("name = ?", name).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) FindByNameLike(ctx context.Context, keyword string) ([]entity.Account, error) {
	var accounts []entity.Account
	if err := r.DB.WithContext(ctx).
		Where("name ILIKE ?", "%"+keyword+"%").
		Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}
