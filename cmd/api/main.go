package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alrasyidin/bwa-backer-startup/config"
	"github.com/alrasyidin/bwa-backer-startup/db"
	_ "github.com/alrasyidin/bwa-backer-startup/docs"
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
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func signalExit(log zerolog.Logger) {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	s := <-quit
	log.Info().Msgf("caught signal: %v", map[string]string{"signal": s.String()})
	os.Exit(0)
}

// @title           BWA Backer Startup
// @version         2.0
// @description     This is a bwa backer api
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8000
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	ginMode := os.Getenv("GIN_MODE")

	if ginMode != gin.ReleaseMode {
		customlog.Logger = customlog.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	config, err := config.LoadConfig(".")
	if err != nil {
		customlog.Fatal().Msgf("failed map config, please provide env with right value: %v", err)
	}
	app := gin.Default()

	app.SetTrustedProxies(nil)
	app.Use(middleware.HTTPLoggerMiddleware("BWA Backer"))
	app.Use(cors.Default())
	app.Use(gin.Recovery())

	dbConn := db.ConnDB(config)
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
	userHandler := userHandle.NewUserHandler(userService, tokenGenerator, config)

	campaignRepo := repository.NewCampaignRepo(dbConn)
	campaignService := service.NewCampaignService(campaignRepo)
	campaignHandler := campaignHandle.NewCampaignHandler(campaignService)

	transactionRepo := repository.NewTransactionRepo(dbConn)
	var midtransEnvTYpe midtrans.EnvironmentType
	if config.Environtment == "PRODUCTION" {
		midtransEnvTYpe = midtrans.Production
	} else {
		midtransEnvTYpe = midtrans.Sandbox
	}
	paymentService := service.NewPayment(&service.PaymentConfig{
		ServerKey: config.MidtransServerKey,
		EnvType:   midtransEnvTYpe,
	}, transactionRepo, campaignRepo)
	transactionService := service.NewTransactionService(transactionRepo, campaignRepo, paymentService)
	transactionHandler := transactionHandle.NewTransactionHandler(transactionService, paymentService)

	// static assets avatar
	app.Static("/images", "./images")

	v1 := app.Group("/api/v1")

	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

	srv := &http.Server{
		Addr:    PORT,
		Handler: app,
	}

	shutdownError := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit
		customlog.Info().Msgf("caught signal: %v", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	err = srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		customlog.Fatal().Msgf("listen: %s\n", err)
	}
	err = <-shutdownError

	if err != nil {
		customlog.Fatal().Msgf("listen: %s\n", err)
	}

	customlog.Info().Msgf("stopped server: http://localhost:%v", PORT)

	// app.Run(PORT)
}
