package router

import (
	"net/http"

	"github.com/suryaapandi28/kasircore/internal/http/handler"
	"github.com/suryaapandi28/kasircore/pkg/route"
)

const (
	SuperAdmin = "superadmin"
	Admin      = "admin"
	Staff      = "staff"
)

var (
	// allRoles  = []string{Admin, User}
	onlySuperAdmin = []string{SuperAdmin}
	onlyAdmin      = []string{Admin}
	onlyStaff      = []string{Staff}
)

func PublicRoutes(AccountproviderHandler handler.AccountproviderHandler, otpHandler handler.OtpHandler, MerchantHandler handler.MerchantHandler, AccountUserHandler handler.AccountUserHandler) []*route.Route {
	return []*route.Route{
		{
			Method:  http.MethodPost,
			Path:    "/create-account-provider",
			Handler: AccountproviderHandler.CreateAdmin,
		},

		{
			Method:  http.MethodPost,
			Path:    "/create-account-user",
			Handler: AccountUserHandler.CreateAccountUser,
		},

		{
			Method:  http.MethodPost,
			Path:    "/login-user",
			Handler: AccountUserHandler.LoginUser,
		},

		{
			Method:  http.MethodPost,
			Path:    "/login-provider",
			Handler: AccountproviderHandler.LoginProvider,
		},

		{
			Method:  http.MethodPost,
			Path:    "/create-otp-verify",
			Handler: otpHandler.GenerateOtp,
		},

		{
			Method:  http.MethodPost,
			Path:    "/otp-verify",
			Handler: otpHandler.VerifyOtpRequest,
		},

		{
			Method:  http.MethodPost,
			Path:    "/create-otp-verify-user",
			Handler: otpHandler.GenerateOtp,
		},

		{
			Method:  http.MethodPost,
			Path:    "/otp-verify-user",
			Handler: otpHandler.VerifyOtpRequest,
		},

		// Merchant Routes can be added here

		{
			Method:  http.MethodPost,
			Path:    "/create-merchant",
			Handler: MerchantHandler.CreateMerchant,
		},
	}
}

func PrivateRoutes() []*route.Route {
	return []*route.Route{}
}
