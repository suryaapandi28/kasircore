package repository

import (
	"github.com/suryaapandi28/kasircore/internal/entity"
	"github.com/suryaapandi28/kasircore/pkg/cache"
	"gorm.io/gorm"
)

type MerchantRepository interface {
	Create(merchant *entity.Merchant) error
	FindAll() ([]entity.Merchant, error)
	FindByID(id string) (*entity.Merchant, error)
	Update(merchant *entity.Merchant) error
	Delete(id string) error
}

type merchantRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

func NewMerchantRepository(db *gorm.DB, cacheable cache.Cacheable) MerchantRepository {
	return &merchantRepository{db, cacheable}
}

func (r *merchantRepository) Create(merchant *entity.Merchant) error {
	return r.db.Create(merchant).Error
}

func (r *merchantRepository) FindAll() ([]entity.Merchant, error) {
	var merchants []entity.Merchant
	err := r.db.Find(&merchants).Error
	return merchants, err
}

func (r *merchantRepository) FindByID(id string) (*entity.Merchant, error) {
	var merchant entity.Merchant
	err := r.db.First(&merchant, "f_kode_merchant = ?", id).Error
	return &merchant, err
}

func (r *merchantRepository) Update(merchant *entity.Merchant) error {
	return r.db.Save(merchant).Error
}

func (r *merchantRepository) Delete(id string) error {
	return r.db.Delete(&entity.Merchant{}, "f_kode_merchant = ?", id).Error
}
