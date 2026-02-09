package models

import "time"

type Product struct {
	Model
	SupplierID            string          `json:"supplier_id" gorm:"type:varchar(255)"`
	Name                  string          `json:"name" gorm:"type:varchar(255)"`
	Category              string          `json:"category" gorm:"type:varchar(100)"`
	Type                  string          `json:"type" gorm:"type:varchar(50)"`
	Description           string          `json:"description" gorm:"type:text"`
	BaseUnitPrice         int64           `json:"base_unit_price" gorm:"type:int"`
	Unit                  string          `json:"unit" gorm:"type:varchar(50)"`
	MinimumOrderQuantity  int64           `json:"minimum_order_quantity" gorm:"type:int"`
	PaymentTerms          string          `json:"payment_terms" gorm:"type:varchar(50)"` //prepayment or pay_on_delivery
	PaymentMethods        string          `json:"payment_methods" gorm:"type:text"`      // comma separated values
	CurrentStockQuantity  int64           `json:"current_stock_quantity" gorm:"type:int"`
	LowStockAlertLevel    int64           `json:"low_stock_alert_level" gorm:"type:int"`
	FulfilmentType        string          `json:"fulfilment_type" gorm:"type:varchar(50)"` //delivery,customer_pick_up,both
	EstimatedDeliveryTime string          `json:"estimated_delivery_time" gorm:"type:varchar(50)"`
	Status                string          `json:"status" gorm:"type:varchar(25);default:'pending'"`
	ApprovalStatus        string          `json:"approval_status" gorm:"type:varchar(25);default:'pending'"`
	DateApproved          time.Time       `json:"date_approved" gorm:"type:date"`
	DateRejected          time.Time       `json:"date_rejected" gorm:"type:date"`
	ApprovedBy            string          `json:"approved_by" gorm:"type:varchar(100)"`
	RejectedBy            string          `json:"rejected_by" gorm:"type:varchar(100)"`
	Rating                int             `json:"rating" gorm:"type:int"`
	RejectReason          string          `json:"reject_reason" gorm:"type:varchar(100)"`
	ProductUploads        []ProductUpload `json:"product_uploads" gorm:"foreignKey:ProductID"`
	Supplier              User            `json:"supplier" gorm:"foreignKey:SupplierID"`
}

type ProductUpload struct {
	Model
	ProductID  string `json:"product_id" gorm:"index;not null"`
	SupplierID string `json:"supplier_id" gorm:"type:varchar(100)"`
	FileURL    string `json:"file_url" gorm:"type:text"`
	FileType   string `json:"file_type" gorm:"type:varchar(50)"`
}

type ApproveOrRejectSupplierProduct struct {
	Action    string `json:"action"`
	Comment   string `json:"comment"`
	ProductID string `json:"product_id"`
}

type ProductStats struct {
	TotalProducts  int64 `json:"total_products"`
	ActiveListings int64 `json:"active_listings"`
	PendingReview  int64 `json:"pending_review"`
	Rejected       int64 `json:"rejected"`
	Deactivated    int64 `json:"deactivated"`
}

type SupplierProductStats struct {
	TotalProducts   int64 `json:"total_products"`
	ActiveProducts  int64 `json:"active_products"`
	LowStockItems   int64 `json:"low_stock_items"`
	OutOfStockItems int64 `json:"out_of_stock_items"`
}
