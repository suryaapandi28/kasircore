package repository

import (
	"encoding/json"

	"github.com/Kevinmajesta/depublic-backend/internal/entity"
	"github.com/Kevinmajesta/depublic-backend/pkg/cache"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	CreatePayment(payment *entity.Payments) (*entity.Payments, error)
	FindPayByID(payment_id uuid.UUID) (*entity.Payments, error)
}

type paymentRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

func NewPaymentRepository(db *gorm.DB, cacheable cache.Cacheable) PaymentRepository {
	return &paymentRepository{db: db, cacheable: cacheable}
}

func (r *paymentRepository) CreatePayment(payment *entity.Payments) (*entity.Payments, error) {

	if err := r.db.Create(&payment).Error; err != nil {
		return payment, err
	}
	return payment, nil

}

func (r *paymentRepository) FindPayByID(payment_id uuid.UUID) (*entity.Payments, error) {
	var pay entity.Payments

	paysdata := &pay
	key := "FindPayByID"

	data, _ := r.cacheable.Get(key)

	if data == "" {

		err := r.db.Raw("SELECT * FROM payments WHERE payment_id = ?", payment_id).Scan(&pay).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, nil
			}
			return nil, err
		}
		return &pay, nil

	} else {
		// Data ditemukan di Redis, unmarshal data ke transaction
		err := json.Unmarshal([]byte(data), &paysdata)
		if err != nil {
			return paysdata, err
		}
	}
	return &pay, nil

}
