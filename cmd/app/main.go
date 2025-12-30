package main

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/suryaapandi28/kasircore/configs"
	"github.com/suryaapandi28/kasircore/internal/builder"
	"github.com/suryaapandi28/kasircore/internal/entity"
	"github.com/suryaapandi28/kasircore/pkg/cache"
	"github.com/suryaapandi28/kasircore/pkg/encrypt"
	"github.com/suryaapandi28/kasircore/pkg/middleware"
	"github.com/suryaapandi28/kasircore/pkg/postgres"
	"github.com/suryaapandi28/kasircore/pkg/server"
	"github.com/suryaapandi28/kasircore/pkg/token"
	"github.com/suryaapandi28/kasircore/pkg/utils"
)

func main() {
	// Load configurations from .env file
	cfg, err := configs.NewConfig(".env")
	checkError(err)

	// Initialize PostgreSQL database connection
	db, err := postgres.InitPostgres(&cfg.Postgres)
	checkError(err)

	// Initialize Redis cache connection
	redisDB := cache.InitCache(&cfg.Redis)

	// Initialize encryption tool
	encryptTool := encrypt.NewEncryptTool(cfg.Encrypt.SecretKey, cfg.Encrypt.IV)

	// Initialize JWT token use case
	tokenUseCase := token.NewTokenUseCase(cfg.JWT.SecretKey)

	// Convert configs.Config to *entity.Config
	entityCfg := convertToEntityConfig(cfg)

	// Build public and private routes
	publicRoutes := builder.BuildPublicRoutes(db, redisDB, tokenUseCase, encryptTool, entityCfg)
	privateRoutes := builder.BuildPrivateRoutes()

	// Middleware rate limiter
	mw := middleware.NewRateLimiter(redisDB, 10, 30*time.Second)
	for _, route := range publicRoutes {
		route.Handler = mw(route.Handler)
	}
	for _, route := range privateRoutes {
		route.Handler = mw(route.Handler)
	}

	// Initialize server
	srv := server.NewServer("app", publicRoutes, privateRoutes, cfg.JWT.SecretKey)

	// Endpoint reset penalty
	srv.Echo.POST("/reset/:ip", func(c echo.Context) error {
		ip := c.Param("ip")
		err := utils.ResetPenalty(redisDB, ip)
		if err != nil {
			return c.JSON(500, map[string]string{"error": err.Error()})
		}
		return c.JSON(200, map[string]string{"message": "Penalty untuk " + ip + " sudah direset"})
	})

	// Endpoint cek IP real
	srv.Echo.GET("/myip", func(c echo.Context) error {
		ip := c.RealIP()
		return c.JSON(200, map[string]string{"your_ip": ip})
	})

	// Run server
	srv.Run()
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func convertToEntityConfig(cfg *configs.Config) *entity.Config {
	return &entity.Config{
		SMTP: entity.SMTPConfig{
			Host:     cfg.SMTP.Host,
			Port:     cfg.SMTP.Port,
			Password: cfg.SMTP.Password,
		},
	}
}
