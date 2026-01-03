package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/suryaapandi28/kasircore/internal/entity"
	"github.com/suryaapandi28/kasircore/internal/http/binder"
	"github.com/suryaapandi28/kasircore/internal/service"
	"github.com/suryaapandi28/kasircore/pkg/response"
)

type AccountUserHandler struct {
	accountuserService service.AccountUserService
}

func NewAccountUserHandler(accountuserService service.AccountUserService) AccountUserHandler {
	return AccountUserHandler{accountuserService: accountuserService}
}

func (h *AccountUserHandler) CreateAccountUser(c echo.Context) error {
	input := binder.AccountUserCreateRequest{}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(30, "there is an input error"))
	}

	if h.accountuserService.EmailExists(input.F_email_account) {
		return c.JSON(http.StatusBadRequest, response.DuplicateEmailResponse(30, "email is already in use"))
	}
	if input.F_email_account == "" || input.F_password == "" {
		return c.JSON(http.StatusBadRequest,
			response.ErrorResponse(26, "Field wajib tidak boleh kosong"))
	}
	roleaccount := "user"
	verificationaccount := false

	newacountuser := entity.NewAccountUser(input.F_nama_account, input.F_email_account, input.F_password, roleaccount, input.F_phone_account, verificationaccount)
	accountuseradd, err := h.accountuserService.CreateAccountUser(newacountuser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// Balikin hanya field tertentu, bukan semuanya
	respData := map[string]interface{}{
		"nama_account_user":  accountuseradd.F_nama_account,
		"email_account_user": accountuseradd.F_email_account,
		"created_at":         accountuseradd.CreatedAt,
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(200, "Successfully created a new account user", respData))
}

func (h *AccountUserHandler) LoginUser(c echo.Context) error {
	input := binder.AccountUserLoginRequest{}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(30, "there is an input error"))
	}

	loginData, err := h.accountuserService.LoginUser(input.F_email_account, input.F_password)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(30, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "login success", map[string]interface{}{
		"email":      loginData.F_email_account,
		"roles":      loginData.F_role_accout,
		"token":      loginData.F_jwt_token,
		"expired_at": loginData.F_jwt_token_expired,
	}))

}
