package entity

import (
	"time"

	"github.com/google/uuid"
)

type OtpVerify struct {
	F_kd_otp      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"f_kd_otp"`
	F_kd_account  uuid.UUID `gorm:"type:uuid;not null" json:"f_kd_account"`
	F_kode_otp    string    `gorm:"type:varchar(10);not null" json:"f_kode_otp"`
	F_otp_expired time.Time `gorm:"not null" json:"f_otp_expired"`
	F_otp_via     string    `gorm:"type:varchar(50);not null" json:"f_otp_via"`

	Auditable
}

func (OtpVerify) TableName() string {
	return "otp_verify"
}
