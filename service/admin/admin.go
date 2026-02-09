package admin

import (
	"bambamload/logger"
	"bambamload/service/email"
	"bambamload/service/postgresrepository"
	"bambamload/service/redisService"
	"errors"
)

type ServiceAdmin struct {
	RedisService       redisService.RedisService
	PostgresRepository *postgresrepository.PostgresRepository
	EmailService       email.Email
}

func NewServiceAdmin(redisService redisService.RedisService, postgresRepository *postgresrepository.PostgresRepository, emailService email.Email) *ServiceAdmin {
	return &ServiceAdmin{
		RedisService:       redisService,
		PostgresRepository: postgresRepository,
		EmailService:       emailService,
	}
}

func (sa *ServiceAdmin) DashboardCards() (any, error) {
	cards, err := sa.PostgresRepository.AdminDashboardCards()
	if err != nil {
		logger.Logger.Errorf("[DashboardCards]Failed to get dashboard cards: %v", err)
		return nil, errors.New("unable to get dashboard cards")
	}
	return cards, nil
}
