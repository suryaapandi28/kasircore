package service

import (
	"math/rand"
	"time"

	"github.com/suryaapandi28/kasircore/internal/entity"
	"github.com/suryaapandi28/kasircore/internal/repository"

	"github.com/google/uuid"
)

type OtpService interface {
	GenerateOtp(accountID uuid.UUID, via string) (*entity.OtpVerify, error)
}

type otpService struct {
	otpRepo repository.OtpRepository
}

func NewOtpService(repo repository.OtpRepository) OtpService {
	return &otpService{otpRepo: repo}
}
func (s *otpService) GenerateOtp(accountID uuid.UUID, via string) (*entity.OtpVerify, error) {
	// generate random 6 digit OTP
	rand.Seed(time.Now().UnixNano())
	otpCode := rand.Intn(900000) + 100000 // 6 digit

	otp := &entity.OtpVerify{
		F_kd_account:  accountID,
		F_kode_otp:    string(rune(otpCode)),
		F_otp_expired: time.Now().Add(5 * time.Minute),
		F_otp_via:     via,
	}

	if err := s.otpRepo.SaveOtp(otp); err != nil {
		return nil, err
	}

	return otp, nil
}
