package middleware

import (
	"bambamload/logger"
	"bytes"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

const maxLogSize = 700 * 1024 // 700KB

type APILog struct {
	Method       string
	Path         string
	IP           string
	UserAgent    string
	StartTime    time.Time
	EndTime      time.Time
	Duration     float64
	StatusCode   int
	RequestBody  interface{}
	ResponseBody interface{}
}

var apiLogChan = make(chan APILog, 5000)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func StartBackgroundRequestLogger() {
	go func() {
		for log := range apiLogChan {
			logger.RequestLogger.WithFields(logrus.Fields{
				"method":        log.Method,
				"path":          log.Path,
				"ip":            log.IP,
				"user_agent":    log.UserAgent,
				"start_time":    log.StartTime,
				"end_time":      log.EndTime,
				"duration":      log.Duration,
				"status_code":   log.StatusCode,
				"request_body":  log.RequestBody,
				"response_body": log.ResponseBody,
			}).Info("API LOG")
		}
	}()
}

func APILogger() fiber.Handler { //nolint:typecheck
	return func(c *fiber.Ctx) error {
		startTime := time.Now()

		// --- Capture request body if JSON ---
		var reqLog interface{}
		reqContentType := strings.ToLower(c.Get("Content-Type"))
		if strings.Contains(reqContentType, "application/json") {
			reqCopy := append([]byte(nil), c.Request().Body()...)
			if len(reqCopy) > maxLogSize {
				reqLog = "request body too large (truncated)"
				reqCopy = reqCopy[:maxLogSize]
			} else {
				reqLog = prettyPrintJSON(reqCopy)
			}
		} else {
			reqLog = "non-JSON request body"
		}

		// Continue request
		err := c.Next()

		// --- Capture response if JSON ---
		var respLog interface{}
		respContentType := strings.ToLower(string(c.Response().Header.ContentType()))
		if strings.Contains(respContentType, "application/json") {
			respCopy := append([]byte(nil), c.Response().Body()...)
			if len(respCopy) > maxLogSize {
				respLog = "response body too large (truncated)"
				respCopy = respCopy[:maxLogSize]
			} else {
				respLog = prettyPrintJSON(respCopy)
			}
		} else {
			respLog = "non-JSON response body"
		}

		// --- Capture safe values before goroutine ---
		logData := APILog{
			Method:       c.Method(),
			Path:         c.Path(),
			IP:           c.IP(),
			UserAgent:    c.Get("User-Agent"),
			StartTime:    startTime,
			EndTime:      time.Now(),
			Duration:     time.Since(startTime).Seconds(),
			StatusCode:   c.Response().StatusCode(),
			RequestBody:  reqLog,
			ResponseBody: respLog,
		}

		// Async logging without touching c after request
		go func(data APILog) {
			apiLogChan <- data
		}(logData)

		return err
	}
}

func prettyPrintJSON(input []byte) map[string]interface{} {
	if len(input) > maxLogSize {
		input = input[:maxLogSize]
	}
	var out map[string]interface{}
	_ = json.Unmarshal(input, &out)
	return out
}
