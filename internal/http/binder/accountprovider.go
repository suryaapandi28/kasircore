package binder

// Request untuk login provider
type ProviderLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Request untuk membuat akun provider baru
type ProviderCreateRequest struct {
	Name         string `json:"name" validate:"required"` // Nama penyedia, ex: APSM Indonesia Global
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required"`
	Role         string `json:"role" validate:"omitempty,oneof=superadmin admin staff"`
	Phone        string `json:"phone" validate:"required"`
	Verification bool   `json:"verification"`
}

// Request untuk update akun provider
type ProviderUpdateRequest struct {
	ID           string `param:"id" validate:"required,uuid"` // UUID provider
	Name         string `json:"name" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"omitempty"` // Password boleh kosong (tidak update)
	Role         string `json:"role" validate:"omitempty,oneof=superadmin admin staff"`
	Phone        string `json:"phone" validate:"required"`
	Verification bool   `json:"verification"`
}

// Request untuk delete akun provider
type ProviderDeleteRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}
