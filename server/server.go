package server

import (
	"bambamload/constant"
	"bambamload/handler"
	adminhandler "bambamload/handler/admin"
	buyerhandler "bambamload/handler/buyer"
	supplierhandler "bambamload/handler/supplier"
	utilitieshandler "bambamload/handler/utilities"
	"bambamload/logger"
	"bambamload/middleware"
	"bambamload/route"
	"bambamload/service/admin"
	"bambamload/service/buyer"
	"bambamload/service/email"
	"bambamload/service/postgresrepository"
	"bambamload/service/redisService"
	"bambamload/service/supplier"
	uploadservice "bambamload/service/uploadService"
	"bambamload/service/utilities"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	f "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func Start() {
	logger.InitLogger(os.Getenv(constant.AppEnv))
	middleware.StartBackgroundRequestLogger()

	rs := redisService.NewRedisService()
	pg := postgresrepository.NewPostgresRepository()
	defer pg.Close()
	emailService := email.NewEmailService()
	//smsService := smsservice.NewSMSService()
	uploadService := uploadservice.NewUploadService()
	adminService := admin.NewServiceAdmin(rs, pg, *emailService)
	supplierService := supplier.NewServiceSupplier(rs, pg, *emailService, uploadService)
	buyerService := buyer.NewServiceBuyer(rs, pg, *emailService, uploadService)
	utilitiesService := utilities.NewServiceUtilities(rs, pg, *emailService, uploadService)
	apiHandler := handler.NewHandler(rs, *pg, *emailService, uploadService, adminService, supplierService, buyerService, utilitiesService)
	adminHandler := adminhandler.NewAdminHandler(apiHandler)
	supplierHandler := supplierhandler.NewSupplierHandler(apiHandler)
	buyerHandler := buyerhandler.NewBuyerHandler(apiHandler)
	utilitiesHandler := utilitieshandler.NewUtilitiesHandler(apiHandler)

	app := f.New()

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     fmt.Sprintf("%s", os.Getenv("ALLOWED_ORIGINS")),
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     fmt.Sprintf("%s,%s,%s,%s,%s", constant.Origin, constant.Authorization, constant.ContentType, constant.ContentLength, constant.XRequestedWith),
		ExposeHeaders:    fmt.Sprintf("%s,%s,%s", constant.ContentType, constant.ContentLength, constant.Authorization),
		AllowCredentials: false,
	}))

	// Rate Limiter
	app.Use(limiter.New(limiter.Config{
		Expiration: 5 * time.Second,
		Max:        10,
	}))

	app.Use(middleware.APILogger())

	// Routes
	app.Get("/", apiHandler.WelcomeHandler)
	app.Get("/health", apiHandler.Health)
	route.AdminRoutes(app, adminHandler)
	route.SupplierRoutes(app, supplierHandler)
	route.BuyerRoutes(app, buyerHandler)
	route.UtilitiesRoutes(app, utilitiesHandler)

	// Not found
	app.Use(apiHandler.NotFoundHandler)

	// Graceful shutdown
	go func() {
		if err := app.Listen(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))); err != nil {
			logger.Logger.Fatalf("failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Logger.Infof("shutting down server...")
	if err := app.Shutdown(); err != nil {
		logger.Logger.Fatalf("failed to stop server: %v", err)
	}
}
