package buyer

import (
	"bambamload/constant"
	"bambamload/handler"
	"bambamload/logger"
	"bambamload/models"
	"bambamload/utils"
	"net/http"
	"strings"

	f "github.com/gofiber/fiber/v2"
)

type Handler struct {
	*handler.Handler
}

func NewBuyerHandler(apiHandler *handler.Handler) *Handler {
	return &Handler{
		apiHandler,
	}
}

func (h *Handler) Register(c *f.Ctx) error {

	var req models.BuyerRegisterRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "invalid request body", nil)
	}
	if req.Name == "" {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "name cannot be empty", nil)
	}
	if req.Email == "" {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "email cannot be empty", nil)
	}
	if req.PhoneNumber == "" {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "phone number cannot be empty", nil)
	}

	if req.Password == "" {
		return utils.WriteResponse(c, http.StatusBadRequest, false, "password cannot be empty", nil)
	}

	msg, err := h.BuyerService.Register(req) //nolint:typecheck
	if err != nil {
		return utils.WriteResponse(c, http.StatusInternalServerError, false, msg, nil)
	}
	sEmail := strings.ToLower(strings.TrimSpace(req.Email))

	//trigger otp after registration
	err = h.UtilitiesService.SendOtp(constant.RegisterOtp, sEmail)
	if err != nil {
		logger.Logger.Errorf("[Register]Send Otp error: %v", err)
	}

	//login user
	loginData, err := h.UtilitiesService.Login(models.LoginRequest{
		Email:    strings.ToLower(req.Email),
		Password: req.Password,
	})
	if err != nil {
		logger.Logger.Errorf("[Register]Login error: %v", err)
	}

	return utils.WriteResponse(c, http.StatusOK, true, "successful", loginData)
}

func (h *Handler) LogoutBuyer(c *f.Ctx) error {

	err := h.UtilitiesService.Logout(c.Locals(constant.Token).(string), c.Locals(constant.User).(*models.User)) //nolint:typecheck
	if err != nil {
		return utils.WriteResponse(c, http.StatusInternalServerError, false, err.Error(), nil)
	}

	return utils.WriteResponse(c, http.StatusOK, true, "successful", nil)
}
