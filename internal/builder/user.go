package builder

import (
	"github.com/suryaapandi28/kasircore/internal/entity"
	"github.com/suryaapandi28/kasircore/internal/http/handler"
	"github.com/suryaapandi28/kasircore/internal/http/router"
	"github.com/suryaapandi28/kasircore/internal/repository"
	"github.com/suryaapandi28/kasircore/internal/service"
	"github.com/suryaapandi28/kasircore/pkg/cache"
	"github.com/suryaapandi28/kasircore/pkg/email"
	"github.com/suryaapandi28/kasircore/pkg/encrypt"
	"github.com/suryaapandi28/kasircore/pkg/route"
	"github.com/suryaapandi28/kasircore/pkg/token"

	// "github.com/labstack/echo/"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func BuildPublicRoutes(db *gorm.DB, redisDB *redis.Client, tokenUseCase token.TokenUseCase, encryptTool encrypt.EncryptTool,
	entityCfg *entity.Config) []*route.Route {
	cacheable := cache.NewCacheable(redisDB)
	emailService := email.NewEmailSender(entityCfg)
	userRepository := repository.NewUserRepository(db, nil)

	notificationRepository := repository.NewNotificationRepository(db, cacheable)
	notificationService := service.NewNotificationService(notificationRepository, tokenUseCase, userRepository)
	userService := service.NewUserService(userRepository, tokenUseCase, encryptTool, emailService, notificationService)
	notificationHandler := handler.NewNotificationHandler(notificationService, userService)

	userHandler := handler.NewUserHandler(userService)

	adminRepository := repository.NewAdminRepository(db, nil)
	adminService := service.NewAdminService(adminRepository, tokenUseCase, encryptTool, emailService, notificationService)
	adminHandler := handler.NewAdminHandler(adminService)

	//Event
	eventRepository := repository.NewEventRepository(db)
	eventService := service.NewEventService(eventRepository)
	eventHandler := handler.NewEventHandler(eventService)

	// Category
	categoryRepository := repository.NewCategoryRepository(db, cacheable)
	categoryService := service.NewCategoryService(categoryRepository)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	return router.PublicRoutes(userHandler, adminHandler, notificationHandler, eventHandler, categoryHandler)
}

func BuildPrivateRoutes(db *gorm.DB, redisDB *redis.Client, encryptTool encrypt.EncryptTool, entityCfg *entity.Config, tokenUseCase token.TokenUseCase) []*route.Route {
	cacheable := cache.NewCacheable(redisDB)
	emailService := email.NewEmailSender(entityCfg)
	userRepository := repository.NewUserRepository(db, cacheable)

	notificationRepository := repository.NewNotificationRepository(db, cacheable)
	notificationService := service.NewNotificationService(notificationRepository, nil, userRepository)
	userService := service.NewUserService(userRepository, nil, encryptTool, nil, notificationService)
	notificationHandler := handler.NewNotificationHandler(notificationService, userService)

	userHandler := handler.NewUserHandler(userService)

	adminRepository := repository.NewAdminRepository(db, cacheable)
	adminService := service.NewAdminService(adminRepository, nil, encryptTool, nil, notificationService)
	adminHandler := handler.NewAdminHandler(adminService)

	eventRepository := repository.NewEventRepository(db)
	eventService := service.NewEventService(eventRepository)
	eventHandler := handler.NewEventHandler(eventService)

	wishlistRepository := repository.NewWishlistRepository(db, cacheable)
	wishlistService := service.NewWishlistService(wishlistRepository, eventRepository, userRepository, notificationService)
	wishlistHandler := handler.NewWishlistHandler(wishlistService)

	cartRepository := repository.NewCartRepository(db, cacheable)
	cartService := service.NewCartService(cartRepository, eventRepository, userRepository, notificationService)
	cartHandler := handler.NewCartHandler(cartService)

	paymentRepository := repository.NewPaymentRepository(db, cacheable)
	paymentService := service.NewPaymentService(paymentRepository)

	transactionRepository := repository.NewTransactionRepository(db, cacheable)
	transactionService := service.NewTransactionService(transactionRepository, emailService, userRepository)
	transactionHandler := handler.NewTransactionHandler(transactionService, tokenUseCase, paymentService)

	categoryRepository := repository.NewCategoryRepository(db, cacheable)
	categoryService := service.NewCategoryService(categoryRepository)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	ticketRepository := repository.NewTicketRepository(db, cacheable)
	ticketService := service.NewTicketService(ticketRepository, nil)
	ticketHandler := handler.NewTicketHandler(ticketService)

	return router.PrivateRoutes(userHandler, adminHandler, transactionHandler,
		cartHandler, wishlistHandler, notificationHandler, eventHandler, categoryHandler, ticketHandler)
}
