package models

type PaginationMetadata struct {
	TotalRecords int64  `json:"total_records"`
	TotalPages   int    `json:"total_pages"`
	PageSize     int    `json:"page_size"`
	Page         int    `json:"page"`
	SupplierID   string `json:"-"`
}
