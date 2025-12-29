package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/suryaapandi28/kasircore/internal/entity"
	"github.com/suryaapandi28/kasircore/internal/http/binder"
)

type MerchantHandler struct {
	merchantRepo service.merchantService
}

func NewMerchantHandler(merchantRepo service.merchantService) *MerchantHandler {
	return &MerchantHandler{merchantRepo: merchantRepo}
}
func (h *MerchantHandler) CreateMerchant(c echo.Context) error {
	var req binder.CreateMerchantRequest

	// 🔹 Bind request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Payload tidak valid",
		})
	}

	// 🔹 (Opsional) Validasi
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	// 🔹 Mapping request → entity
	merchant := entity.Merchant{
		F_nama_merchant:   req.F_nama_merchant,
		F_jenis_usaha:     req.F_jenis_usaha,
		F_email_merchant:  req.F_email_merchant,
		F_phone_merchant:  req.F_phone_merchant,
		F_alamat_merchant: req.F_alamat_merchant,
		F_kota:            req.F_kota,
		F_provinsi:        req.F_provinsi,
		F_kode_pos:        req.F_kode_pos,

		// default POS setting
		F_currency:        "IDR",
		F_ppn_enabled:     false,
		F_ppn_persen:      11,
		F_status_merchant: true,
	}

	// 🔹 Simpan
	if err := h.merchantRepo.Create(&merchant); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Gagal membuat merchant",
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Merchant berhasil dibuat",
		"data": echo.Map{
			"id": merchant.F_kode_merchant,
		},
	})
}

func (h *MerchantHandler) GetMerchants(c echo.Context) error {
	data, err := h.merchantRepo.FindAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Gagal mengambil data merchant",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": data,
	})
}

func (h *MerchantHandler) GetMerchantByID(c echo.Context) error {
	id := c.Param("id")

	merchant, err := h.merchantRepo.FindByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "Merchant tidak ditemukan",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": merchant,
	})
}

func (h *MerchantHandler) UpdateMerchant(c echo.Context) error {
	id := c.Param("id")

	merchant, err := h.merchantRepo.FindByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "Merchant tidak ditemukan",
		})
	}

	if err := c.Bind(merchant); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Payload tidak valid",
		})
	}

	if err := h.merchantRepo.Update(merchant); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Gagal update merchant",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Merchant berhasil diupdate",
		"data":    merchant,
	})
}

func (h *MerchantHandler) DeleteMerchant(c echo.Context) error {
	id := c.Param("id")

	if err := h.merchantRepo.Delete(id); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Gagal menghapus merchant",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Merchant berhasil dihapus",
	})
}
