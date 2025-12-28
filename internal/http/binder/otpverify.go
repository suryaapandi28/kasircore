package binder

type GenerateOtpRequest struct {
	F_email_account string `json:"f_email_account" validate:"required"`
	F_otp_via       string `json:"f_otp_via" validate:"required"`
}

type VerifyOtpRequest struct {
	F_email_account string `json:"f_email_account" validate:"required"`
	F_kode_otp      string `json:"f_kode_otp" validate:"required"`
}
