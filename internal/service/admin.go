package service

import (
	"errors" // Import log package
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/suryaapandi28/kasircore/internal/entity"
	"github.com/suryaapandi28/kasircore/internal/repository"
	"github.com/suryaapandi28/kasircore/pkg/email"
	"github.com/suryaapandi28/kasircore/pkg/encrypt"
	"github.com/suryaapandi28/kasircore/pkg/token"
	"golang.org/x/crypto/bcrypt"
)

type AdminService interface {
	LoginAdmin(email string, password string) (string, error)
	FindAllUser() ([]entity.User, error)
	CreateAdmin(admin *entity.Admin) (*entity.Admin, error)
	UpdateAdmin(admin *entity.Admin) (*entity.Admin, error)
	DeleteAdmin(admin uuid.UUID) (bool, error)
	EmailExists(email string) bool
	CheckUserExists(id uuid.UUID) (bool, error)
}

type adminService struct {
	adminRepository     repository.AdminRepository
	tokenUseCase        token.TokenUseCase
	encryptTool         encrypt.EncryptTool
	emailSender         *email.EmailSender
	notificationService NotificationService
}

func NewAdminService(adminRepository repository.AdminRepository, tokenUseCase token.TokenUseCase,
	encryptTool encrypt.EncryptTool, emailSender *email.EmailSender, notificationService NotificationService) *adminService {
	return &adminService{
		adminRepository:     adminRepository,
		tokenUseCase:        tokenUseCase,
		encryptTool:         encryptTool,
		emailSender:         emailSender,
		notificationService: notificationService,
	}
}

func (s *adminService) LoginAdmin(email string, password string) (string, error) {
	admin, err := s.adminRepository.FindAdminByEmail(email)
	if err != nil {
		return "", errors.New("wrong input email/password")
	}
	if admin.Role != "admin" {
		return "", errors.New("you dont have access")
	}
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
	if err != nil {
		return "", errors.New("wrong input email/password")
	}

	expiredTime := time.Now().Local().Add(24 * time.Hour)

	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err)
	}

	// Dekripsi nomor telepon jika perlu
	admin.Phone, _ = s.encryptTool.Decrypt(admin.Phone)
	expiredTimeInJakarta := expiredTime.In(location)
	// Buat claims JWT
	claims := token.JwtCustomClaims{
		ID:    admin.User_ID.String(),
		Email: admin.Email,
		Role:  "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Depublic",
			ExpiresAt: jwt.NewNumericDate(expiredTimeInJakarta),
		},
	}

	// Generate JWT token
	jwtToken, err := s.tokenUseCase.GenerateAccessToken(claims)
	if err != nil {
		return "", errors.New("there is an error in the system")
	}

	// nyimpen token JWT dan waktu ke database
	admin.JwtToken = jwtToken
	admin.JwtTokenExpiresAt = expiredTime

	// update token JWT dan waktu di database
	if err := s.adminRepository.UpdateAdminJwtToken(admin.User_ID, jwtToken, expiredTime); err != nil {
		return "", errors.New("failed to update user token info")
	}

	// double check buat bandingin JWT yang dibuat dengan JWT yang tersimpan di database
	if admin.JwtToken != jwtToken {
		return "", errors.New("JWT token mismatch")
	}
	return jwtToken, nil
}

func (s *adminService) FindAllUser() ([]entity.User, error) {
	admin, err := s.adminRepository.FindAllUser()
	if err != nil {
		return nil, err
	}

	formattedAdmin := make([]entity.User, 0)
	for _, v := range admin {
		v.Phone, _ = s.encryptTool.Decrypt(v.Phone)
		formattedAdmin = append(formattedAdmin, v)
	}

	return formattedAdmin, nil
}

func (s *adminService) CreateAdmin(admin *entity.Admin) (*entity.Admin, error) {
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

	newAdmin, err := s.adminRepository.CreateAdmin(admin)
	if err != nil {
		return nil, err
	}
	emailAddr := newAdmin.Email
	err = s.emailSender.SendWelcomeEmail(emailAddr, newAdmin.Fullname, "")

	if err != nil {
		return nil, err
	}

	resetCode := generateResetCode()
	err = s.emailSender.SendVerificationEmail(newAdmin.Email, newAdmin.Fullname, resetCode)
	if err != nil {
		return nil, err
	}
	err = s.adminRepository.SaveVerifCode(newAdmin.User_ID, resetCode)
	if err != nil {
		return nil, err
	}

	notification := &entity.Notification{
		UserID:  newAdmin.User_ID,
		Type:    "Registration",
		Message: "Admin registration successful",
		IsRead:  false,
	}
	err = s.notificationService.CreateNotification(notification)
	if err != nil {
		return nil, err
	}

	return newAdmin, nil
}

func (s *adminService) CheckUserExists(id uuid.UUID) (bool, error) {
	return s.adminRepository.CheckUserExists(id)
}

func (s *adminService) UpdateAdmin(admin *entity.Admin) (*entity.Admin, error) {
	if admin.Email == "" {
		return nil, errors.New("email cannot be empty")
	}
	if admin.Password == "" {
		return nil, errors.New("password cannot be empty")
	}
	if admin.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		admin.Password = string(hashedPassword)
	}
	if admin.Phone != "" {
		admin.Phone, _ = s.encryptTool.Encrypt(admin.Phone)
	}

	updatedAdmin, err := s.adminRepository.UpdateAdmin(admin)
	if err != nil {
		return nil, err
	}

	notification := &entity.Notification{
		UserID:  updatedAdmin.User_ID,
		Type:    "Update Profile",
		Message: "Update Profile successful",
		IsRead:  false,
	}
	err = s.notificationService.CreateNotification(notification)
	if err != nil {
		return nil, err
	}
	return s.adminRepository.UpdateAdmin(admin)
}

func (s *adminService) DeleteAdmin(user_Id uuid.UUID) (bool, error) {
	user, err := s.adminRepository.FindAdminByID(user_Id)
	if err != nil {
		return false, err
	}

	return s.adminRepository.DeleteAdmin(user)
}

func (s *adminService) EmailExists(email string) bool {
	_, err := s.adminRepository.FindAdminByEmail(email)
	if err != nil {
		return false
	}
	return true
}
