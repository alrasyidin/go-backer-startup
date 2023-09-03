package main

import (
	defaultLog "log"
	"os"
	"time"

	campaignHandle "github.com/alrasyidin/bwa-backer-startup/handler/campaign"
	userHandle "github.com/alrasyidin/bwa-backer-startup/handler/user"
	"github.com/alrasyidin/bwa-backer-startup/middleware"
	"github.com/alrasyidin/bwa-backer-startup/pkg/tokenization"
	"github.com/alrasyidin/bwa-backer-startup/repository"
	"github.com/alrasyidin/bwa-backer-startup/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	ginMode := os.Getenv("GIN_MODE")

	if ginMode != gin.ReleaseMode {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	app := gin.Default()

	app.SetTrustedProxies(nil)
	app.Use(middleware.HTTPLoggerMiddleware("BWA Backer"))
	app.Use(gin.Recovery())

	dsn := "host=localhost user=root password=postgres dbname=bwabackerdb port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	dbLogger := logger.New(
		defaultLog.New(os.Stdout, "\r\n", defaultLog.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: dbLogger,
	})

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

	tokenGenerator := tokenization.NewJWTGenerator("initokeninitokeninitokeninitoken")

	userRepo := repository.NewUserRepo(db)
	userService := service.NewUserService(userRepo)
	userHandler := userHandle.NewUserHandler(userService, tokenGenerator)

	campaignRepo := repository.NewCampaignRepo(db)
	campaignService := service.NewCampaignService(campaignRepo)
	campaignHandler := campaignHandle.NewCampaignHandler(campaignService)

	// static assets avatar
	app.Static("/images", "./images")

	v1 := app.Group("/api/v1")
	{
		v1.POST("/users/register", userHandler.Register)
		v1.POST("/users/session", userHandler.Login)
		v1.POST("/users/email-check", userHandler.CheckEmailAvailability)
		v1.POST("/users/avatar", middleware.AuthMiddlware(userService, tokenGenerator), userHandler.UploadAvatar)

		// campaigns
		v1.GET("/campaigns", campaignHandler.GetCampaigns)
		v1.GET("/campaigns/:id", campaignHandler.GetCampaign)
		v1.POST("/campaigns", middleware.AuthMiddlware(userService, tokenGenerator), campaignHandler.CreateCampaign)
		v1.PUT("/campaigns/:id", middleware.AuthMiddlware(userService, tokenGenerator), campaignHandler.UpdateCampaign)
		v1.POST("/campaign-images", middleware.AuthMiddlware(userService, tokenGenerator), campaignHandler.UploadImage)
	}

	const PORT = ":8000"
	log.Info().Msg("API has started at " + PORT)

	app.Run(PORT)
}
