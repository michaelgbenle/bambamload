package admin

import (
	"bambamload/constant"
	"bambamload/handler"
	"bambamload/models"
	"bambamload/utils"
	"net/http"

	f "github.com/gofiber/fiber/v2"
)

type Handler struct {
	*handler.Handler
}

func NewAdminHandler(apiHandler *handler.Handler) *Handler {
	return &Handler{
		apiHandler,
	}
}

func (h *Handler) ApproveOrRejectSupplierKyc(c *f.Ctx) error {
	user := c.Locals("user").(*models.User)
	var req models.ApproveOrRejectSupplierRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "invalid request body", nil)
	}

	if req.DocumentKey != constant.CacCertificate && req.DocumentKey != constant.ValidPersonalID && req.DocumentKey != constant.UtilityBill && req.DocumentKey != constant.TinDocument {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "document_key can only be 'cac_certificate','valid_personal_id','utility_bill','tin_document'", nil)
	}

	if req.Action != constant.Approve && req.Action != constant.Reject {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "action can only be approve/reject", nil)
	}

	err := h.AdminService.ApproveOrRejectSupplierKyc(req, user)
	if err != nil {
		return utils.WriteResponse(c, http.StatusInternalServerError, false, err.Error(), err)
	}
	return utils.WriteResponse(c, http.StatusOK, true, "success", nil)
}

func (h *Handler) ApproveOrRejectSupplier(c *f.Ctx) error {
	user := c.Locals("user").(*models.User)
	var req models.ApproveOrRejectSupplierRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "invalid request body", nil)
	}

	if req.SupplierID == "" && req.Action == "" {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "action / supplier id required", nil)
	}

	if req.Action != constant.Approve && req.Action != constant.Reject {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "action can only be approve/reject", nil)
	}

	err := h.AdminService.ApproveOrRejectSupplier(req, user)
	if err != nil {
		return utils.WriteResponse(c, http.StatusInternalServerError, false, err.Error(), err)
	}
	return utils.WriteResponse(c, http.StatusOK, true, "success", nil)
}

func (h *Handler) InviteSupplier(c *f.Ctx) error {
	user := c.Locals("user").(*models.User)
	var req models.InviteSupplier

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "invalid request body", nil)
	}
	if req.BusinessName == "" {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "business name is required", nil)
	}
	if req.ContactPerson == "" {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "contact person is required", nil)
	}
	if req.Email == "" {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "email is required", nil)
	}
	if req.PhoneNumber == "" {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "phone number is required", nil)
	}

	msg, err := h.AdminService.InviteSupplier(req, user)
	if err != nil {
		return utils.WriteResponse(c, http.StatusInternalServerError, false, msg, nil)
	}
	return utils.WriteResponse(c, http.StatusOK, true, msg, nil)
}

func (h *Handler) ResendInviteSupplierEmail(c *f.Ctx) error {
	reference := c.Query("reference")
	if reference == "" {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "reference is required", nil)
	}

	msg, err := h.AdminService.ResendSupplierInviteEmail(reference)
	if err != nil {
		return utils.WriteResponse(c, http.StatusInternalServerError, false, msg, nil)
	}
	return utils.WriteResponse(c, http.StatusOK, true, msg, nil)
}

func (h *Handler) LogoutAdmin(c *f.Ctx) error {

	err := h.UtilitiesService.Logout(c.Locals(constant.Token).(string), c.Locals(constant.User).(*models.User)) //nolint:typecheck
	if err != nil {
		return utils.WriteResponse(c, http.StatusInternalServerError, false, err.Error(), nil)
	}

	return utils.WriteResponse(c, http.StatusOK, true, "successful", nil)
}

func (h *Handler) SupplierDashboardCards(c *f.Ctx) error {

	cards, err := h.AdminService.SupplierDashboardCards()
	if err != nil {
		return utils.WriteResponse(c, http.StatusInternalServerError, false, err.Error(), nil)
	}
	return utils.WriteResponse(c, http.StatusOK, true, "success", cards)
}

func (h *Handler) GetSuppliers(c *f.Ctx) error {

	page := c.Query(constant.Page, "1")
	pageSize := c.Query(constant.PageSize, "10")
	pm := utils.InitPaginationMetadata(page, pageSize)
	searchText := c.Query("search_text", "")
	status := c.Query("status", "")

	suppliers, paginationMeta, err := h.AdminService.GetSuppliers(pm, status, searchText)
	if err != nil {
		return utils.WriteResponse(c, http.StatusInternalServerError, false, err.Error(), nil)
	}

	resp := map[string]interface{}{
		"pagination_meta": paginationMeta,
		"suppliers":       suppliers,
	}
	return utils.WriteResponse(c, http.StatusOK, true, "success", resp)
}

func (h *Handler) GetSupplier(c *f.Ctx) error {

	id := c.Params("id")
	if id == "" {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "id is required", nil)
	}

	supplier, err := h.AdminService.GetSupplier(id)
	if err != nil {
		return utils.WriteResponse(c, http.StatusInternalServerError, false, err.Error(), nil)
	}

	return utils.WriteResponse(c, http.StatusOK, true, "success", supplier)
}

func (h *Handler) DashboardCards(c *f.Ctx) error {

	cards, err := h.AdminService.DashboardCards()
	if err != nil {
		return utils.WriteResponse(c, http.StatusInternalServerError, false, err.Error(), nil)
	}
	return utils.WriteResponse(c, http.StatusOK, true, "success", cards)
}

func (h *Handler) ChangeSupplierCommissionRate(c *f.Ctx) error {
	user := c.Locals("user").(*models.User)
	var req models.EditSupplierCommissionRate

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "invalid request body", nil)
	}

	if req.SupplierID == "" {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "supplier id is required", nil)
	}
	if req.Rate <= 0 {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "rate is required", nil)
	}

	err := h.AdminService.ChangeSupplierCommissionRate(req, user)
	if err != nil {
		return utils.WriteResponse(c, http.StatusInternalServerError, false, err.Error(), nil)
	}

	return utils.WriteResponse(c, http.StatusOK, true, "success", nil)
}

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
func (h *Handler) GetAdminProductCards(c *f.Ctx) error {
	_ = c.Locals("user").(*models.User)

	data, err := h.AdminService.GetProductStats()
	if err != nil {
		return utils.WriteResponse(c, http.StatusInternalServerError, false, err.Error(), nil)
	}

	return utils.WriteResponse(c, http.StatusOK, true, "success", data)
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
func (h *Handler) ApproveOrRejectSupplierProduct(c *f.Ctx) error {

	user := c.Locals("user").(*models.User)
	var req models.ApproveOrRejectSupplierProduct

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "invalid request body", nil)
	}

	if req.Action != constant.Approve && req.Action != constant.Reject {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "action can only be approve/reject", nil)
	}
	if req.ProductID == "" {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "product id is required", nil)
	}
	err := h.AdminService.ApproveOrRejectSupplierProduct(req, user)
	if err != nil {
		return utils.WriteResponse(c, http.StatusInternalServerError, false, err.Error(), nil)
	}
	return utils.WriteResponse(c, http.StatusOK, true, "success", nil)
}
