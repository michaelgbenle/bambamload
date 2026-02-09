package route

import (
	"bambamload/enum"
	supplierHandler "bambamload/handler/supplier"
	"bambamload/middleware"

	f "github.com/gofiber/fiber/v2"
)

func SupplierRoutes(app *f.App, h *supplierHandler.Handler) {

	app.Post("/auth/supplier/register", h.Register)

	supplier := app.Group("/api/supplier", middleware.Authenticate(enum.Supplier, h.PostgresRepository, h.RedisService))

	supplier.Post("/upload_kyc_docs", h.UploadKycDocuments)
	supplier.Post("/submit_business_profile", h.SubmitBusinessProfile)

	supplier.Post("/create_product", h.CreateProduct)
	supplier.Put("/product/:id", h.EditProduct)
	supplier.Get("/product/:id", h.GetProduct)
	supplier.Get("/products", h.GetProducts)
	supplier.Get("/products/stats", h.GetSupplierProductStats)

	supplier.Post("/product/images/:id", h.UploadProductImages)

	supplier.Post("/logout", h.LogoutSupplier)
}
