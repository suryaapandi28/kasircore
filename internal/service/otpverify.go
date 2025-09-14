package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/suryaapandi28/kasircore/internal/entity"
	"github.com/suryaapandi28/kasircore/internal/repository"
	"github.com/suryaapandi28/kasircore/pkg/email"
)

type OtpService interface {
	GenerateOtp(F_email_account string, via string) (*entity.OtpVerify, error)
}

type otpService struct {
	otpRepo     repository.OTPRepository
	emailSender *email.EmailSender
}

func NewOtpService(repo repository.OTPRepository, emailSender *email.EmailSender) OtpService {
	return &otpService{
		otpRepo:     repo,
		emailSender: emailSender,
	}
}

func (s *otpService) GenerateOtp(F_email_account string, via string) (*entity.OtpVerify, error) {
	// cek account by email
	DatOTPAccount, err := s.otpRepo.FindByEmail(F_email_account)
	if err != nil {
		return nil, errors.New("Email tidak terdaftar")
	}

	// generate random 6 digit OTP
	otpCode := fmt.Sprintf("%06d", rand.Intn(1000000))

	// cek apakah sudah ada OTP untuk account ini
	existingOtp, err := s.otpRepo.FindbyKdAccount(DatOTPAccount.F_kd_account)
	if err == nil && existingOtp != nil {
		// update record yang sudah ada
		existingOtp.F_kode_otp = otpCode
		existingOtp.F_otp_expired = time.Now().Add(5 * time.Minute)
		existingOtp.F_otp_via = via

		if err := s.otpRepo.UpdateOtp(existingOtp); err != nil {
			return nil, err
		}

		// kirim email OTP
		if via == "email" {
			err = s.emailSender.SendVerificationEmail(DatOTPAccount.F_email_account, DatOTPAccount.F_nama_account, otpCode)
			if err != nil {
				return nil, err
			}
		}
		return existingOtp, nil
	}

	// kalau tidak ada, buat baru
	otp := &entity.OtpVerify{
		F_kd_account:  DatOTPAccount.F_kd_account,
		F_kode_otp:    otpCode,
		F_otp_expired: time.Now().Add(5 * time.Minute),
		F_otp_via:     via,
	}

	if err := s.otpRepo.SaveOtp(otp); err != nil {
		return nil, err
	}

	// kirim email OTP
	if via == "email" {
		err = s.emailSender.SendVerificationEmail(DatOTPAccount.F_email_account, DatOTPAccount.F_nama_account, otpCode)
		if err != nil {
			return nil, err
		}
	}

	return otp, nil
}
