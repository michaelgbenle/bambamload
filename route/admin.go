package route

import (
	"bambamload/enum"
	adminhandler "bambamload/handler/admin"
	"bambamload/middleware"

	f "github.com/gofiber/fiber/v2"
)

func AdminRoutes(app *f.App, h *adminhandler.Handler) {

	app.Get("/v1/logs/:date", h.GetLogs)

	//app.Post("/auth/admin/register", h.RegisterAdmin)

	admin := app.Group("/api/admin", middleware.Authenticate(enum.Admin, h.PostgresRepository, h.RedisService))

	admin.Get("/me", h.Me)
	admin.Post("/invite/supplier", h.InviteSupplier)
	admin.Post("/invite/resend", h.ResendInviteSupplierEmail)

	admin.Get("/dashboard/cards", h.DashboardCards)

	//suppliers
	admin.Post("/supplier/approve_or_reject", h.ApproveOrRejectSupplier)
	admin.Post("/supplier/kyc/approve_or_reject", h.ApproveOrRejectSupplierKyc)

	admin.Get("/supplier/:id", h.GetSupplier)
	admin.Get("/suppliers", h.GetSuppliers)
	admin.Get("/suppliers/cards", h.SupplierDashboardCards)
	admin.Post("/suppliers/commision_rate", h.ChangeSupplierCommissionRate)

	//products
	admin.Get("/product/:id", h.GetProduct)
	admin.Get("/products", h.GetProducts)
	admin.Get("/products/cards", h.GetAdminProductCards)

	admin.Post("/products/approve_or_reject", h.ApproveOrRejectSupplierProduct)

	admin.Post("/logout", h.LogoutAdmin)

}
