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
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "there is an input error"))
	}

	if h.accountproviderService.EmailExists(input.Email) {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "email is already in use"))
	}

	newAdmin := entity.NewProviderAccount(input.Name, input.Email, input.Password, input.Role, input.Phone, input.Verification)
	admin, err := h.accountproviderService.CreateAdmin(newAdmin)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Successfully created a new admin, the email has been sent", admin))
}
