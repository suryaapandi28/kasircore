package repository

import (
	"github.com/suryaapandi28/kasircore/internal/entity"
	"github.com/suryaapandi28/kasircore/pkg/cache"
	"gorm.io/gorm"
)

type MerchantRepository interface {
	CreateMerchant(merchant *entity.Merchant) (*entity.Merchant, error)
}

type merchantRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

func NewMerchantRepository(db *gorm.DB, cacheable cache.Cacheable) *merchantRepository {
	return &merchantRepository{db: db, cacheable: cacheable}
}

func (r *merchantRepository) CreateMerchant(merchant *entity.Merchant) (*entity.Merchant, error) {
	if err := r.db.Create(&merchant).Error; err != nil {
		return merchant, err
	}
	return merchant, nil
}
