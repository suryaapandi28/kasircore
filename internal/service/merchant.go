package service

import (
	"errors"

	"github.com/suryaapandi28/kasircore/internal/entity"

	"github.com/suryaapandi28/kasircore/internal/http/binder"

	"github.com/suryaapandi28/kasircore/internal/repository"
)

type merchantService struct {
	repo repository.MerchantRepository
}

func NewMerchantService(repo repository.MerchantRepository) *merchantService {
	return &merchantService{repo}
}

func (s *merchantService) CreateMerchant(req binder.CreateMerchantRequest) (*entity.Merchant, error) {

	merchant := &entity.Merchant{
		F_nama_merchant:   req.F_nama_merchant,
		F_jenis_usaha:     req.F_jenis_usaha,
		F_email_merchant:  req.F_email_merchant,
		F_phone_merchant:  req.F_phone_merchant,
		F_alamat_merchant: req.F_alamat_merchant,
		F_kota:            req.F_kota,
		F_provinsi:        req.F_provinsi,
		F_kode_pos:        req.F_kode_pos,

		F_currency:        "IDR",
		F_ppn_enabled:     false,
		F_ppn_persen:      11,
		F_status_merchant: true,
	}

	if err := s.repo.Create(merchant); err != nil {
		return nil, errors.New("gagal membuat merchant")
	}

	return merchant, nil
}

func (s *merchantService) GetAllMerchant() ([]entity.Merchant, error) {
	return s.repo.FindAll()
}

func (s *merchantService) GetMerchantByID(id string) (*entity.Merchant, error) {
	return s.repo.FindByID(id)
}
