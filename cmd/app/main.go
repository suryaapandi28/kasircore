// package main

// import (
// 	"github.com/suryaapandi28/kasircore/configs"
// 	"github.com/suryaapandi28/kasircore/internal/builder"
// 	"github.com/suryaapandi28/kasircore/internal/entity"
// 	"github.com/suryaapandi28/kasircore/pkg/cache"
// 	"github.com/suryaapandi28/kasircore/pkg/encrypt"
// 	"github.com/suryaapandi28/kasircore/pkg/postgres"
// 	"github.com/suryaapandi28/kasircore/pkg/server"
// 	"github.com/suryaapandi28/kasircore/pkg/token"
// )

// func main() {
// 	// Load configurations from .env file
// 	cfg, err := configs.NewConfig(".env")
// 	checkError(err)

// 	// Initialize PostgreSQL database connection
// 	db, err := postgres.InitPostgres(&cfg.Postgres)
// 	checkError(err)

// 	// Initialize Redis cache connection
// 	redisDB := cache.InitCache(&cfg.Redis)

// 	// Initialize encryption tool
// 	encryptTool := encrypt.NewEncryptTool(cfg.Encrypt.SecretKey, cfg.Encrypt.IV)

// 	// Initialize JWT token use case
// 	tokenUseCase := token.NewTokenUseCase(cfg.JWT.SecretKey)

// 	// Convert configs.Config to *entity.Config
// 	entityCfg := convertToEntityConfig(cfg)

// 	// Build public and private routes
// 	publicRoutes := builder.BuildPublicRoutes(db, redisDB, tokenUseCase, encryptTool, entityCfg)
// 	privateRoutes := builder.BuildPrivateRoutes()

// 	// Initialize and run the server
// 	srv := server.NewServer("app", publicRoutes, privateRoutes, cfg.JWT.SecretKey)
// 	srv.Run()
// }

// func checkError(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }

// // Example function to convert configs.Config to *entity.Config
// func convertToEntityConfig(cfg *configs.Config) *entity.Config {
// 	return &entity.Config{
// 		SMTP: entity.SMTPConfig{
// 			Host:     cfg.SMTP.Host,
// 			Port:     cfg.SMTP.Port,
// 			Password: cfg.SMTP.Password,
// 		},
// 		// Add other fields as needed
// 	}
// }

package main

import (
	"time"

	"github.com/suryaapandi28/kasircore/configs"
	"github.com/suryaapandi28/kasircore/internal/builder"
	"github.com/suryaapandi28/kasircore/internal/entity"
	"github.com/suryaapandi28/kasircore/pkg/cache"
	"github.com/suryaapandi28/kasircore/pkg/encrypt"
	"github.com/suryaapandi28/kasircore/pkg/middleware" // <- import middleware
	"github.com/suryaapandi28/kasircore/pkg/postgres"
	"github.com/suryaapandi28/kasircore/pkg/server"
	"github.com/suryaapandi28/kasircore/pkg/token"
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

	// **Tambahkan middleware rate limiter sebelum server dijalankan**
	mw := middleware.NewRateLimiter(redisDB, 100, time.Minute) // max 100 request / menit per IP
	for _, route := range publicRoutes {
		route.Handler = mw(route.Handler)
	}

	for _, route := range privateRoutes {
		route.Handler = mw(route.Handler)
	}

	// Initialize and run the server
	srv := server.NewServer("app", publicRoutes, privateRoutes, cfg.JWT.SecretKey)
	srv.Run()
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// Example function to convert configs.Config to *entity.Config
func convertToEntityConfig(cfg *configs.Config) *entity.Config {
	return &entity.Config{
		SMTP: entity.SMTPConfig{
			Host:     cfg.SMTP.Host,
			Port:     cfg.SMTP.Port,
			Password: cfg.SMTP.Password,
		},
		// Add other fields as needed
	}
}
