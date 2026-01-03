package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountUser struct {
	F_kd_account    uuid.UUID `json:"f_kd_account" gorm:"type:uuid;primaryKey"`
	F_nama_account  string    `json:"f_nama_account" gorm:"size:150;not null"`
	F_email_account string    `json:"f_email_account" gorm:"unique;not null"`
	F_password      string    `json:"f_password" gorm:"not null"`
	F_role_accout   string    `json:"F_role_accout" gorm:"type:enum('superadmin','admin','staff');default:'admin'"`
	F_phone_account string    `json:"f_phone_account"`
	Auditable
	F_verification_account bool      `json:"f_verification_account"`
	F_jwt_token            string    `json:"f_jwt_token,omitempty"`
	F_jwt_token_expired    time.Time `json:"f_jwt_token_expired,omitempty"`
}

func (p *AccountUser) BeforeCreate(tx *gorm.DB) (err error) {
	if p.F_kd_account == uuid.Nil {
		p.F_kd_account = uuid.New()
	}
	if p.F_role_accout == "" {
		p.F_role_accout = "admin"
	}
	return
}

func NewAccountUser(f_nama_account, f_email_account, f_password, F_role_accout, f_phone_account string, f_verification_account bool) *AccountUser {
	return &AccountUser{
		F_kd_account:           uuid.New(),
		F_nama_account:         f_nama_account,
		F_email_account:        f_email_account,
		F_password:             f_password,
		F_role_accout:          F_role_accout,
		F_phone_account:        f_phone_account,
		Auditable:              NewAuditable(),
		F_verification_account: f_verification_account,
	}
}

func UpdateAccountUser(f_kd_account uuid.UUID, f_nama_account, f_email_account, f_password, F_role_accout, f_phone_account string, F_verification_account bool) *AccountUser {
	return &AccountUser{
		F_kd_account:           f_kd_account,
		F_nama_account:         f_nama_account,
		F_email_account:        f_email_account,
		F_password:             f_password,
		F_role_accout:          F_role_accout,
		F_phone_account:        f_phone_account,
		Auditable:              UpdateAuditable(),
		F_verification_account: F_verification_account,
	}
}

func (AccountUser) TableName() string {
	return "accounts_providers"
}
