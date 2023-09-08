package main

import (
	"os"
	"time"

	"github.com/alrasyidin/bwa-backer-startup/db"
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
	customlog "github.com/rs/zerolog/log"
)

func main() {
	ginMode := os.Getenv("GIN_MODE")

	if ginMode != gin.ReleaseMode {
		customlog.Logger = customlog.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	app := gin.Default()

	app.SetTrustedProxies(nil)
	app.Use(middleware.HTTPLoggerMiddleware("BWA Backer"))
	app.Use(cors.Default())
	app.Use(gin.Recovery())

	dbConn := db.ConnDB()
	sqlDB, err := dbConn.DB()
	if err != nil {
		customlog.Fatal().Msgf("failed get connection db instance: %v", err)
	}
	defer sqlDB.Close()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)
	tokenGenerator := tokenization.NewJWTGenerator("initokeninitokeninitokeninitoken")

	userRepo := repository.NewUserRepo(dbConn)
	userService := service.NewUserService(userRepo)
	userHandler := userHandle.NewUserHandler(userService, tokenGenerator)

	campaignRepo := repository.NewCampaignRepo(dbConn)
	campaignService := service.NewCampaignService(campaignRepo)
	campaignHandler := campaignHandle.NewCampaignHandler(campaignService)

	transactionRepo := repository.NewTransactionRepo(dbConn)
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

	freeRouter.POST("/transactions/notification", middleware.BeginDBTransaction(dbConn), transactionHandler.ProcessPaymentNotification)

	requiredRouter := v1.Use(middleware.AuthMiddlware(userService, tokenGenerator))

	requiredRouter.POST("/users/avatar", userHandler.UploadAvatar)
	requiredRouter.GET("/users/me", userHandler.Me)

	requiredRouter.POST("/campaigns", campaignHandler.CreateCampaign)
	requiredRouter.PUT("/campaigns/:id", campaignHandler.UpdateCampaign)
	requiredRouter.POST("/campaign-images", campaignHandler.UploadImage)
	requiredRouter.GET("/campaigns/:id/transactions", transactionHandler.GetCampaignTransactions)

	requiredRouter.GET("/transactions", transactionHandler.GetUserTransactions)
	requiredRouter.POST("/transactions", transactionHandler.CreateTransaction)

	const PORT = ":8000"
	customlog.Info().Msg("API has started at " + PORT)

	app.Run(PORT)
}
