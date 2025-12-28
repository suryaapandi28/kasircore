package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/suryaapandi28/kasircore/internal/entity"
	"github.com/suryaapandi28/kasircore/internal/repository"
	"github.com/suryaapandi28/kasircore/pkg/email"
	"github.com/suryaapandi28/kasircore/pkg/whatsapp"
)

type OtpService interface {
	GenerateOtp(F_email_account string, via string) (*entity.OtpVerify, error)
	VerifyOtpRequest(F_email_account string, F_kode_otp string) (*entity.OtpVerify, error)
}

type otpService struct {
	otpRepo     repository.OTPRepository
	emailSender *email.EmailSender
	WaSender    *whatsapp.WhatsappSender
}

func NewOtpService(repo repository.OTPRepository, emailSender *email.EmailSender, WaSender *whatsapp.WhatsappSender) OtpService {
	return &otpService{
		otpRepo:     repo,
		emailSender: emailSender,
		WaSender:    WaSender,
	}
}

func (s *otpService) GenerateOtp(F_email_account string, via string) (*entity.OtpVerify, error) {
	// cek via email or whatsapp
	if via != "email" && via != "whatsapp" {
		return nil, errors.New("via harus email atau whatsapp")
	}

	// cek account by email
	DatOTPAccount, err := s.otpRepo.FindByEmail(F_email_account)
	if err != nil {
		return nil, errors.New("Email tidak terdaftar")
	}

	// ambil phone berdasarkan kd_account (UNTUK WA)
	var phone string
	if via == "whatsapp" {
		phone, err = s.otpRepo.FindPhoneByKdAccount(
			DatOTPAccount.F_kd_account,
		)
		if err != nil || phone == "" {
			return nil, errors.New("nomor whatsapp tidak ditemukan")
		}
	}
	// generate random 6 digit OTP
	otpCode := fmt.Sprintf("%06d", rand.Intn(1000000))

	// cek apakah sudah ada OTP untuk account ini
	existingOtp, err := s.otpRepo.FindOTPbyKdAccount(DatOTPAccount.F_kd_account)
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

		// kirim whatsapp OTP
		if via == "whatsapp" {
			message := fmt.Sprintf(
				"Halo %s,\n\nTerima kasih telah mendaftar di APSM Indonesia Global.\n\n"+
					"Silakan gunakan kode berikut untuk memverifikasi akun Anda:\n\n"+
					"```%s```\n\n"+
					"Kode ini akan kedaluwarsa dalam 5 menit.\n\n"+
					"Jika Anda tidak meminta ini, mohon abaikan pesan ini.\n\n"+
					"Salam,\nTim APSM Indonesia Global",
				DatOTPAccount.F_nama_account,
				otpCode,
			)

			err = s.WaSender.SendMessage(
				DatOTPAccount.F_phone_account,
				message,
			)
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

	// kirim whatsapp OTP
	if via == "whatsapp" {
		message := fmt.Sprintf(
			"Halo %s,\n\nTerima kasih telah mendaftar di APSM Indonesia Global.\n\n"+
				"Silakan gunakan kode berikut untuk memverifikasi akun Anda:\n\n"+
				"```%s```\n\n"+
				"Kode ini akan kedaluwarsa dalam 5 menit.\n\n"+
				"Jika Anda tidak meminta ini, mohon abaikan pesan ini.\n\n"+
				"Salam,\nTim APSM Indonesia Global",
			DatOTPAccount.F_nama_account,
			otpCode,
		)

		err = s.WaSender.SendMessage(
			DatOTPAccount.F_phone_account,
			message,
		)
		if err != nil {
			return nil, err
		}
	}

	return otp, nil

}

func (s *otpService) VerifyOtpRequest(F_email_account string, F_kode_otp string) (*entity.OtpVerify, error) {
	// 1. validasi input
	if F_email_account == "" {
		return nil, errors.New("Email tidak terbaca")
	}
	if F_kode_otp == "" {
		return nil, errors.New("OTP tidak terbaca")
	}

	// 2. cek account
	account, err := s.otpRepo.FindByEmail(F_email_account)
	if err != nil {
		return nil, errors.New("Email tidak terdaftar")
	}

	// 3. ambil OTP
	otpRecord, err := s.otpRepo.FindOTPbyKdAccount(account.F_kd_account)
	if err != nil {
		return nil, errors.New("Data OTP tidak ditemukan")
	}

	// 4. cek OTP cocok
	if otpRecord.F_kode_otp != F_kode_otp {
		return nil, errors.New("Kode OTP tidak cocok")
	}

	// 5. cek expired
	if time.Now().After(otpRecord.F_otp_expired) {
		return nil, errors.New("Kode OTP sudah kedaluwarsa")
	}

	// ===============================
	// 6. UPDATE ACCOUNT → VERIFIED
	// ===============================
	if err := s.otpRepo.UpdateAccountVerified(account.F_kd_account, true); err != nil {
		return nil, errors.New("Gagal memverifikasi akun")
	}

	// ===============================
	// 7. HAPUS OTP (HARD DELETE)
	// ===============================
	_ = s.otpRepo.DeleteOTPByKdAccount(account.F_kd_account)

	// ===============================
	// 8. NOTIFIKASI (OPTIONAL)
	// ===============================

	// WhatsApp
	if account.F_phone_account != "" {
		message := fmt.Sprintf(
			"Halo %s,\n\n"+
				"Akun Anda di APSM Indonesia Global berhasil diverifikasi ✅\n\n"+
				"Terima kasih telah menggunakan layanan kami.\n\n"+
				"Salam,\nTim APSM Indonesia Global",
			account.F_nama_account,
		)
		_ = s.WaSender.SendMessage(account.F_phone_account, message)
	}

	// Email
	if account.F_email_account != "" {
		_ = s.emailSender.SendSuccessVerificationEmail(
			account.F_email_account,
			account.F_nama_account,
		)
	}

	// 9. sukses
	return otpRecord, nil

}
