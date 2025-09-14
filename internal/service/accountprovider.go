package service

import (
	"errors" // Import log package
	"time"

	"github.com/suryaapandi28/kasircore/internal/entity"
	"github.com/suryaapandi28/kasircore/internal/repository"
	"github.com/suryaapandi28/kasircore/pkg/email"
	"github.com/suryaapandi28/kasircore/pkg/encrypt"
	"github.com/suryaapandi28/kasircore/pkg/token"
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
}

func NewAccountproviderService(accountproviderRepository repository.AccountproviderRepository, tokenUseCase token.TokenUseCase,
	encryptTool encrypt.EncryptTool, emailSender *email.EmailSender) *accountproviderService {
	return &accountproviderService{
		accountproviderRepository: accountproviderRepository,
		tokenUseCase:              tokenUseCase,
		encryptTool:               encryptTool,
		emailSender:               emailSender,
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

	err = bcrypt.CompareHashAndPassword([]byte(DatAccount.F_password), []byte(F_password))
	if err != nil {
		return nil, errors.New("Password salah")
	}

	if !DatAccount.F_verification_account {
		return nil, errors.New("Akun provider belum diverifikasi")
	}

	return DatAccount, nil
}
