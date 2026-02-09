package postgresrepository

import (
	"bambamload/constant"
	"bambamload/enum"
	"bambamload/logger"
	"bambamload/models"
	"errors"

	"gorm.io/gorm/clause"
)

// CreateUser
func (p *PostgresRepository) CreateUser(req *models.User) error {
	return p.db.Omit("bvn", "nin").Create(req).Error
}

// GetUser fetches a user by any identifier provided
func (p *PostgresRepository) GetUser(id, identifier string) (*models.User, error) {
	var (
		user *models.User
		err  error
	)

	switch identifier {
	case constant.ID:
		err = p.db.Preload(clause.Associations).Where("id = ?", id).First(&user).Error

	case constant.PhoneNumber:
		err = p.db.Preload(clause.Associations).Where("phone_number = ?", id).First(&user).Error

	case constant.Email:
		err = p.db.Preload(clause.Associations).Where("email = ?", id).First(&user).Error

	case constant.Reference:
		err = p.db.Where("reference = ?", id).First(&user).Error

	default:
		return nil, errors.New("identifier is not valid")
	}

	if err != nil {
		logger.Logger.Errorf("error getting client by %s: %s", identifier, err)
		return user, err
	}

	return user, nil
}

func (p *PostgresRepository) UserExists(id, identifier string) bool {
	var (
		err   error
		count int64
	)

	switch identifier {
	case constant.PhoneNumber:
		err = p.db.Model(&models.User{}).Where("phone_number = ?", id).Count(&count).Error

	case constant.Email:
		err = p.db.Model(&models.User{}).Where("email = ?", id).Count(&count).Error

	case constant.Bvn:
		err = p.db.Model(&models.User{}).Where("bvn = ?", id).Count(&count).Error

	case constant.Nin:
		err = p.db.Model(&models.User{}).Where("nin = ?", id).Count(&count).Error

	default:
		err = errors.New("identifier is not valid")
	}
	if err != nil {
		logger.Logger.Errorf("[UserExists]error getting user by %s: %s", identifier, err)
		return false
	}
	return count > 0
}

func (p *PostgresRepository) UpdateUser(id, identifier string, updates map[string]interface{}) error {
	var err error
	switch identifier {
	case constant.ID:
		err = p.db.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error

	case constant.PhoneNumber:
		err = p.db.Model(&models.User{}).Where("phone_number = ?", id).Updates(updates).Error

	case constant.Email:
		err = p.db.Model(&models.User{}).Where("email = ?", id).Updates(updates).Error

	case constant.Bvn:
		err = p.db.Model(&models.User{}).Where("bvn = ?", id).Updates(updates).Error

	case constant.Nin:
		err = p.db.Model(&models.User{}).Where("nin = ?", id).Updates(updates).Error

	default:
		return errors.New("identifier is not valid")
	}
	if err != nil {
		logger.Logger.Errorf("[UpdateUser]error updating user by %s: %s", identifier, err)
		return err
	}
	return nil
}

func (p *PostgresRepository) AdminDashboardCards() (any, error) {

	var stats models.DashboardStats

	// ---- SUPPLIERS ----
	if err := p.db.Model(&models.User{}).
		Where("role = ?", enum.Supplier).
		Count(&stats.TotalSuppliers).Error; err != nil {
		return nil, err
	}

	p.db.Model(&models.User{}).
		Where("role = ? AND status = ?", enum.Supplier, constant.Approved).
		Count(&stats.VerifiedSuppliers)

	p.db.Model(&models.User{}).
		Where("role = ? AND is_active = true", enum.Supplier).
		Count(&stats.ActiveSuppliers)

	// ---- BUYERS ----
	p.db.Model(&models.User{}).
		Where("role = ?", enum.Buyer).
		Count(&stats.TotalBuyers)

	p.db.Model(&models.User{}).
		Where("role = ? AND is_active = true", enum.Buyer).
		Count(&stats.ActiveBuyers)

	// ---- PRODUCTS ----
	p.db.Model(&models.Product{}).
		Count(&stats.TotalProducts)

	p.db.Model(&models.Product{}).
		Where("status = ?", constant.Approved).
		Count(&stats.ApprovedProducts)

	p.db.Model(&models.Product{}).
		Where("status = ?", constant.Pending).
		Count(&stats.PendingProducts)

	// ---- ORDERS ----
	p.db.Model(&models.Order{}).
		Count(&stats.TotalOrders)

	p.db.Model(&models.Order{}).
		Where("status = ?", constant.Completed).
		Count(&stats.CompletedOrders)

	p.db.Model(&models.Order{}).
		Where("status = ?", constant.InTransit).
		Count(&stats.InTransitOrders)

	return &stats, nil
}

func (p *PostgresRepository) GetSupplierCards() (any, error) {
	var rows []models.SupplierAgg
	stats := &models.SupplierStats{}

	err := p.db.Model(&models.User{}).
		Select("status, COUNT(*) as count").
		Where("role = ?", enum.Supplier).
		Group("status").
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	for _, r := range rows {
		stats.TotalSuppliers += r.Count

		switch r.Status {
		case constant.Invited:
			stats.Invited = r.Count
		case constant.Registering:
			stats.Registering = r.Count
		case constant.Pending:
			stats.Pending = r.Count
		case constant.Verified:
			stats.Verified = r.Count
		}
	}

	return stats, nil
}

func (p *PostgresRepository) GetSuppliers(pm *models.PaginationMetadata, status, searchText string) ([]models.User, *models.PaginationMetadata, error) {

	var (
		suppliers []models.User
		//	total int64
	)

	query := p.db.Model(&models.User{}).Where("role = ?", enum.Supplier).Order("created_at desc")

	if searchText != "" {
		search := "%" + searchText + "%"
		query = query.Where(`
			business_name ILIKE ?
			OR name ILIKE ?
			OR email ILIKE ?
			OR phone ILIKE ?
		`, search, search, search, search)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Scopes(Paginator(pm, models.User{}, query)).Find(&suppliers).Error
	if err != nil {
		return nil, pm, err
	}
	//if err = query.Count(&total).Error; err != nil {
	//	logger.Logger.Errorf("[GetSuppliers]count error: %s", err)
	//}
	//pm.TotalRecords = total
	//pm.TotalPages = int((total + int64(pm.PageSize) - 1) / int64(pm.PageSize))

	return suppliers, pm, nil
}

func (p *PostgresRepository) SupplierEditProduct(productID string, updateMap map[string]interface{}) error {

	err := p.db.Model(&models.Product{}).Where("id = ?", productID).Updates(updateMap).Error
	if err != nil {
		return err
	}

	return nil
}
