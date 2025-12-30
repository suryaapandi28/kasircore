package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Merchant struct {
	// PRIMARY KEY
	F_kode_merchant uuid.UUID `json:"f_kode_merchant" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`

	// IDENTITAS
	F_nama_merchant string `json:"f_nama_merchant" gorm:"size:150;not null"`
	F_jenis_usaha   string `json:"f_jenis_usaha" gorm:"size:50;not null"` // F&B, Retail, Jasa

	// KONTAK
	F_email_merchant string `json:"f_email_merchant" gorm:"size:150"`
	F_phone_merchant string `json:"f_phone_merchant" gorm:"size:30"`

	// LOKASI
	F_alamat_merchant string `json:"f_alamat_merchant" gorm:"type:text"`
	F_kota            string `json:"f_kota" gorm:"size:100"`
	F_provinsi        string `json:"f_provinsi" gorm:"size:100"`
	F_kode_pos        string `json:"f_kode_pos" gorm:"size:10"`

	// PENGATURAN POS
	F_currency    string  `json:"f_currency" gorm:"size:10;default:'IDR'"`
	F_ppn_enabled bool    `json:"f_ppn_enabled" gorm:"default:false"`
	F_ppn_persen  float64 `json:"f_ppn_persen" gorm:"type:numeric(5,2);default:11.00"`

	// STATUS
	F_status_merchant bool `json:"f_status_merchant" gorm:"default:true"`

	// AUDIT
	Auditable
}

func (m *Merchant) BeforeCreate(tx *gorm.DB) (err error) {
	if m.F_kode_merchant == uuid.Nil {
		m.F_kode_merchant = uuid.New()
	}
	return
}
func (Merchant) TableName() string {
	return "merchants"
}
func NewMerchant(nama, jenis, email, phone, alamat, kota, provinsi, kodePos string) *Merchant {
	return &Merchant{
		F_kode_merchant:   uuid.New(),
		F_nama_merchant:   nama,
		F_jenis_usaha:     jenis,
		F_email_merchant:  email,
		F_phone_merchant:  phone,
		F_alamat_merchant: alamat,
		F_kota:            kota,
		F_provinsi:        provinsi,
		F_kode_pos:        kodePos,
		F_currency:        "IDR",
		F_ppn_enabled:     false,
		F_ppn_persen:      11.00,
		F_status_merchant: true,
	}
}
