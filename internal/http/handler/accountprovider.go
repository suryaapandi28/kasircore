package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/suryaapandi28/kasircore/internal/entity"
	"github.com/suryaapandi28/kasircore/internal/http/binder"
	"github.com/suryaapandi28/kasircore/internal/service"
	"github.com/suryaapandi28/kasircore/pkg/response"
)

type AccountproviderHandler struct {
	accountproviderService service.AccountproviderService
}

func NewAccountproviderHandler(accountproviderService service.AccountproviderService) AccountproviderHandler {
	return AccountproviderHandler{accountproviderService: accountproviderService}
}

func (h *AccountproviderHandler) CreateAdmin(c echo.Context) error {
	input := binder.ProviderCreateRequest{}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(0, "there is an input error"))
	}

	if h.accountproviderService.EmailExists(input.F_email_account) {
		return c.JSON(http.StatusBadRequest, response.DuplicateEmailResponse(30, "email is already in use"))
	}
	if input.F_email_account == "" || input.F_password == "" {
		return c.JSON(http.StatusBadRequest,
			response.ErrorResponse(30, "Field wajib tidak boleh kosong"))
	}

	newAccountProvider := entity.NewProviderAccount(input.F_nama_account, input.F_email_account, input.F_password, input.F_role_accout, input.F_phone_account, input.F_verification_account)
	addAccountProvider, err := h.accountproviderService.CreateAdmin(newAccountProvider)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// Balikin hanya field tertentu, bukan semuanya
	respData := map[string]interface{}{
		"nama_account_provider":  addAccountProvider.F_nama_account,
		"email_account_provider": addAccountProvider.F_email_account,
		"created_at":             addAccountProvider.CreatedAt,
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(200, "Successfully created a new account provider", respData))
}
