package builder

import (
	"github.com/suryaapandi28/kasircore/internal/http/handler"
	"github.com/suryaapandi28/kasircore/internal/http/router"
	"github.com/suryaapandi28/kasircore/internal/repository"
	"github.com/suryaapandi28/kasircore/internal/service"
	"github.com/suryaapandi28/kasircore/pkg/cache"
	"github.com/suryaapandi28/kasircore/pkg/encrypt"
	"github.com/suryaapandi28/kasircore/pkg/route"
	"github.com/suryaapandi28/kasircore/pkg/token"

	// "github.com/labstack/echo/"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func BuildPublicRoutes(db *gorm.DB, redisDB *redis.Client, tokenUseCase token.TokenUseCase, encryptTool encrypt.EncryptTool) []*route.Route {
	cacheable := cache.NewCacheable(redisDB)

	accountproviderRepository := repository.NewAccountproviderRepository(db, cacheable)
	accountproviderService := service.NewAccountproviderService(accountproviderRepository, tokenUseCase, encryptTool)
	AccountproviderHandler := handler.NewAccountproviderHandler(accountproviderService)

	return router.PublicRoutes(AccountproviderHandler)
}

func BuildPrivateRoutes() []*route.Route {

	return router.PrivateRoutes()
}
