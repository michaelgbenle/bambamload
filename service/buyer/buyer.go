package buyer

import (
	"bambamload/constant"
	"bambamload/enum"
	"bambamload/logger"
	"bambamload/models"
	"bambamload/service/email"
	"bambamload/service/postgresrepository"
	"bambamload/service/redisService"
	"bambamload/service/uploadService"
	"bambamload/utils"
	"errors"
	"strings"
)

type ServiceBuyer struct {
	RedisService       redisService.RedisService
	PostgresRepository *postgresrepository.PostgresRepository
	EmailService       email.Email
	UploadService      *uploadService.UploadService
}

func NewServiceBuyer(redisService redisService.RedisService, postgresRepository *postgresrepository.PostgresRepository, emailService email.Email, uploadService *uploadService.UploadService) *ServiceBuyer {
	return &ServiceBuyer{
		RedisService:       redisService,
		PostgresRepository: postgresRepository,
		EmailService:       emailService,
		UploadService:      uploadService,
	}
}

func (sb *ServiceBuyer) Register(req models.BuyerRegisterRequest) (string, error) {
	sEmail := strings.ToLower(strings.TrimSpace(req.Email))
	phoneNumber := utils.StandardiseMSISDN(strings.TrimSpace(req.PhoneNumber))

	if sb.PostgresRepository.UserExists(sEmail, constant.Email) {
		return "User with this email already exists", errors.New("email already exists")
	}

	if sb.PostgresRepository.UserExists(phoneNumber, constant.PhoneNumber) {
		return "User with this phone number already exists", errors.New("phone number already exists")
	}

	password, err := utils.HashPassword(req.Password)
	if err != nil {
		logger.Logger.Errorf("unable to hash password: %v", err)
		return "unable to register, please try again later", err
	}

	ref := utils.GenerateReference("")
	buyer := &models.User{
		Name:         req.Name,
		Email:        sEmail,
		PhoneNumber:  phoneNumber,
		BusinessName: "",
		Status:       constant.Registering,
		Role:         enum.Buyer,
		Reference:    ref,
		Password:     password,
	}
	err = sb.PostgresRepository.CreateUser(buyer)
	if err != nil {
		logger.Logger.Errorf("Failed to create buyer: %v", err)
		return "unable to register, please try again later", err
	}
	return "success", nil
}
