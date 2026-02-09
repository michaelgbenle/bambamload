package postgresrepository

import (
	"bambamload/constant"
	"bambamload/logger"
	"bambamload/models"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CreateProduct
func (p *PostgresRepository) CreateProduct(req *models.Product) error {
	return p.db.Create(req).Error
}

// GetProduct fetches a product by any identifier provided
func (p *PostgresRepository) GetProduct(id, identifier string) (*models.Product, error) {
	var (
		product *models.Product
		err     error
	)

	switch identifier {
	case constant.ID:
		err = p.db.Preload(clause.Associations).Where("id = ?", id).First(&product).Error

	default:
		return nil, errors.New("identifier is not valid")
	}

	if err != nil {
		logger.Logger.Errorf("error getting product by %s: %s", identifier, err)
		return product, err
	}

	return product, nil
}

func (p *PostgresRepository) UpdateProduct(id, identifier string, updates map[string]interface{}) error {
	var err error
	switch identifier {
	case constant.ID:
		err = p.db.Model(&models.Product{}).Where("id = ?", id).Updates(updates).Error

	default:
		return errors.New("identifier is not valid")
	}
	if err != nil {
		logger.Logger.Errorf("[UpdateUser]error updating product by %s: %s", identifier, err)
		return err
	}
	return nil
}

func (p *PostgresRepository) GetProducts(pm *models.PaginationMetadata, status, searchText, productType string) ([]models.Product, *models.PaginationMetadata, error) {

	var (
		products []models.Product
		query    *gorm.DB
	)

	if pm.SupplierID == "" {
		query = p.db.Model(&models.Product{}).Order("created_at desc")
	} else {
		query = p.db.Model(&models.Product{}).Where("supplier_id = ?", pm.SupplierID).Order("created_at desc")
	}

	if searchText != "" {
		search := "%" + searchText + "%"
		query = query.Where(`
			name ILIKE ?
			OR type ILIKE ?
			OR category ILIKE ?
		`, search, search, search)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if productType != "" {
		query = query.Where("type = ?", productType)
	}

	err := query.Scopes(Paginator(pm, &models.Product{}, query)).Find(&products).Error
	if err != nil {
		return nil, pm, err
	}

	return products, pm, nil
}

func (p *PostgresRepository) BatchInsertProductUploads(uploads []models.ProductUpload) error {
	if len(uploads) == 0 {
		return nil
	}

	return p.db.Create(&uploads).Error
}

func (p *PostgresRepository) GetProductStats() (*models.ProductStats, error) {
	var stats models.ProductStats

	err := p.db.Raw(`
		SELECT
			COUNT(*) AS total_products,
			SUM(CASE WHEN status = 'active' AND approval_status = 'approved' THEN 1 ELSE 0 END) AS active_listings,
			SUM(CASE WHEN approval_status = 'pending' THEN 1 ELSE 0 END) AS pending_review,
			SUM(CASE WHEN approval_status = 'rejected' THEN 1 ELSE 0 END) AS rejected,
			SUM(CASE WHEN status = 'deactivated' THEN 1 ELSE 0 END) AS deactivated
		FROM products
	`).Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (p *PostgresRepository) GetSupplierProductStats(supplierID string) (*models.SupplierProductStats, error) {
	var stats models.SupplierProductStats

	err := p.db.Raw(`
		SELECT
			COUNT(*)::INTEGER AS total_products,
			COALESCE(SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END)::INTEGER, 0) AS active_products,
			COALESCE(SUM(
				CASE 
					WHEN current_stock_quantity > 0 
					AND current_stock_quantity <= low_stock_alert_level 
					THEN 1 
					ELSE 0 
				END
			)::INTEGER, 0) AS low_stock_items,
			COALESCE(SUM(
				CASE 
					WHEN current_stock_quantity = 0 
					THEN 1 
					ELSE 0 
				END
			)::INTEGER, 0) AS out_of_stock_items
		FROM products
		WHERE supplier_id = $1
	`, supplierID).Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	return &stats, nil
}
