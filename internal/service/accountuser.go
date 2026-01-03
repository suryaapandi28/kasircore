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

type AccountUserService interface {
	CreateAccountUser(accountuser *entity.AccountUser) (*entity.AccountUser, error)
	LoginUser(F_email_account string, F_password string) (*entity.AccountUser, error)
	EmailExists(email string) bool
}

type accountuserService struct {
	accountuserRepository repository.AccountUserRepository
	tokenUseCase          token.TokenUseCase
	encryptTool           encrypt.EncryptTool
	emailSender           *email.EmailSender

	WaSender *whatsapp.WhatsappSender
}

func NewAccountUserService(accountuserRepository repository.AccountUserRepository, tokenUseCase token.TokenUseCase,
	encryptTool encrypt.EncryptTool, emailSender *email.EmailSender, WaSender *whatsapp.WhatsappSender) *accountuserService {
	return &accountuserService{
		accountuserRepository: accountuserRepository,
		tokenUseCase:          tokenUseCase,
		encryptTool:           encryptTool,
		emailSender:           emailSender,

		WaSender: WaSender,
	}
}

func (s *accountuserService) CreateAccountUser(accountuser *entity.AccountUser) (*entity.AccountUser, error) {
	if accountuser.F_email_account == "" {
		return nil, errors.New("email cannot be empty")
	}
	if accountuser.F_password == "" {
		return nil, errors.New("password cannot be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(accountuser.F_password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	accountuser.F_password = string(hashedPassword)

	newAdmin, err := s.accountuserRepository.CreateAccountUser(accountuser)
	if err != nil {
		return nil, err
	}

	createaccounttime := time.Now()

	err = s.emailSender.SendWelcomeEmail(accountuser.F_email_account, accountuser.F_nama_account, createaccounttime)
	if err != nil {
		return nil, err
	}

	return newAdmin, nil
}

func (s *accountuserService) EmailExists(F_email_account string) bool {
	_, err := s.accountuserRepository.FindAdminByEmail(F_email_account)
	if err != nil {
		return false
	}
	return true
}

func (s *accountuserService) LoginUser(F_email_account string, F_password string) (*entity.AccountUser, error) {
	DatAccount, err := s.accountuserRepository.FindAdminByEmail(F_email_account)
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
		return nil, errors.New("Akun belum diverifikasi")
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
	err = s.accountuserRepository.UpdateJwtToken(
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
	DatAccount.F_role_accout = claims.Role
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
