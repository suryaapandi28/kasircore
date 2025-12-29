package binder

type CreateMerchantRequest struct {
	F_nama_merchant   string `json:"f_nama_merchant" form:"f_nama_merchant" validate:"required"`
	F_jenis_usaha     string `json:"f_jenis_usaha" form:"f_jenis_usaha" validate:"required"`
	F_email_merchant  string `json:"f_email_merchant" form:"f_email_merchant" validate:"omitempty,email"`
	F_phone_merchant  string `json:"f_phone_merchant" form:"f_phone_merchant"`
	F_alamat_merchant string `json:"f_alamat_merchant" form:"f_alamat_merchant"`
	F_kota            string `json:"f_kota" form:"f_kota"`
	F_provinsi        string `json:"f_provinsi" form:"f_provinsi"`
	F_kode_pos        string `json:"f_kode_pos" form:"f_kode_pos"`
	F_status          bool   `json:"f_status" form:"f_status"`
}
