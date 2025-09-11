package main

import (
	"github.com/Kevinmajesta/depublic-backend/configs"
	"github.com/Kevinmajesta/depublic-backend/internal/builder"
	"github.com/Kevinmajesta/depublic-backend/internal/entity"
	"github.com/Kevinmajesta/depublic-backend/pkg/cache"
	"github.com/Kevinmajesta/depublic-backend/pkg/encrypt"
	"github.com/Kevinmajesta/depublic-backend/pkg/postgres"
	"github.com/Kevinmajesta/depublic-backend/pkg/server"
	"github.com/Kevinmajesta/depublic-backend/pkg/token"
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
	privateRoutes := builder.BuildPrivateRoutes(db, redisDB, encryptTool, entityCfg, tokenUseCase)

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
