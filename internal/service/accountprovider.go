package service

import (
	"errors" // Import log package

	"github.com/suryaapandi28/kasircore/internal/entity"
	"github.com/suryaapandi28/kasircore/internal/repository"
	"github.com/suryaapandi28/kasircore/pkg/encrypt"
	"github.com/suryaapandi28/kasircore/pkg/token"
	"golang.org/x/crypto/bcrypt"
)

type AccountproviderService interface {
	CreateAdmin(admin *entity.ProviderAccount) (*entity.ProviderAccount, error)
	EmailExists(email string) bool
}

type accountproviderService struct {
	accountproviderRepository repository.AccountproviderRepository
	tokenUseCase              token.TokenUseCase
	encryptTool               encrypt.EncryptTool
}

func NewAccountproviderService(accountproviderRepository repository.AccountproviderRepository, tokenUseCase token.TokenUseCase,
	encryptTool encrypt.EncryptTool) *accountproviderService {
	return &accountproviderService{
		accountproviderRepository: accountproviderRepository,
		tokenUseCase:              tokenUseCase,
		encryptTool:               encryptTool,
	}
}

func (s *accountproviderService) CreateAdmin(admin *entity.ProviderAccount) (*entity.ProviderAccount, error) {
	if admin.Email == "" {
		return nil, errors.New("email cannot be empty")
	}
	if admin.Password == "" {
		return nil, errors.New("password cannot be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	admin.Password = string(hashedPassword)

	newAdmin, err := s.accountproviderRepository.CreateAdmin(admin)
	if err != nil {
		return nil, err
	}

	return newAdmin, nil
}

func (s *accountproviderService) EmailExists(email string) bool {
	_, err := s.accountproviderRepository.FindAdminByEmail(email)
	if err != nil {
		return false
	}
	return true
}
