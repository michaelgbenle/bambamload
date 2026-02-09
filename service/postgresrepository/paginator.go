package postgresrepository

import (
	"bambamload/models"

	"gorm.io/gorm"
)

// Paginator ...
func Paginator(p *models.PaginationMetadata, model interface{}, db *gorm.DB) func(*gorm.DB) *gorm.DB {
	var total int64
	if model != nil {
		db.Count(&total)
	} else {
		db.Model(model).Count(&total)
	}

	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 10
	}
	p.TotalRecords = total
	p.TotalPages = int((total + int64(p.PageSize) - 1) / int64(p.PageSize))

	return func(db *gorm.DB) *gorm.DB {
		offset := (p.Page - 1) * p.PageSize
		return db.Offset(offset).Limit(p.PageSize)
	}
}
