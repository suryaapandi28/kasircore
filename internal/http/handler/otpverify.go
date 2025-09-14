package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/suryaapandi28/kasircore/internal/http/binder"
	"github.com/suryaapandi28/kasircore/internal/service"
	"github.com/suryaapandi28/kasircore/pkg/response"
)

type OtpHandler struct {
	otpService service.OtpService
}

func NewOtpHandler(OtpService service.OtpService) OtpHandler {
	return OtpHandler{otpService: OtpService}
}

func (h *OtpHandler) GenerateOtp(c echo.Context) error {
	var input binder.GenerateOtpRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(30, "there is an input error"))
	}

	hasilotp, err := h.otpService.GenerateOtp(input.F_email_account, input.F_otp_via)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(30, err.Error()))
	}

	// return c.JSON(http.StatusOK, otp)
	respData := map[string]interface{}{
		"kode_account": hasilotp.F_kd_account,
		"kode_otp":     hasilotp.F_kode_otp,
		"expired_at":   hasilotp.F_otp_expired,
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(200, "Successfully created a new otp verify", respData))
}
