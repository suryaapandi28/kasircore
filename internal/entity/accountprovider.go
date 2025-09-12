package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProviderAccount struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name     string    `json:"name" gorm:"size:150;not null"`
	Email    string    `json:"email" gorm:"unique;not null"`
	Password string    `json:"-" gorm:"not null"`
	Role     string    `json:"role" gorm:"type:enum('superadmin','admin','staff');default:'admin'"`
	Phone    string    `json:"phone"`
	Auditable
	Verification      bool      `json:"verification"`
	JwtToken          string    `json:"jwt_token,omitempty"`
	JwtTokenExpiresAt time.Time `json:"jwt_token_expires_at,omitempty"`
}

func (p *ProviderAccount) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	if p.Role == "" {
		p.Role = "admin"
	}
	return
}

func NewProviderAccount(name, email, password, role, phone string, verification bool) *ProviderAccount {
	return &ProviderAccount{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		Password:     password,
		Role:         role,
		Phone:        phone,
		Auditable:    NewAuditable(),
		Verification: verification,
	}
}

func UpdateProviderAccount(id uuid.UUID, name, email, password, role, phone string, verification bool) *ProviderAccount {
	return &ProviderAccount{
		ID:           id,
		Name:         name,
		Email:        email,
		Password:     password,
		Role:         role,
		Phone:        phone,
		Auditable:    UpdateAuditable(),
		Verification: verification,
	}
}

func (ProviderAccount) TableName() string {
	return "accounts_provider"
}
