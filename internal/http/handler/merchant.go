package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/suryaapandi28/kasircore/internal/entity"
	"github.com/suryaapandi28/kasircore/internal/http/binder"
	"github.com/suryaapandi28/kasircore/internal/service"
	"github.com/suryaapandi28/kasircore/pkg/response"
)

type MerchantHandler struct {
	merchantService service.MerchantService
}

func NewMerchantHandler(merchantService service.MerchantService) MerchantHandler {
	return MerchantHandler{

		merchantService: merchantService,
	}
}

func (h *MerchantHandler) CreateMerchant(c echo.Context) error {
	input := binder.CreateMerchantRequest{}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(30, "there is an input error"))
	}

	NewMerchant := entity.NewMerchant(input.F_nama_merchant, input.F_jenis_usaha, input.F_email_merchant, input.F_phone_merchant, input.F_alamat_merchant, input.F_kota, input.F_provinsi, input.F_kode_pos)
	addMerchant, err := h.merchantService.CreateMerchant(NewMerchant)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// Balikin hanya field tertentu, bukan semuanya
	respData := map[string]interface{}{
		"nama_merchant":  addMerchant.F_nama_merchant,
		"email_merchant": addMerchant.F_email_merchant,
		"created_at":     addMerchant.CreatedAt,
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(200, "Successfully created a new account provider", respData))
}
