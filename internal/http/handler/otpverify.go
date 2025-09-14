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

func NewOtpHandler(s service.OtpService) *OtpHandler {
	return &OtpHandler{otpService: s}
}
func (h *OtpHandler) GenerateOtp(c echo.Context) error {
	var input binder.GenerateOtpRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(30, "there is an input error"))
	}

	otp, err := h.otpService.GenerateOtp(input.F_email_account, input.F_otp_via)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, otp)
}
