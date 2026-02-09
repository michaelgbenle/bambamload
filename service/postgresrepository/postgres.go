package postgresrepository

import (
	"bambamload/constant"
	"bambamload/enum"
	"bambamload/logger"
	"bambamload/models"
	"bambamload/utils"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository() *PostgresRepository {
	db, err := gorm.Open(postgres.Open(os.Getenv("POSTGRES_DSN")), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})

	if err != nil {
		logger.Logger.Fatalf("failed to connect to database: %v", err)
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Logger.Fatalf("failed to get generic database object: %v", err)
		return nil
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(15)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	p := &PostgresRepository{db: db}

	err = p.Migrate()
	if err != nil {
		logger.Logger.Fatalf("failed to migrate database schemas: %v", err)
	}
	logger.Logger.Infof("successfully migrated database schemas")

	//create superadmin
	err = p.CreateSuperAdmin()
	if err != nil {
		logger.Logger.Fatalf("failed to create super admin: %v", err)
	}

	return p
}

func (p *PostgresRepository) CreateSuperAdmin() error {
	exists := p.UserExists(os.Getenv("SUPERADMIN_EMAIL"), constant.Email)
	if exists {
		return nil
	}
	exists = p.UserExists(os.Getenv("SUPERADMIN_PHONE"), constant.PhoneNumber)
	if exists {
		return nil
	}

	pass, _ := utils.HashPassword(os.Getenv("SUPERADMIN_PASSWORD"))

	p.db.Create(&models.User{
		Name:          os.Getenv("SUPERADMIN_NAME"),
		Email:         os.Getenv("SUPERADMIN_EMAIL"),
		PhoneNumber:   os.Getenv("SUPERADMIN_PHONE"),
		Reference:     utils.GenerateReference(""),
		Password:      pass,
		Status:        constant.Approved,
		Role:          enum.SuperAdmin,
		IsBlocked:     false,
		IsActive:      true,
		LastLoginTime: time.Time{},
	})
	return nil
}

func (p *PostgresRepository) Close() error {
	db, err := p.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (p *PostgresRepository) Migrate() error {
	return p.db.AutoMigrate(&models.User{}, &models.Product{}, &models.ProductUpload{}, &models.Order{}, &models.Cart{})
}

func (p *PostgresRepository) Ping() error {

	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Ping()
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresRepository) DB() *gorm.DB {
	return p.db
}
