package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type LoggerData struct {
	AppName    string
	Protocol   string
	Method     string
	Path       string
	StatusCode int
	StatusText string
	Duration   time.Duration
	MsgStr     string
}

func HTTPLoggerMiddleware(appName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()
		duration := time.Since(startTime)

		lData := &LoggerData{
			AppName:    appName,
			Protocol:   "http",
			Method:     c.Request.Method,
			Path:       c.Request.RequestURI,
			StatusCode: c.Writer.Status(),
			StatusText: http.StatusText(c.Writer.Status()),
			Duration:   duration,
			MsgStr:     "HTTP request",
		}

		logSwitch(lData)
	}
}

func logSwitch(data *LoggerData) {
	switch {
	case data.StatusCode >= 400 && data.StatusCode < 500:
		{
			log.Warn().Str("app_name", data.AppName).Str("method", data.Method).Str("path", data.Path).Dur("duration", data.Duration).Int("status_code", data.StatusCode).Str("status_text", data.StatusText).Msg(data.MsgStr)
		}
	case data.StatusCode >= 500:
		{
			log.Error().Str("app_name", data.AppName).Str("method", data.Method).Str("path", data.Path).Dur("duration", data.Duration).Int("status_code", data.StatusCode).Str("status_text", data.StatusText).Msg(data.MsgStr)
		}
	default:
		log.Info().Str("app_name", data.AppName).Str("method", data.Method).Str("path", data.Path).Dur("duration", data.Duration).Int("status_code", data.StatusCode).Str("status_text", data.StatusText).Msg(data.MsgStr)
	}
}
