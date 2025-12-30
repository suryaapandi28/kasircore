package service

import (
	"errors"
	"strings"

	"github.com/suryaapandi28/kasircore/internal/entity"
	"github.com/suryaapandi28/kasircore/internal/repository"
)

type MerchantService interface {
	CreateMerchant(merchant *entity.Merchant) (*entity.Merchant, error)
}
type merchantService struct {
	merchantRepository repository.MerchantRepository
}

func NewMerchantService(merchantRepository repository.MerchantRepository) *merchantService {
	return &merchantService{
		merchantRepository: merchantRepository,
	}
}

func (s *merchantService) CreateMerchant(merchant *entity.Merchant) (*entity.Merchant, error) {

	// ===== VALIDASI BUSINESS =====
	if strings.TrimSpace(merchant.F_nama_merchant) == "" {
		return nil, errors.New("nama merchant wajib diisi")
	}

	if strings.TrimSpace(merchant.F_jenis_usaha) == "" {
		return nil, errors.New("jenis usaha wajib diisi")
	}

	// Email opsional, tapi kalau ada harus valid (contoh sederhana)
	if merchant.F_email_merchant != "" && !strings.Contains(merchant.F_email_merchant, "@") {
		return nil, errors.New("format email tidak valid")
	}

	// ===== DEFAULT VALUE (POS SETTING) =====
	if merchant.F_currency == "" {
		merchant.F_currency = "IDR"
	}

	if merchant.F_ppn_persen == 0 {
		merchant.F_ppn_persen = 11
	}

	// ===== SIMPAN KE DB =====
	// newMerchant, err := s.merchantRepository.CreateMerchant(merchant)
	newMerchant, err := s.merchantRepository.CreateMerchant(merchant)
	if err != nil {
		return nil, err
	}

	return newMerchant, nil
}
