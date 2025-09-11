package repository

import (
	"encoding/json"
	"time"

	"github.com/Kevinmajesta/depublic-backend/internal/entity"
	"github.com/Kevinmajesta/depublic-backend/pkg/cache"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketRepository interface {
	FindAllTicket() ([]entity.Tickets, error)
	FindTicketsByEventID(eventID uuid.UUID) ([]entity.Tickets, error)
	FindTicketsByQRCode(QRCode uuid.UUID) ([]entity.Tickets, error)
	CheckTicketExists(id uuid.UUID) (bool, error)
	CheckTicketCodeQRExists(id uuid.UUID) (bool, error)
}

type ticketRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

func NewTicketRepository(db *gorm.DB, cacheable cache.Cacheable) *ticketRepository {
	return &ticketRepository{db: db, cacheable: cacheable}
}

func (r *ticketRepository) CheckTicketExists(id uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.Model(&entity.Tickets{}).Where("event_id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ticketRepository) CheckTicketCodeQRExists(id uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.Model(&entity.Tickets{}).Where("code_qr = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ticketRepository) FindAllTicket() ([]entity.Tickets, error) {
	var tickets []entity.Tickets

	key := "FindAllTicket"

	data, err := r.cacheable.Get(key)
	if err == nil && data != "" {
		err = json.Unmarshal([]byte(data), &tickets)
		if err == nil {
			return tickets, nil
		}
	}
	result := r.db.Find(&tickets)
	if result.Error != nil {
		return nil, result.Error
	}

	marshalledTickets, err := json.Marshal(tickets)
	if err == nil {
		err = r.cacheable.Set(key, marshalledTickets, 5*time.Minute)
	}

	return tickets, err
}

func (r *ticketRepository) FindTicketsByEventID(eventID uuid.UUID) ([]entity.Tickets, error) {
	// Check if event with given ID exists first
	var tickets []entity.Tickets
	if err := r.db.Where("event_id = ?", eventID).Find(&tickets).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return tickets, nil
}

func (r *ticketRepository) FindTicketsByQRCode(QRCode uuid.UUID) ([]entity.Tickets, error) {
	var tickets []entity.Tickets
	if err := r.db.Where("code_qr = ?", QRCode).Find(&tickets).Error; err != nil {
		return nil, err
	}
	return tickets, nil
}
