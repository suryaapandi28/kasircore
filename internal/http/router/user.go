package router

import (
	"net/http"

	"github.com/Kevinmajesta/depublic-backend/internal/http/handler"
	"github.com/Kevinmajesta/depublic-backend/pkg/route"
)

const (
	Admin = "admin"
	User  = "user"
)

var (
	allRoles  = []string{Admin, User}
	onlyAdmin = []string{Admin}
	onlyUser  = []string{User}
)

func PublicRoutes(userHandler handler.UserHandler,
	adminHandler handler.AdminHandler, notificationHandler handler.NotificationHandler, eventHandler handler.EventHandler, categoryHandler handler.CategoryHandler) []*route.Route {
	return []*route.Route{
		{
			Method:  http.MethodPost,
			Path:    "/login",
			Handler: userHandler.LoginUser,
		},
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: userHandler.CreateUser,
		},
		{
			Method:  http.MethodPost,
			Path:    "/login/admin",
			Handler: adminHandler.LoginAdmin,
		},
		{
			Method:  http.MethodPost,
			Path:    "/admins",
			Handler: adminHandler.CreateAdmin,
		},
		{
			Method:  http.MethodPost,
			Path:    "/password-reset-request",
			Handler: userHandler.RequestPasswordReset,
		},
		{
			Method:  http.MethodPost,
			Path:    "/verification-account",
			Handler: userHandler.VerifUser,
		},
		{
			Method:  http.MethodPost,
			Path:    "/password-reset",
			Handler: userHandler.ResetPassword,
		},
		//event
		{
			Method:  http.MethodPost,
			Path:    "/event",
			Handler: eventHandler.AddEvent,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodGet,
			Path:    "/event",
			Handler: eventHandler.GetAllEvent,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/event/",
			Handler: eventHandler.SearchEvents,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/event/filter/",
			Handler: eventHandler.FilterEvents,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/event/sort/",
			Handler: eventHandler.SortEvents,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/event/:id",
			Handler: eventHandler.GetEventByID,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPut,
			Path:    "/event/:id",
			Handler: eventHandler.UpdateEvent,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/event/:id",
			Handler: eventHandler.DeleteEventByID,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodGet,
			Path:    "/category",
			Handler: categoryHandler.GetAllCategory,
			Roles:   allRoles,
		},
		// By ID
		{
			Method:  http.MethodGet,
			Path:    "/category/:id",
			Handler: categoryHandler.GetCategoryByID,
			Roles:   onlyAdmin,
		},
		// By Param
		{
			Method:  http.MethodGet,
			Path:    "/category/",
			Handler: categoryHandler.GetCategoryByParam,
			Roles:   allRoles,
		},
		// TODO ROUTE POST
		{
			Method:  http.MethodPost,
			Path:    "/category",
			Handler: categoryHandler.AddCategory,
			Roles:   onlyAdmin,
		},
		// TODO ROUTE PUT
		{
			Method:  http.MethodPut,
			Path:    "/category/:id",
			Handler: categoryHandler.UpdateCategoryByID,
			Roles:   onlyAdmin,
		},
		// TODO ROUTE DELETE
		{
			Method:  http.MethodDelete,
			Path:    "/category/:id",
			Handler: categoryHandler.DeleteCategoryByID,
			Roles:   onlyAdmin,
		},
	}
}

func PrivateRoutes(userHandler handler.UserHandler,
	adminHandler handler.AdminHandler,
	transactionHandler handler.TransactionHandler, cartHandler handler.CartHandler,
	wishlistHandler handler.WishlistHandler,
	notificationHandler handler.NotificationHandler,
	eventHandler handler.EventHandler, categoryHandler handler.CategoryHandler,
	ticketHandler handler.TicketHandler) []*route.Route {
	return []*route.Route{
		//user
		{
			Method:  http.MethodPut,
			Path:    "/users/:user_id",
			Handler: userHandler.UpdateUser,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/users/:user_id",
			Handler: userHandler.DeleteUser,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/:user_id",
			Handler: userHandler.GetUserProfile,
			Roles:   allRoles,
		},
		//admin
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: adminHandler.FindAllUser,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodPut,
			Path:    "/admins/:user_id",
			Handler: adminHandler.UpdateAdmin,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/admins/:user_id",
			Handler: adminHandler.DeleteAdmin,
			Roles:   onlyAdmin,
		},
		//transaction
		{
			Method:  http.MethodPost,
			Path:    "transaction/create",
			Handler: transactionHandler.CreateTransaction,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "transaction/all",
			Handler: transactionHandler.FindAllTransaction,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPost,
			Path:    "transaction/check-pay",
			Handler: transactionHandler.CheckPayTransaction,
			Roles:   allRoles,
		},
		//wishlist
		{
			Method:  http.MethodGet,
			Path:    "/wishlist",
			Handler: wishlistHandler.GetAllWishlist,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodGet,
			Path:    "/wishlist/:id",
			Handler: wishlistHandler.GetWishlistByUserId,
			Roles:   onlyUser,
		},
		{
			Method:  http.MethodPost,
			Path:    "/wishlist/create",
			Handler: wishlistHandler.AddWishlist,
			Roles:   onlyUser,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/wishlist/remove",
			Handler: wishlistHandler.RemoveWishlist,
			Roles:   onlyUser,
		},
		//cart
		{
			Method:  http.MethodGet,
			Path:    "/cart",
			Handler: cartHandler.GetAllCarts,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodPost,
			Path:    "/cart",
			Handler: cartHandler.AddToCart,
			Roles:   onlyUser,
		},
		{
			Method:  http.MethodPost,
			Path:    "/cart/less",
			Handler: cartHandler.UpdateQuantityLess,
			Roles:   onlyUser,
		},
		{
			Method:  http.MethodPost,
			Path:    "/cart/add",
			Handler: cartHandler.UpdateQuantityAdd,
			Roles:   onlyUser,
		},
		{
			Method:  http.MethodGet,
			Path:    "/cart/:id",
			Handler: cartHandler.GetCartByUserId,
			Roles:   onlyUser,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/cart/:id",
			Handler: cartHandler.RemoveCart,
		},
		//notification
		{
			Method:  http.MethodPost,
			Path:    "/notification",
			Handler: notificationHandler.CreateNotification,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodGet,
			Path:    "/user/notification",
			Handler: notificationHandler.GetUserNotifications,
			Roles:   allRoles,
		},

		// TODO EVENT
		// Create
		{
			Method:  http.MethodPost,
			Path:    "/event",
			Handler: eventHandler.AddEvent,
			Roles:   onlyAdmin,
		},
		// Get All
		{
			Method:  http.MethodGet,
			Path:    "/event",
			Handler: eventHandler.GetAllEvent,
			Roles:   allRoles,
		},
		// Search
		{
			Method:  http.MethodGet,
			Path:    "/event/",
			Handler: eventHandler.SearchEvents,
			Roles:   allRoles,
		},
		// Filter
		{
			Method:  http.MethodGet,
			Path:    "/event/filter/",
			Handler: eventHandler.FilterEvents,
			Roles:   allRoles,
		},
		// Sort
		{
			Method:  http.MethodGet,
			Path:    "/event/sort/",
			Handler: eventHandler.SortEvents,
			Roles:   allRoles,
		},
		// By ID
		{
			Method:  http.MethodGet,
			Path:    "/event/:id",
			Handler: eventHandler.GetEventByID,
			Roles:   allRoles,
		},
		// Update
		{
			Method:  http.MethodPut,
			Path:    "/event/:id",
			Handler: eventHandler.UpdateEvent,
			Roles:   onlyAdmin,
		},
		// Delete
		{
			Method:  http.MethodDelete,
			Path:    "/event/:id",
			Handler: eventHandler.DeleteEventByID,
			Roles:   onlyAdmin,
		},

		// TODO CATEGORY
		// TODO ROUTE GET
		{
			Method:  http.MethodGet,
			Path:    "/category",
			Handler: categoryHandler.GetAllCategory,
			Roles:   allRoles,
		},
		// By ID
		{
			Method:  http.MethodGet,
			Path:    "/category/:id",
			Handler: categoryHandler.GetCategoryByID,
			Roles:   allRoles,
		},
		// By Param
		{
			Method:  http.MethodGet,
			Path:    "/category/",
			Handler: categoryHandler.GetCategoryByParam,
			Roles:   allRoles,
		},
		// TODO ROUTE POST
		{
			Method:  http.MethodPost,
			Path:    "/category",
			Handler: categoryHandler.AddCategory,
			Roles:   onlyAdmin,
		},
		// TODO ROUTE PUT
		{
			Method:  http.MethodPut,
			Path:    "/category/:id",
			Handler: categoryHandler.UpdateCategoryByID,
			Roles:   onlyAdmin,
		},
		// TODO ROUTE DELETE
		{
			Method:  http.MethodDelete,
			Path:    "/category/:id",
			Handler: categoryHandler.DeleteCategoryByID,
			Roles:   onlyAdmin,
		},
		//notif
		{
			Method:  http.MethodPost,
			Path:    "/notification",
			Handler: notificationHandler.CreateNotification,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodGet,
			Path:    "/user/notification/all",
			Handler: notificationHandler.GetUserNotifications,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/user/notification",
			Handler: notificationHandler.GetUserNotificationNoRead,
			Roles:   allRoles,
		},
		//ticket
		{
			Method:  http.MethodGet,
			Path:    "/tickets/event/:eventID",
			Handler: ticketHandler.FindTicketsByEventID,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tickets/:QRCode",
			Handler: ticketHandler.FindTicketsByQRCode,
			Roles:   onlyAdmin,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tickets",
			Handler: ticketHandler.FindAllTicket,
			Roles:   onlyAdmin,
		},
	}
}
