package utils

import (
	"bambamload/models"
	"strconv"
)

func InitPaginationMetadata(pageString, pageSizeString string) *models.PaginationMetadata {
	pm := &models.PaginationMetadata{}

	if pageString != "" {
		pageInt, err := strconv.ParseInt(pageString, 10, 64)
		if err == nil {
			pm.Page = int(pageInt)
		}
	} else {
		pm.Page = 1
	}

	if pageSizeString != "" {
		pageSizeInt, err := strconv.ParseInt(pageSizeString, 10, 64)
		if err == nil {
			pm.PageSize = int(pageSizeInt)
		}
	} else {
		pm.PageSize = 10
	}

	return pm
}
