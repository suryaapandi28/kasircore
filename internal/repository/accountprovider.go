package repository

import (
	"github.com/suryaapandi28/kasircore/internal/entity"
	"github.com/suryaapandi28/kasircore/pkg/cache"
	"gorm.io/gorm"
)

type AccountproviderRepository interface {
	CreateAdmin(ProviderAccount *entity.ProviderAccount) (*entity.ProviderAccount, error)
	FindAdminByEmail(email string) (*entity.ProviderAccount, error)
}

type accountproviderRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

func NewAccountproviderRepository(db *gorm.DB, cacheable cache.Cacheable) *accountproviderRepository {
	return &accountproviderRepository{db: db, cacheable: cacheable}
}

func (r *accountproviderRepository) CreateAdmin(ProviderAccount *entity.ProviderAccount) (*entity.ProviderAccount, error) {
	if err := r.db.Create(&ProviderAccount).Error; err != nil {
		return ProviderAccount, err
	}
	return ProviderAccount, nil
}
func (r *accountproviderRepository) FindAdminByEmail(email string) (*entity.ProviderAccount, error) {
	admin := new(entity.ProviderAccount)
	if err := r.db.Where("email = ?", email).Take(admin).Error; err != nil {
		return admin, err
	}
	return admin, nil
}
