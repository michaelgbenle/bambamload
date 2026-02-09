package route

import (
	utilitiesHandler "bambamload/handler/utilities"

	f "github.com/gofiber/fiber/v2"
)

func UtilitiesRoutes(app *f.App, h *utilitiesHandler.Handler) {

	app.Post("/utilities/auth/user/login", h.Login)
	app.Post("/utilities/auth/password/forgot", h.ForgotPassword)
	app.Post("/utilities/auth/password/verify", h.VerifyAndChangeForgotPassword)

	app.Post("/utilities/auth/verify_registration_otp", h.VerifyRegistrationOtp)

	app.Post("/utilities/auth/otp/resend", h.ResendOtp)

	app.Get("/utilities/user_by_reference", h.GetUserRegistrationDetailsByReference)

}
