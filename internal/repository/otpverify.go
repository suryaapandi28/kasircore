package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/suryaapandi28/kasircore/internal/entity"

	"github.com/suryaapandi28/kasircore/pkg/cache"
	"gorm.io/gorm"
)

type OTPRepository interface {
	Create(ctx context.Context, otp *entity.OtpVerify) error
	FindByEmail(email string) (*entity.ProviderAccount, error)
	FindbyKdAccount(F_kd_account uuid.UUID) (*entity.OtpVerify, error)

	SaveOtp(otp *entity.OtpVerify) error
	UpdateOtp(otp *entity.OtpVerify) error
}

type otpRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

func NewOTPRepository(db *gorm.DB, cacheable cache.Cacheable) OTPRepository {
	return &otpRepository{db: db, cacheable: cacheable}
}

func (r *otpRepository) Create(ctx context.Context, otp *entity.OtpVerify) error {
	return r.db.WithContext(ctx).Create(otp).Error
}
func (r *otpRepository) FindByEmail(F_email_account string) (*entity.ProviderAccount, error) {
	accountprovider := new(entity.ProviderAccount)
	if err := r.db.Where("f_email_account = ?", F_email_account).Take(accountprovider).Error; err != nil {
		return accountprovider, err
	}
	return accountprovider, nil
}

func (r *otpRepository) SaveOtp(otp *entity.OtpVerify) error {
	return r.db.Create(otp).Error
}

func (r *otpRepository) UpdateOtp(otp *entity.OtpVerify) error {
	return r.db.Save(otp).Error
}

func (r *otpRepository) FindbyKdAccount(F_kd_account uuid.UUID) (*entity.OtpVerify, error) {
	otpverify := new(entity.OtpVerify)
	if err := r.db.Where("f_kd_account = ?", F_kd_account).Take(otpverify).Error; err != nil {
		return otpverify, err
	}
	return otpverify, nil
}
