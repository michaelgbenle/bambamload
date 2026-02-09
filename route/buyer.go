package route

import (
	"bambamload/enum"
	buyerHandler "bambamload/handler/buyer"
	"bambamload/middleware"

	f "github.com/gofiber/fiber/v2"
)

func BuyerRoutes(app *f.App, h *buyerHandler.Handler) {

	app.Post("/auth/buyer/register", h.Register)

	buyer := app.Group("/api/buyer", middleware.Authenticate(enum.Buyer, h.PostgresRepository, h.RedisService))

	buyer.Get("/product/:id", h.GetProduct)
	buyer.Get("/products", h.GetProducts)

	buyer.Post("/logout", h.LogoutBuyer)
}
