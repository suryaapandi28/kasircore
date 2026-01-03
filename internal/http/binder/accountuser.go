package binder

// Request untuk login provider
type AccountUserLoginRequest struct {
	F_email_account string `json:"f_email_account" validate:"required,email"`
	F_password      string `json:"f_password" validate:"required"`
}

// Request untuk membuat akun provider baru
type AccountUserCreateRequest struct {
	F_nama_account         string `json:"f_nama_account" validate:"required"` // Nama penyedia, ex: APSM Indonesia Global
	F_email_account        string `json:"f_email_account" validate:"required,email"`
	F_password             string `json:"f_password" validate:"required"`
	F_role_accout          string `json:"f_role_accout" validate:"omitempty,oneof=superadmin admin staff"`
	F_phone_account        string `json:"f_phone_account" validate:"required"`
	F_verification_account bool   `json:"f_verification_account"`
}

// Request untuk update akun provider
type AccountUserUpdateRequest struct {
	ID           string `param:"id" validate:"required,uuid"` // UUID provider
	Name         string `json:"name" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"omitempty"` // Password boleh kosong (tidak update)
	Role         string `json:"role" validate:"omitempty,oneof=superadmin admin staff"`
	Phone        string `json:"phone" validate:"required"`
	Verification bool   `json:"verification"`
}

// Request untuk delete akun provider
type AccountUserDeleteRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}
