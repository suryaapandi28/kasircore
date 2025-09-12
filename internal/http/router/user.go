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

func PublicRoutes(AccountproviderHandler handler.AccountproviderHandler) []*route.Route {
	return []*route.Route{
		{
			Method:  http.MethodPost,
			Path:    "/admins",
			Handler: AccountproviderHandler.CreateAdmin,
		},
	}
}

func PrivateRoutes() []*route.Route {
	return []*route.Route{}
}
