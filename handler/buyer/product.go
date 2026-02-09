package buyer

import (
	"bambamload/constant"
	"bambamload/models"
	"bambamload/utils"
	"net/http"

	f "github.com/gofiber/fiber/v2"
)

func (h *Handler) GetProduct(c *f.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "id cannot be empty", nil)
	}

	product, err := h.SupplierService.SupplierGetProduct(id, nil)
	if err != nil {
		return utils.WriteResponse(c, http.StatusInternalServerError, false, "get product failed", nil)
	}
	return utils.WriteResponse(c, http.StatusOK, true, "success", product)
}

func (h *Handler) GetProducts(c *f.Ctx) error {
	_ = c.Locals("user").(*models.User)

	page := c.Query(constant.Page, "1")
	pageSize := c.Query(constant.PageSize, "10")
	pm := utils.InitPaginationMetadata(page, pageSize)
	searchText := c.Query("search_text", "")
	status := c.Query("status", "")
	productType := c.Query("type", "")

	products, paginationMeta, err := h.SupplierService.GetProducts(pm, status, searchText, productType)
	if err != nil {
		return utils.WriteResponse(c, http.StatusInternalServerError, false, err.Error(), nil)
	}

	resp := map[string]interface{}{
		"pagination_meta": paginationMeta,
		"products":        products,
	}
	return utils.WriteResponse(c, http.StatusOK, true, "success", resp)
}
