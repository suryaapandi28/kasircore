package service

import (
	"errors" // Import log package
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/suryaapandi28/kasircore/internal/entity"
	"github.com/suryaapandi28/kasircore/internal/repository"
	"github.com/suryaapandi28/kasircore/pkg/email"
	"github.com/suryaapandi28/kasircore/pkg/encrypt"
	"github.com/suryaapandi28/kasircore/pkg/token"
	"github.com/suryaapandi28/kasircore/pkg/whatsapp"
	"golang.org/x/crypto/bcrypt"
)

type AccountproviderService interface {
	CreateAdmin(accountprovider *entity.ProviderAccount) (*entity.ProviderAccount, error)
	LoginProvider(F_email_account string, F_password string) (*entity.ProviderAccount, error)
	EmailExists(email string) bool
}

type accountproviderService struct {
	accountproviderRepository repository.AccountproviderRepository
	tokenUseCase              token.TokenUseCase
	encryptTool               encrypt.EncryptTool
	emailSender               *email.EmailSender

	WaSender *whatsapp.WhatsappSender
}
type LoginProviderTokenResponse struct {
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewAccountproviderService(accountproviderRepository repository.AccountproviderRepository, tokenUseCase token.TokenUseCase,
	encryptTool encrypt.EncryptTool, emailSender *email.EmailSender, WaSender *whatsapp.WhatsappSender) *accountproviderService {
	return &accountproviderService{
		accountproviderRepository: accountproviderRepository,
		tokenUseCase:              tokenUseCase,
		encryptTool:               encryptTool,
		emailSender:               emailSender,

		WaSender: WaSender,
	}
}

func (s *accountproviderService) CreateAdmin(accountprovider *entity.ProviderAccount) (*entity.ProviderAccount, error) {
	if accountprovider.F_email_account == "" {
		return nil, errors.New("email cannot be empty")
	}
	if accountprovider.F_password == "" {
		return nil, errors.New("password cannot be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(accountprovider.F_password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	accountprovider.F_password = string(hashedPassword)

	newAdmin, err := s.accountproviderRepository.CreateAccountProvider(accountprovider)
	if err != nil {
		return nil, err
	}

	createaccounttime := time.Now()

	err = s.emailSender.SendWelcomeEmail(accountprovider.F_email_account, accountprovider.F_nama_account, createaccounttime)
	if err != nil {
		return nil, err
	}

	return newAdmin, nil
}

func (s *accountproviderService) EmailExists(F_email_account string) bool {
	_, err := s.accountproviderRepository.FindAdminByEmail(F_email_account)
	if err != nil {
		return false
	}
	return true
}

func (s *accountproviderService) LoginProvider(F_email_account string, F_password string) (*entity.ProviderAccount, error) {
	DatAccount, err := s.accountproviderRepository.FindAdminByEmail(F_email_account)
	if err != nil {
		return nil, errors.New("Email tidak terdaftar")
	}

	// Validasi password
	if err := bcrypt.CompareHashAndPassword(
		[]byte(DatAccount.F_password),
		[]byte(F_password),
	); err != nil {
		return nil, errors.New("Password salah")
	}

	// Validasi verifikasi akun
	if !DatAccount.F_verification_account {
		return nil, errors.New("Akun provider belum diverifikasi")
	}

	// ===== JWT CLAIMS =====
	expiredAt := time.Now().Add(24 * time.Hour)

	claims := token.JwtCustomClaims{
		ID:    DatAccount.F_kd_account.String(),
		Email: DatAccount.F_email_account, // 👈 EMAIL
		Role:  DatAccount.F_role_accout,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredAt), // 👈 EXP
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Generate token
	tokenString, err := s.tokenUseCase.GenerateAccessToken(claims)
	if err != nil {
		return nil, errors.New("Gagal menghasilkan token")
	}

	// ===== UPDATE KE DATABASE =====
	err = s.accountproviderRepository.UpdateJwtToken(
		DatAccount.F_kd_account,
		tokenString,
		expiredAt,
	)
	if err != nil {
		return nil, errors.New("Gagal menyimpan token")
	}

	// ===== TEMPEL KE ENTITY =====
	DatAccount.F_jwt_token = tokenString
	DatAccount.F_email_account = claims.Email
	DatAccount.F_jwt_token_expired = expiredAt
	message := fmt.Sprintf(
		"Halo %s,\n\nKami mendeteksikan aktivitas login baru pada akun Anda.\n\n"+
			"Jika ini adalah Anda, silakan abaikan pesan ini.\n\n"+
			"Jika bukan Anda, segera amankan akun Anda.\n\n"+
			"Terima kasih,\nTim APSM Indonesia Global",
		DatAccount.F_nama_account,
	)

	err = s.WaSender.SendMessage(
		DatAccount.F_phone_account,
		message,
	)
	if err != nil {
		return nil, err
	}
	return DatAccount, nil
}
