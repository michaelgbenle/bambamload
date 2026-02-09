package supplier

import (
	"bambamload/constant"
	"bambamload/logger"
	"bambamload/models"
)

func (ss *ServiceSupplier) CreateProduct(req *models.Product, user *models.User) error {
	req.SupplierID = user.ID
	return ss.PostgresRepository.CreateProduct(req)
}

func (ss *ServiceSupplier) SupplierGetProduct(id string, user *models.User) (*models.Product, error) {
	product, err := ss.PostgresRepository.GetProduct(id, constant.ID)
	if err != nil {
		logger.Logger.Errorf("Supplier GetProduct err: %v", err)
		return nil, err
	}
	return product, nil
}

func (ss *ServiceSupplier) GetProducts(pm *models.PaginationMetadata, status, searchText, productType string) ([]models.Product, *models.PaginationMetadata, error) {

	products, paginationMetaData, err := ss.PostgresRepository.GetProducts(pm, status, searchText, productType)
	if err != nil {
		logger.Logger.Errorf("Supplier Get Products Error: %s", err)
		return nil, paginationMetaData, err
	}
	return products, paginationMetaData, nil
}

func (ss *ServiceSupplier) EditProduct(id string, product *models.Product, user *models.User) error {

	updateMap := make(map[string]interface{})

	if product.Name != "" {
		updateMap["name"] = product.Name
	}

	if product.Category != "" {
		updateMap["category"] = product.Category
	}

	if product.Type != "" {
		updateMap["type"] = product.Type
	}

	if product.Description != "" {
		updateMap["description"] = product.Description
	}

	if product.BaseUnitPrice > 0 {
		updateMap["base_unit_price"] = product.BaseUnitPrice
	}
	if product.Unit != "" {
		updateMap["unit"] = product.Unit
	}
	if product.MinimumOrderQuantity > 0 {
		updateMap["minimum_order_quantity"] = product.MinimumOrderQuantity
	}
	if product.PaymentTerms != "" {
		updateMap["payment_terms"] = product.PaymentTerms
	}
	if product.PaymentMethods != "" {
		updateMap["payment_methods"] = product.PaymentMethods
	}

	if product.CurrentStockQuantity > 0 {
		updateMap["current_stock_quantity"] = product.CurrentStockQuantity
	}

	if product.LowStockAlertLevel > 0 {
		updateMap["low_stock_alert_level"] = product.LowStockAlertLevel
	}
	if product.FulfilmentType != "" {
		updateMap["fulfilment_type"] = product.FulfilmentType
	}

	if product.EstimatedDeliveryTime != "" {
		updateMap["estimated_delivery_time"] = product.EstimatedDeliveryTime
	}

	err := ss.PostgresRepository.SupplierEditProduct(id, updateMap)
	if err != nil {
		logger.Logger.Errorf("SupplierEditProduct Error: %v", err)
		return err
	}

	return nil

}

func (ss *ServiceSupplier) GetSupplierProductStats(user *models.User) (any, error) {

	data, err := ss.PostgresRepository.GetSupplierProductStats(user.ID)
	if err != nil {
		logger.Logger.Errorf("Supplier GetSupplierProductStats Error: %v", err)
		return nil, err
	}
	return data, nil
}
