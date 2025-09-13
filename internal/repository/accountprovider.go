package repository

import (
	"github.com/suryaapandi28/kasircore/internal/entity"
	"github.com/suryaapandi28/kasircore/pkg/cache"
	"gorm.io/gorm"
)

type AccountproviderRepository interface {
	CreateAccountProvider(ProviderAccount *entity.ProviderAccount) (*entity.ProviderAccount, error)
	FindAdminByEmail(email string) (*entity.ProviderAccount, error)
}

type accountproviderRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

func NewAccountproviderRepository(db *gorm.DB, cacheable cache.Cacheable) *accountproviderRepository {
	return &accountproviderRepository{db: db, cacheable: cacheable}
}

func (r *accountproviderRepository) CreateAccountProvider(ProviderAccount *entity.ProviderAccount) (*entity.ProviderAccount, error) {
	if err := r.db.Create(&ProviderAccount).Error; err != nil {
		return ProviderAccount, err
	}
	return ProviderAccount, nil
}
func (r *accountproviderRepository) FindAdminByEmail(F_email_account string) (*entity.ProviderAccount, error) {
	accountprovider := new(entity.ProviderAccount)
	if err := r.db.Where("f_email_account = ?", F_email_account).Take(accountprovider).Error; err != nil {
		return accountprovider, err
	}
	return accountprovider, nil
}
