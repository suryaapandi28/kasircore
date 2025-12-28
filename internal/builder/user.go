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
	"github.com/suryaapandi28/kasircore/pkg/whatsapp"

	// "github.com/labstack/echo/"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func BuildPublicRoutes(db *gorm.DB, redisDB *redis.Client, tokenUseCase token.TokenUseCase, encryptTool encrypt.EncryptTool,
	entityCfg *entity.Config) []*route.Route {
	cacheable := cache.NewCacheable(redisDB)
	emailService := email.NewEmailSender(entityCfg)
	WaSender := whatsapp.NewWhatsappSender(entityCfg)

	accountproviderRepository := repository.NewAccountproviderRepository(db, cacheable)
	accountproviderService := service.NewAccountproviderService(accountproviderRepository, tokenUseCase, encryptTool, emailService)

	AccountproviderHandler := handler.NewAccountproviderHandler(accountproviderService)

	otpRepository := repository.NewOTPRepository(db, cacheable)
	otpService := service.NewOtpService(otpRepository, emailService, WaSender)
	otpHandler := handler.NewOtpHandler(otpService)

	return router.PublicRoutes(AccountproviderHandler, otpHandler)
}

func BuildPrivateRoutes() []*route.Route {

	return router.PrivateRoutes()
}
