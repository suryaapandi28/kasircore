package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/suryaapandi28/kasircore/internal/entity"
	"github.com/suryaapandi28/kasircore/pkg/cache"
	"gorm.io/gorm"
)

type AccountUserRepository interface {
	CreateAccountUser(AccountUser *entity.AccountUser) (*entity.AccountUser, error)
	UpdateJwtToken(f_kd_account uuid.UUID, token string, expiredAt time.Time) error
	FindAdminByEmail(email string) (*entity.AccountUser, error)
}

type accountuserRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

func NewAccountUserRepository(db *gorm.DB, cacheable cache.Cacheable) *accountuserRepository {
	return &accountuserRepository{db: db, cacheable: cacheable}
}

func (r *accountuserRepository) CreateAccountUser(AccountUser *entity.AccountUser) (*entity.AccountUser, error) {
	if err := r.db.Create(&AccountUser).Error; err != nil {
		return AccountUser, err
	}
	return AccountUser, nil
}
func (r *accountuserRepository) FindAdminByEmail(F_email_account string) (*entity.AccountUser, error) {
	accountprovider := new(entity.AccountUser)
	if err := r.db.Where("f_email_account = ?", F_email_account).Take(accountprovider).Error; err != nil {
		return accountprovider, err
	}
	return accountprovider, nil
}

func (r *accountuserRepository) UpdateJwtToken(
	f_kd_account uuid.UUID,
	token string,
	expiredAt time.Time,
) error {

	return r.db.Model(&entity.AccountUser{}).
		Where("f_kd_account = ?", f_kd_account).
		Updates(map[string]interface{}{
			"f_jwt_token":         token,
			"f_jwt_token_expired": expiredAt,
		}).Error
}
