package main

import (
	"os"
	"time"

	handler "github.com/alrasyidin/bwa-backer-startup/handler/user"
	"github.com/alrasyidin/bwa-backer-startup/repository"
	"github.com/alrasyidin/bwa-backer-startup/service"
	ginzerolog "github.com/dn365/gin-zerolog"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	app := gin.Default()

	app.SetTrustedProxies(nil)
	app.Use(ginzerolog.Logger("gin"))

	dsn := "host=localhost user=root password=postgres dbname=bwabackerdb port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		log.Fatal().Msgf("failed connect to db: %v", err)
	}

	log.Info().Msg("connected to db")
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal().Msgf("failed get connection db instance: %v", err)
	}
	defer sqlDB.Close()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	userRepo := repository.NewUserRepo(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(*userService)
	v1 := app.Group("/api/v1")
	{
		v1.POST("/users/register", userHandler.Register)
		v1.POST("/users/session", userHandler.Login)
		v1.POST("/users/email-check", userHandler.CheckEmailAvailability)
	}

	const PORT = ":8000"
	log.Info().Msg("API has started at " + PORT)

	app.Run(PORT)
}
