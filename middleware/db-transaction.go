package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	customlog "github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func StatusInList(status int, statusList []int) bool {
	for _, i := range statusList {
		if i == status {
			return true
		}
	}
	return false
}

func BeginDBTransaction(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tx := db.Begin()
		customlog.Info().Msg("beginning database transaction")

		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		c.Set("db_trx", tx)
		c.Next()

		if StatusInList(c.Writer.Status(), []int{http.StatusOK, http.StatusCreated}) {
			customlog.Info().Msg("committing transactions")
			if err := tx.Commit().Error; err != nil {
				customlog.Error().Msgf("trx commit error: %v", err)
			}
		} else {
			customlog.Error().Msgf("rolling back transaction due to status code: %v", c.Writer.Status())
			tx.Rollback()
		}
	}
}
