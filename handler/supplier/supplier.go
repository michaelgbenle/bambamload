package supplier

import (
	"bambamload/constant"
	"bambamload/handler"
	"bambamload/logger"
	"bambamload/models"
	"bambamload/utils"
	"net/http"
	"path/filepath"
	"strings"

	f "github.com/gofiber/fiber/v2"
)

type Handler struct {
	*handler.Handler
}

func NewSupplierHandler(apiHandler *handler.Handler) *Handler {
	return &Handler{
		apiHandler,
	}
}

func (h *Handler) Register(c *f.Ctx) error {

	var req models.SupplierRegisterRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "invalid request body", nil)
	}
	if req.Password == "" {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "password cannot be empty", nil)
	}
	if req.Reference == "" {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "reference cannot be empty", nil)
	}
	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		logger.Logger.Errorf("[Supplier] Register: hash password error: %v", err)
		return utils.WriteResponse(c, http.StatusBadRequest, false, "hash password error", nil)
	}
	unHashedPassword := req.Password
	req.Password = hashPassword

	msg, email, err := h.SupplierService.Register(&req) //nolint:typecheck
	if err != nil {
		return utils.WriteResponse(c, http.StatusInternalServerError, false, msg, nil)
	}

	//trigger otp after registration
	err = h.UtilitiesService.SendOtp(constant.RegisterOtp, email)
	if err != nil {
		logger.Logger.Errorf("[Register]Send Otp error: %v", err)
	}

	//login user
	loginData, err := h.UtilitiesService.Login(models.LoginRequest{
		Email:    strings.ToLower(strings.TrimSpace(email)),
		Password: unHashedPassword,
	})
	if err != nil {
		logger.Logger.Errorf("[Register]Login error: %v", err)
	}

	return utils.WriteResponse(c, http.StatusOK, true, "successful", loginData)
}

func (h *Handler) LogoutSupplier(c *f.Ctx) error {

	err := h.UtilitiesService.Logout(c.Locals(constant.Token).(string), c.Locals(constant.User).(*models.User)) //nolint:typecheck
	if err != nil {
		return utils.WriteResponse(c, http.StatusInternalServerError, false, err.Error(), nil)
	}

	return utils.WriteResponse(c, http.StatusOK, true, "successful", nil)
}

func (h *Handler) SubmitBusinessProfile(c *f.Ctx) error {
	user := c.Locals(constant.User).(*models.User)
	var req models.SubmitBusinessProfileRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "invalid request body", nil)
	}

	err := h.SupplierService.SubmitBusinessProfile(&req, user)
	if err != nil {
		return utils.WriteResponse(c, http.StatusInternalServerError, false, err.Error(), nil)
	}

	return utils.WriteResponse(c, http.StatusOK, true, "successful", nil)
}
func (h *Handler) UploadKycDocuments(c *f.Ctx) error {

	user := c.Locals(constant.User).(*models.User)
	//if user.KycStatus != constant.IdentityVerification {
	//	return utils.WriteResponse(c, http.StatusForbidden, false, "you can't upload kyc documents yet", nil)
	//}

	documentFields := []string{"cac_certificate", "valid_personal_id", "utility_bill", "tin_document"}

	// To store successful upload info
	results := make(map[string]interface{})
	errs := make(map[string]string)

	for _, field := range documentFields {

		file, err := c.FormFile(field)
		if err != nil {
			logger.Logger.Errorf("[UploadKycDocuments] FormFile error: %v", err)
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
			logger.Logger.Errorf("[UploadKycDocuments]Open file error: %v", e)
			continue
		}

		//  Send it to Backblaze B2
		// uploadedUrl, err := uploadToB2(file)

		url, err := h.UploadService.Upload(ff, file.Filename)
		if err != nil {
			logger.Logger.Errorf("[UploadKycDocuments]Upload error: %v", err)
			errs[field] = err.Error()
			return err
		}
		results[field] = url
	}
	if len(errs) > 0 {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "upload kyc documents error", errs)
	}

	results["kyc_status"] = constant.InReview
	err := h.PostgresRepository.UpdateUser(user.ID, constant.ID, results)
	if err != nil {
		logger.Logger.Errorf("[UploadKycDocuments]UpdateUser error: %v", err)
		return utils.WriteResponse(c, http.StatusInternalServerError, false, err.Error(), nil)
	}

	return utils.WriteResponse(c, http.StatusOK, true, "successful", nil)
}
