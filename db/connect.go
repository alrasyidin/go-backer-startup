package db

import (
	"log"
	"os"
	"time"

	customlog "github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnDB() *gorm.DB {
	dsn := "host=localhost user=root password=postgres dbname=bwabackerdb port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
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
		customlog.Fatal().Msgf("failed connect to db: %v", err)
	}

	// db = db.Set("gorm:auto_preload", false)

	customlog.Info().Msg("connected to db")

	return db
}
