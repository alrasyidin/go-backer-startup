package main

import (
	defaultLog "log"
	"os"
	"time"

	campaignHandle "github.com/alrasyidin/bwa-backer-startup/handler/campaign"
	transactionHandle "github.com/alrasyidin/bwa-backer-startup/handler/transaction"
	userHandle "github.com/alrasyidin/bwa-backer-startup/handler/user"
	"github.com/alrasyidin/bwa-backer-startup/middleware"
	"github.com/alrasyidin/bwa-backer-startup/pkg/tokenization"
	"github.com/alrasyidin/bwa-backer-startup/repository"
	"github.com/alrasyidin/bwa-backer-startup/service"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go"
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
	app.Use(cors.Default())
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

	// db = db.Set("gorm:auto_preload", false)

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

	transactionRepo := repository.NewTransactionRepo(db)
	paymentService := service.NewPayment(&service.PaymentConfig{
		ServerKey: "SB-Mid-server-thLn0-Fjl5Nu9-eEEbfM_56n",
		EnvType:   midtrans.Sandbox,
	}, transactionRepo, campaignRepo)
	transactionService := service.NewTransactionService(transactionRepo, campaignRepo, paymentService)
	transactionHandler := transactionHandle.NewTransactionHandler(transactionService, paymentService)

	// static assets avatar
	app.Static("/images", "./images")

	v1 := app.Group("/api/v1")

	freeRouter := v1
	freeRouter.POST("/users/register", userHandler.Register)
	freeRouter.POST("/users/session", userHandler.Login)
	freeRouter.POST("/users/email-check", userHandler.CheckEmailAvailability)

	freeRouter.GET("/campaigns", campaignHandler.GetCampaigns)
	freeRouter.GET("/campaigns/:id", campaignHandler.GetCampaign)

	freeRouter.POST("/transactions/notification", transactionHandler.ProcessPaymentNotification)

	requiredRouter := v1.Use(middleware.AuthMiddlware(userService, tokenGenerator))

	requiredRouter.POST("/users/avatar", userHandler.UploadAvatar)

	requiredRouter.POST("/campaigns", campaignHandler.CreateCampaign)
	requiredRouter.PUT("/campaigns/:id", campaignHandler.UpdateCampaign)
	requiredRouter.POST("/campaign-images", campaignHandler.UploadImage)
	requiredRouter.GET("/campaigns/:id/transactions", transactionHandler.GetCampaignTransactions)

	requiredRouter.GET("/transactions", transactionHandler.GetUserTransactions)
	requiredRouter.POST("/transactions", transactionHandler.CreateTransaction)

	const PORT = ":8000"
	log.Info().Msg("API has started at " + PORT)

	app.Run(PORT)
}
