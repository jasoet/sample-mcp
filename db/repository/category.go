package repository

import (
	"context"
	"gorm.io/gorm"
	"sample-mcp/db/entity"
)

type CategoryRepository struct {
	*BaseRepository[entity.Category]
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		BaseRepository: &BaseRepository[entity.Category]{DB: db},
	}
}

func (r *CategoryRepository) FindByType(ctx context.Context, categoryType string) ([]entity.Category, error) {
	var categories []entity.Category
	if err := r.DB.WithContext(ctx).Where("category_type = ?", categoryType).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
func (r *CategoryRepository) FindByNameLike(ctx context.Context, keyword string) ([]entity.Category, error) {
	var categories []entity.Category
	if err := r.DB.WithContext(ctx).
		Where("name ILIKE ?", "%"+keyword+"%").
		Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
