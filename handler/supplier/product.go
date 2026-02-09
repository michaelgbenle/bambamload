package supplier

import (
	"bambamload/constant"
	"bambamload/logger"
	"bambamload/models"
	"bambamload/utils"
	"net/http"
	"path/filepath"

	f "github.com/gofiber/fiber/v2"
)

func (h *Handler) CreateProduct(c *f.Ctx) error {
	user := c.Locals("user").(*models.User)
	var req models.Product

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "invalid request body", nil)
	}
	if req.Name == "" {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "name cannot be empty", nil)
	}
	if req.Category == "" {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "category cannot be empty", nil)
	}
	if req.Type == "" {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "phone type cannot be empty", nil)
	}

	err := h.SupplierService.CreateProduct(&req, user)
	if err != nil {
		return utils.WriteResponse(c, http.StatusInternalServerError, false, "create product failed", nil)
	}

	return utils.WriteResponse(c, http.StatusOK, true, "create product successfully", nil)
}
func (h *Handler) GetProduct(c *f.Ctx) error {
	user := c.Locals("user").(*models.User)

	id := c.Params("id")
	if id == "" {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "id cannot be empty", nil)
	}

	product, err := h.SupplierService.SupplierGetProduct(id, user)
	if err != nil {
		return utils.WriteResponse(c, http.StatusInternalServerError, false, "get product failed", nil)
	}
	return utils.WriteResponse(c, http.StatusOK, true, "success", product)
}

func (h *Handler) GetSupplierProductStats(c *f.Ctx) error {
	user := c.Locals("user").(*models.User)

	data, err := h.SupplierService.GetSupplierProductStats(user)
	if err != nil {
		return utils.WriteResponse(c, http.StatusInternalServerError, false, "get product stats failed", nil)
	}
	return utils.WriteResponse(c, http.StatusOK, true, "success", data)
}

func (h *Handler) GetProducts(c *f.Ctx) error {
	user := c.Locals("user").(*models.User)

	page := c.Query(constant.Page, "1")
	pageSize := c.Query(constant.PageSize, "10")
	pm := utils.InitPaginationMetadata(page, pageSize)
	searchText := c.Query("search_text", "")
	status := c.Query("status", "")
	productType := c.Query("type", "")

	pm.SupplierID = user.ID

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

func (h *Handler) EditProduct(c *f.Ctx) error {
	user := c.Locals("user").(*models.User)
	var req models.Product

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "invalid request body", nil)
	}
	id := c.Params("id")
	if id == "" {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "id cannot be empty", nil)
	}

	err := h.SupplierService.EditProduct(id, &req, user)
	if err != nil {
		return utils.WriteResponse(c, http.StatusInternalServerError, false, "edit product failed", nil)
	}

	return utils.WriteResponse(c, http.StatusOK, true, "success", nil)
}

func (h *Handler) UploadProductImages(c *f.Ctx) error {
	user := c.Locals("user").(*models.User)

	productID := c.Params("id")
	if productID == "" {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "product id cannot be empty", nil)
	}

	imageFields := []string{"image_one", "image_two", "image_three", "image_four", "image_five"}
	errs := make(map[string]string)

	var uploads []models.ProductUpload

	for _, field := range imageFields {
		file, err := c.FormFile(field)
		if err != nil {
			logger.Logger.Errorf("[UploadProduct] FormFile error: %v", err)
			errs[field] = err.Error()
			continue
		}

		// Validate Format (PDF, JPG, PNG)
		ext := filepath.Ext(file.Filename)
		if !utils.IsValidExtension(ext) {
			errs[field] = "file extension is invalid"
			continue
		}

		ff, e := file.Open()
		if e != nil {
			logger.Logger.Errorf("[UploadProduct]Open file error: %v", e)
			continue
		}

		url, err := h.UploadService.Upload(ff, file.Filename)
		if err != nil {
			logger.Logger.Errorf("[UploadProduct]Upload error: %v", err)
			errs[field] = err.Error()
			return err
		}

		uploads = append(uploads, models.ProductUpload{
			SupplierID: user.ID,
			ProductID:  productID,
			FileURL:    url,
			FileType:   ext,
		})

	}
	if len(errs) > 0 {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "upload product images error", errs)
	}

	err := h.PostgresRepository.BatchInsertProductUploads(uploads)
	if err != nil {
		logger.Logger.Errorf("[UploadProduct]BatchInsertProductUploads error: %v", err)
		return utils.WriteResponse(c, http.StatusInternalServerError, false, "batch insert product uploads error", nil)
	}

	return utils.WriteResponse(c, http.StatusOK, true, "success", nil)
}
