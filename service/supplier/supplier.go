package supplier

import (
	"bambamload/constant"
	"bambamload/enum"
	"bambamload/logger"
	"bambamload/models"
	"bambamload/service/email"
	"bambamload/service/postgresrepository"
	"bambamload/service/redisService"
	"bambamload/service/uploadService"
	"errors"

	"gorm.io/gorm"
)

type ServiceSupplier struct {
	RedisService       redisService.RedisService
	PostgresRepository *postgresrepository.PostgresRepository
	EmailService       email.Email
	UploadService      *uploadService.UploadService
}

func NewServiceSupplier(redisService redisService.RedisService, postgresRepository *postgresrepository.PostgresRepository, emailService email.Email, uploadService *uploadService.UploadService) *ServiceSupplier {
	return &ServiceSupplier{
		RedisService:       redisService,
		PostgresRepository: postgresRepository,
		EmailService:       emailService,
		UploadService:      uploadService,
	}
}

//invited - registering - pending -verified

func (ss *ServiceSupplier) Register(req *models.SupplierRegisterRequest) (string, string, error) {

	user, err := ss.PostgresRepository.GetUser(req.Reference, constant.Reference)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "supplier does not exist", "", errors.New("user not found")
		}
		logger.Logger.Errorf("[Supplier] Register: get user error: %v", err)
		return "unable to register, please try again later", "", err
	}

	if user.Role != enum.Supplier {
		return "only supplier can register via this route", "", errors.New("only supplier can register")
	}

	if user.Status != constant.Invited {
		return "supplier is not invited", "", errors.New("supplier is not invited")
	}

	updateMap := make(map[string]interface{})
	updateMap["password"] = req.Password
	updateMap["status"] = constant.Registering
	err = ss.PostgresRepository.UpdateUser(user.ID, constant.ID, updateMap)
	if err != nil {
		logger.Logger.Errorf("[Supplier] Register: update user error: %v", err)
		return "unable to register user, please try again later", "", err
	}

	return "success", user.Email, nil
}

// SubmitBusinessProfile ...
func (ss *ServiceSupplier) SubmitBusinessProfile(req *models.SubmitBusinessProfileRequest, user *models.User) error {
	updateMap := make(map[string]interface{})

	if req.AccountType != "" {
		updateMap["account_type"] = req.AccountType
	}
	if req.BusinessDescription != "" {
		updateMap["business_description"] = req.BusinessDescription
	}
	if req.YearFounded != "" {
		updateMap["year_founded"] = req.YearFounded
	}
	if req.WebsiteUrl != "" {
		updateMap["website_url"] = req.WebsiteUrl
	}
	if req.LinkedInProfile != "" {
		updateMap["linked_in_profile"] = req.LinkedInProfile
	}
	if req.Country != "" {
		updateMap["country"] = req.Country
	}
	if req.State != "" {
		updateMap["state"] = req.State
	}
	if req.Address != "" {
		updateMap["address"] = req.Address
	}
	if req.RegionsServed != "" {
		updateMap["regions_served"] = req.RegionsServed
	}

	updateMap["kyc_status"] = constant.IdentityVerification

	err := ss.PostgresRepository.UpdateUser(user.ID, constant.ID, updateMap)
	if err != nil {
		logger.Logger.Errorf("[Supplier] SubmitBusinessProfile: update user error: %v", err)
		return errors.New("unable to submit business profile, please try again later")
	}

	return nil
}
