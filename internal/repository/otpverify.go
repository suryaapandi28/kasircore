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
	FindOTPbyKdAccount(F_kd_account uuid.UUID) (*entity.OtpVerify, error)
	FindPhoneByKdAccount(F_kd_account uuid.UUID) (string, error)
	UpdateAccountVerified(F_kd_account uuid.UUID, status bool) error
	DeleteOTPByKdAccount(F_kd_account uuid.UUID) error

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

func (r *otpRepository) FindOTPbyKdAccount(F_kd_account uuid.UUID) (*entity.OtpVerify, error) {
	otpverify := new(entity.OtpVerify)
	if err := r.db.Where("f_kd_account = ?", F_kd_account).Take(otpverify).Error; err != nil {
		return otpverify, err
	}
	return otpverify, nil
}

func (r *otpRepository) FindPhoneByKdAccount(F_kd_account uuid.UUID) (string, error) {
	var phone string

	err := r.db.
		Table("accounts_providers").
		Select("f_phone_account").
		Where("f_kd_account = ?", F_kd_account).
		Take(&phone).Error

	if err != nil {
		return "", err
	}

	return phone, nil
}
func (r *otpRepository) UpdateAccountVerified(F_kd_account uuid.UUID, status bool) error {
	return r.db.Model(&entity.ProviderAccount{}).
		Where("f_kd_account = ?", F_kd_account).
		Update("f_verification_account", status).
		Error
}

func (r *otpRepository) DeleteOTPByKdAccount(
	F_kd_account uuid.UUID,
) error {
	return r.db.
		Where("f_kd_account = ?", F_kd_account).
		Delete(&entity.OtpVerify{}).
		Error
}
