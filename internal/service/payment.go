package service

import (
	"github.com/Kevinmajesta/depublic-backend/internal/entity"
	"github.com/Kevinmajesta/depublic-backend/internal/repository"
	"github.com/google/uuid"
)

type PaymentService interface {
	CreatePayment(payment *entity.Payments) (*entity.Payments, error)
	FindPayByID(payment_id uuid.UUID) (*entity.Payments, error)
}

type paymentService struct {
	paymentRepository repository.PaymentRepository
}

func NewPaymentService(paymentRepo repository.PaymentRepository) PaymentService {
	return &paymentService{paymentRepository: paymentRepo}
}

func (s *paymentService) CreatePayment(payment *entity.Payments) (*entity.Payments, error) {
	return s.paymentRepository.CreatePayment(payment)
}
func (s *paymentService) FindPayByID(payment_id uuid.UUID) (*entity.Payments, error) {
	return s.paymentRepository.FindPayByID(payment_id)
}
