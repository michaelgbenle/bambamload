package handler

import (
	"bambamload/constant"
	"bambamload/service/admin"
	"bambamload/service/buyer"
	"bambamload/service/email"
	"bambamload/service/postgresrepository"
	"bambamload/service/redisService"
	"bambamload/service/supplier"
	"bambamload/service/uploadService"
	"bambamload/service/utilities"
	"bambamload/utils"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	f "github.com/gofiber/fiber/v2"
)

type Handler struct {
	RedisService       redisService.RedisService
	PostgresRepository postgresrepository.PostgresRepository
	EmailService       email.Email
	UploadService      *uploadService.UploadService
	//SmsService         smsservice.SMSService
	AdminService     *admin.ServiceAdmin
	SupplierService  *supplier.ServiceSupplier
	BuyerService     *buyer.ServiceBuyer
	UtilitiesService *utilities.ServiceUtilities
}

func NewHandler(redisService redisService.RedisService,
	postgresRepository postgresrepository.PostgresRepository, emailService email.Email, uploadService *uploadService.UploadService, adminService *admin.ServiceAdmin,
	supplierService *supplier.ServiceSupplier, buyerService *buyer.ServiceBuyer, utilitiesService *utilities.ServiceUtilities) *Handler {
	return &Handler{
		RedisService:       redisService,
		PostgresRepository: postgresRepository,
		EmailService:       emailService,
		UploadService:      uploadService,
		AdminService:       adminService,
		SupplierService:    supplierService,
		BuyerService:       buyerService,
		UtilitiesService:   utilitiesService,
	}
}

func (h *Handler) Health(c *f.Ctx) error {

	return utils.WriteResponse(c, http.StatusOK, true, "server is healthy", &f.Map{
		"message": "Healthy",
	})
}

func (h *Handler) WelcomeHandler(c *f.Ctx) error {
	var (
		redisStatus, postgresStatus = constant.Healthy, constant.Healthy
	)
	redisError := h.RedisService.Ping()
	if redisError != nil {

		redisStatus = redisError.Error()
	}

	pgError := h.PostgresRepository.Ping()
	if pgError != nil {
		postgresStatus = pgError.Error()
	}

	return utils.WriteResponse(c, http.StatusOK, true, "Welcome to BAMBAMLOAD API", &f.Map{
		"app":             constant.Healthy,
		"redis_status":    redisStatus,
		"postgres_status": postgresStatus,
	})
}

func (h *Handler) NotFoundHandler(c *f.Ctx) error {

	return utils.WriteResponse(c, http.StatusNotFound, false, "route does not exist", &f.Map{
		"message": fmt.Sprintf("route '%v%v' does not exists", c.Hostname(), c.OriginalURL()),
	})
}

// GetLogs retrieves and returns logs based on the parameters provided in the request context.
func (h *Handler) GetLogs(c *f.Ctx) error {
	date := c.Params("date", time.Now().Format("2006-01-02"))

	if _, err := time.Parse("2006-01-02", date); err != nil {
		return c.Status(f.StatusBadRequest).JSON(f.Map{
			"error": "Invalid date format. Use YYYY-MM-DD",
		})
	}

	var filePath string
	switch c.Query(constant.Data) {
	case constant.Requests:
		filePath = filepath.Join(fmt.Sprintf("%s/messages-%s.log", constant.RequestLogsDir, date))
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return c.Status(http.StatusNotFound).JSON(f.Map{
				"error": "No logs found for the given date",
			})
		}

	default:
		filePath = filepath.Join(fmt.Sprintf("%s/messages-%s.log", constant.ErrorLogsDir, date))
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return c.Status(http.StatusNotFound).JSON(f.Map{
				"error": "No logs found for the given date",
			})
		}
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(f.Map{
			"error": "Error while reading file",
		})
	}

	c.Set("Content-Type", "text/plain; charset=utf-8")
	c.Set("Content-Disposition", "inline")
	return c.Status(http.StatusOK).Send(content)
}
