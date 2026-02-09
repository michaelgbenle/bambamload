package utilities

import (
	"bambamload/constant"
	"bambamload/enum"
	"bambamload/handler"
	"bambamload/models"
	"bambamload/utils"
	"net/http"

	f "github.com/gofiber/fiber/v2"
)

type Handler struct {
	*handler.Handler
}

func NewUtilitiesHandler(apiHandler *handler.Handler) *Handler {
	return &Handler{
		apiHandler,
	}
}

//func (h *Handler) VerifyLoginOtp(c *f.Ctx) error {
//	var req models.VerifyLoginAdminOtpRequest
//
//	// Parse request body
//	if err := c.BodyParser(&req); err != nil {
//		return utils.WriteResponse(c, http.StatusBadRequest, false, "invalid request body", nil)
//	}
//
//	if req.Email == "" || req.Code == "" {
//		return utils.WriteResponse(c, http.StatusBadRequest, false, "email and code are required", nil)
//	}
//
//	err := h.AdminService.VerifyAdminOtp(constant.LoginOtp, req.Code, req.Email, nil) //nolint:typecheck
//	if err != nil {
//		return utils.WriteResponse(c, http.StatusInternalServerError, false, err.Error(), nil)
//	}
//
//	return utils.WriteResponse(c, http.StatusOK, true, "login successful", nil)
//}

func (h *Handler) Login(c *f.Ctx) error {
	var req models.LoginRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "invalid request body", nil)
	}

	if req.Email == "" || req.Password == "" {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "email and password are required", nil)
	}

	loginData, err := h.UtilitiesService.Login(req) //nolint:typecheck
	if err != nil {
		return utils.WriteResponse(c, http.StatusInternalServerError, false, err.Error(), nil)
	}

	return utils.WriteResponse(c, http.StatusOK, true, "login successful", loginData)
}

func (h *Handler) ResendOtp(c *f.Ctx) error {
	//actions : register_otp, forgot_password,
	var req models.ResendOtpRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "invalid request body", nil)
	}

	if req.Email == "" {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "email is required", nil)
	}

	err := h.UtilitiesService.SendOtp(req.Action, req.Email) //nolint:typecheck
	if err != nil {
		return utils.WriteResponse(c, http.StatusInternalServerError, false, err.Error(), nil)
	}

	return utils.WriteResponse(c, http.StatusOK, true, "successful", nil)
}

func (h *Handler) ForgotPassword(c *f.Ctx) error {
	var req models.ForgotPasswordRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "invalid request body", nil)
	}

	if req.Email == "" {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "email is required", nil)
	}

	if err := h.UtilitiesService.SendOtp(constant.ForgotPassword, req.Email); err != nil { //nolint:typecheck
		return utils.WriteResponse(c, http.StatusInternalServerError, false, "unable to send otp, please try again", nil)
	}

	return utils.WriteResponse(c, http.StatusOK, true, "successful", nil)
}

func (h *Handler) VerifyAndChangeForgotPassword(c *f.Ctx) error {
	var req models.VerifyForgotPasswordRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "invalid request body", nil)
	}

	if req.Email == "" || req.Code == "" || req.Password == "" {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "phone number,code and new password is required", nil)
	}

	if err := h.UtilitiesService.VerifyOtp(constant.ForgotPassword, req.Code, req.Email, req.Password); err != nil { //nolint:typecheck
		return utils.WriteResponse(c, http.StatusInternalServerError, false, "unable to verify otp, please try again", nil)
	}

	return utils.WriteResponse(c, http.StatusOK, true, "successful", nil)
}

func (h *Handler) GetUserRegistrationDetailsByReference(c *f.Ctx) error {
	reference := c.Query("reference")
	if reference == "" {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "reference is required", nil)
	}

	user, err := h.PostgresRepository.GetUser(reference, constant.Reference)
	if err != nil {
		return utils.WriteResponse(c, http.StatusInternalServerError, false, err.Error(), nil)
	}
	if user.Role != enum.Supplier {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "user is not a supplier", nil)
	}

	resp := map[string]interface{}{
		"business_name": user.BusinessName,
		"email":         user.Email,
		"name":          user.Name,
		"phone_number":  user.PhoneNumber,
		"reference":     reference,
	}

	return utils.WriteResponse(c, http.StatusOK, true, "successful", resp)
}

func (h *Handler) VerifyRegistrationOtp(c *f.Ctx) error {
	var req models.VerifyRegistrationOtpRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "invalid request body", nil)
	}

	if req.Email == "" || req.Code == "" {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "phone number,code and new password is required", nil)
	}

	if err := h.UtilitiesService.VerifyOtp(constant.RegisterOtp, req.Code, req.Email, nil); err != nil { //nolint:typecheck
		return utils.WriteResponse(c, http.StatusInternalServerError, false, "unable to verify otp, please try again", nil)
	}

	return utils.WriteResponse(c, http.StatusOK, true, "successful", nil)
}
