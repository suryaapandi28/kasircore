package repository

import (
	"context"

	"github.com/suryaapandi28/kasircore/internal/entity"

	"github.com/suryaapandi28/kasircore/pkg/cache"
	"gorm.io/gorm"
)

type OTPRepository interface {
	Create(ctx context.Context, otp *entity.OtpVerify) error
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
