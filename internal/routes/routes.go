package routes

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/latoulicious/siresto-backend/internal/config"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/handler"
	"github.com/latoulicious/siresto-backend/internal/middleware"
	"github.com/latoulicious/siresto-backend/internal/repository"
	"github.com/latoulicious/siresto-backend/internal/service"
	"github.com/latoulicious/siresto-backend/internal/utils"
	"github.com/latoulicious/siresto-backend/pkg/core/logging"
	"github.com/latoulicious/siresto-backend/pkg/logutil"
	"gorm.io/gorm"
)

// Health Check Max Response Time
const responseTimeThreshold = 500 * time.Millisecond

func SetupRoutes(app *fiber.App, db *gorm.DB, logger logging.Logger) {

	//* Menu domain

	// QR Code domain
	qrRepo := &repository.QRCodeRepository{DB: db}
	qrService := &service.QRCodeService{Repo: qrRepo}
	qrHandler := &handler.QRCodeHandler{Service: qrService}

	// Category domain
	categoryRepo := &repository.CategoryRepository{DB: db}
	categoryService := &service.CategoryService{Repo: categoryRepo}
	categoryHandler := &handler.CategoryHandler{Service: categoryService}

	// Initialize R2 Uploader
	var r2Uploader *utils.R2Uploader
	r2Uploader, err := config.NewR2UploaderFromEnv()
	if err != nil {
		logger.LogError("Failed to initialize R2 uploader", logutil.MainCall("init", "r2uploader", map[string]interface{}{
			"error": err.Error(),
		}))
		logger.LogInfo("Products can be created but image upload functionality will be disabled",
			logutil.MainCall("init", "r2uploader", map[string]interface{}{
				"suggestion": "Check your R2 configuration in .env file",
			}))
		// r2Uploader will remain nil, which is fine
	} else {
		logger.LogInfo("R2 uploader initialized successfully",
			logutil.MainCall("init", "r2uploader", map[string]interface{}{
				"status": "ready",
			}))
	}

	// Product domain
	productRepo := &repository.ProductRepository{DB: db}
	productService := &service.ProductService{
		Repo:     productRepo,
		Uploader: r2Uploader,
	}
	productHandler := &handler.ProductHandler{Service: productService}

	// Variation domain
	variationRepo := &repository.VariationRepository{DB: db}
	variationService := &service.VariationService{Repo: variationRepo}
	variationHandler := &handler.VariationHandler{
		Service:        variationService,
		ProductService: productService,
	}

	//* Core Domain

	// User domain
	validate := validator.New()
	userRepo := &repository.UserRepository{DB: db}
	userService := &service.UserService{Repo: userRepo}
	userHandler := handler.NewUserHandler(userService, validate)

	// Permission domain
	permissionRepo := &repository.PermissionRepository{DB: db}
	permissionService := &service.PermissionService{Repo: permissionRepo}
	permissionHandler := &handler.PermissionHandler{Service: permissionService}

	// Role domain
	roleRepo := &repository.RoleRepository{DB: db}
	roleService := &service.RoleService{
		Repo:              roleRepo,
		PermissionService: permissionService,
		DB:                db,
	}
	roleHandler := &handler.RoleHandler{Service: roleService}

	// Order domain
	orderRepo := &repository.OrderRepository{DB: db}
	orderService := &service.OrderService{
		Repo:          orderRepo,
		ProductRepo:   productRepo,
		VariationRepo: variationRepo,
	}
	orderHandler := &handler.OrderHandler{OrderService: orderService}

	// Payment domain
	paymentRepo := &repository.PaymentRepository{DB: db}
	paymentService := &service.PaymentService{Repo: paymentRepo}
	paymentHandler := &handler.PaymentHandler{Service: paymentService}

	//* Utility Domain

	// Theme domain
	themeRepo := &repository.ThemeRepository{DB: db}
	themeService := &service.ThemeService{Repo: themeRepo}
	themeHandler := &handler.ThemeHandler{Service: themeService}

	// API v1
	v1 := app.Group("/api/v1")

	// Public routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"code":    fiber.StatusOK,
			"status":  "success",
			"message": "SiResto API is running",
		})
	})
	logger.LogInfo("GET / root route registered", logutil.Route("GET", "/"))

	app.Get("/health", func(c *fiber.Ctx) error {
		startTime := time.Now()

		// Simulate a delay by introducing random sleep time
		// randomDelay := time.Duration(rand.Intn(10000)) * time.Millisecond
		// time.Sleep(randomDelay) // Random delay between 0 to 1 second

		// Check the response time
		elapsedTime := time.Since(startTime)
		if elapsedTime > responseTimeThreshold {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":        "unhealthy",
				"timestamp":     time.Now().Format(time.RFC3339),
				"service":       "siresto-api",
				"version":       "1.0.0",
				"message":       "Service is unhealthy due to high response time",
				"response_time": elapsedTime.String(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":        "healthy",
			"timestamp":     time.Now().Format(time.RFC3339),
			"service":       "siresto-api",
			"version":       "1.0.0",
			"response_time": elapsedTime.String(),
		})
	})
	logger.LogInfo("GET /health route registered", logutil.Route("GET", "/health"))

	// Auth routes (public)
	v1.Post("/auth/login", userHandler.LoginUser)
	logger.LogInfo("POST /api/v1/auth/login route registered", logutil.Route("POST", "/api/v1/auth/login"))

	// Protected routes
	protected := v1.Use(middleware.Protected())

	// User routes - Admin only
	adminRoutes := protected.Use(middleware.RequireAdmin())
	adminRoutes.Get("/users", userHandler.ListAllUsers)
	logger.LogInfo("GET /api/v1/users route registered (admin)", logutil.Route("GET", "/api/v1/users"))

	adminRoutes.Post("/users", userHandler.CreateUser)
	logger.LogInfo("POST /api/v1/users route registered (admin)", logutil.Route("POST", "/api/v1/users"))

	adminRoutes.Put("/users/:id", userHandler.UpdateUser)
	logger.LogInfo("PUT /api/v1/users/:id route registered (admin)", logutil.Route("PUT", "/api/v1/users/:id"))

	adminRoutes.Delete("/users/:id", userHandler.DeleteUser)
	logger.LogInfo("DELETE /api/v1/users/:id route registered (admin)", logutil.Route("DELETE", "/api/v1/users/:id"))

	// User profile - Any authenticated user can access their own profile
	protected.Get("/users/:id", userHandler.GetUserByID)
	logger.LogInfo("GET /api/v1/users/:id route registered", logutil.Route("GET", "/api/v1/users/:id"))

	// Role & Permission routes - Admin only
	adminRoutes.Get("/roles", roleHandler.ListAllRoles)
	logger.LogInfo("GET /api/v1/roles route registered (admin)", logutil.Route("GET", "/api/v1/roles"))

	adminRoutes.Get("/roles/:id", roleHandler.GetRoleByID)
	logger.LogInfo("GET /api/v1/roles/:id route registered (admin)", logutil.Route("GET", "/api/v1/roles/:id"))

	adminRoutes.Post("/roles", roleHandler.CreateRole)
	logger.LogInfo("POST /api/v1/roles route registered (admin)", logutil.Route("POST", "/api/v1/roles"))

	adminRoutes.Put("/roles/:id", roleHandler.UpdateRole)
	logger.LogInfo("PUT /api/v1/roles/:id route registered (admin)", logutil.Route("PUT", "/api/v1/roles/:id"))

	adminRoutes.Delete("/roles/:id", roleHandler.DeleteRole)
	logger.LogInfo("DELETE /api/v1/roles/:id route registered (admin)", logutil.Route("DELETE", "/api/v1/roles/:id"))

	// Roles Permission routes - Admin only
	adminRoutes.Post("/roles/:id/permissions", roleHandler.AddPermissionsToRole)
	logger.LogInfo("POST /api/v1/roles/:id/permission route registered (admin)", logutil.Route("POST", "/api/v1/roles/:id/permission"))

	adminRoutes.Delete("/roles/:id/permissions", roleHandler.RemovePermissionsFromRole)
	logger.LogInfo("DELETE /api/v1/roles/:id/permission route registered (admin)", logutil.Route("DELETE", "/api/v1/roles/:id/permission"))

	// Permission routes - Admin only
	adminRoutes.Get("/permissions", permissionHandler.ListAllPermissions)
	logger.LogInfo("GET /api/v1/permissions route registered (admin)", logutil.Route("GET", "/api/v1/permissions"))

	adminRoutes.Get("/permissions/:id", permissionHandler.GetPermissionByID)
	logger.LogInfo("GET /api/v1/permissions/:id route registered (admin)", logutil.Route("GET", "/api/v1/permissions/:id"))

	adminRoutes.Post("/permissions", permissionHandler.CreatePermission)
	logger.LogInfo("POST /api/v1/permissions route registered (admin)", logutil.Route("POST", "/api/v1/permissions"))

	adminRoutes.Put("/permissions/:id", permissionHandler.UpdatePermission)
	logger.LogInfo("PUT /api/v1/permissions/:id route registered (admin)", logutil.Route("PUT", "/api/v1/permissions/:id"))

	adminRoutes.Delete("/permissions/:id", permissionHandler.DeletePermission)
	logger.LogInfo("DELETE /api/v1/permissions/:id route registered (admin)", logutil.Route("DELETE", "/api/v1/permissions/:id"))

	// QR Code routes
	protected.Get("/qr-codes", qrHandler.ListAllQRCodesHandler)
	logger.LogInfo("GET /api/v1/qr-codes route registered", logutil.Route("GET", "/api/v1/qr-codes"))

	protected.Get("/qr-codes/:id", qrHandler.GetQRCodeByIDHandler)
	logger.LogInfo("GET /api/v1/qr-codes/:id route registered", logutil.Route("GET", "/api/v1/qr-codes/:id"))

	protected.Get("/qr-codes/store/:store_id", qrHandler.ListQRCodesHandler)
	logger.LogInfo("GET /api/v1/qr-codes/store/:store_id route registered", logutil.Route("GET", "/api/v1/qr-codes/store/:store_id"))

	protected.Post("/qr-codes", qrHandler.CreateQRCodeHandler)
	logger.LogInfo("POST /api/v1/qr-codes route registered", logutil.Route("POST", "/api/v1/qr-codes"))

	protected.Post("/qr-codes/bulk", qrHandler.BulkCreateQRCodeHandler)
	logger.LogInfo("POST /api/v1/qr-codes/bulk route registered", logutil.Route("POST", "/api/v1/qr-codes/bulk"))

	protected.Delete("/qr-codes/:id", qrHandler.DeleteQRCodeHandler)
	logger.LogInfo("DELETE /api/v1/qr-codes/:id route registered", logutil.Route("DELETE", "/api/v1/qr-codes/:id"))

	// Category routes - Public access for listing
	v1.Get("/categories", categoryHandler.ListAllCategories)
	logger.LogInfo("GET /api/v1/categories route registered (public)", logutil.Route("GET", "/api/v1/categories"))

	// Protected Category routes
	protected.Get("/categories/:id", categoryHandler.GetCategoryByID)
	logger.LogInfo("GET /api/v1/categories/:id route registered", logutil.Route("GET", "/api/v1/categories/:id"))

	protected.Post("/categories", categoryHandler.CreateCategory)
	logger.LogInfo("POST /api/v1/categories route registered", logutil.Route("POST", "/api/v1/categories"))

	protected.Put("/categories/:id", categoryHandler.UpdateCategory)
	logger.LogInfo("PUT /api/v1/categories/:id route registered", logutil.Route("PUT", "/api/v1/categories/:id"))

	protected.Delete("/categories/:id", categoryHandler.DeleteCategory)
	logger.LogInfo("DELETE /api/v1/categories/:id route registered", logutil.Route("DELETE", "/api/v1/categories/:id"))

	// Product routes
	protected.Get("/products", productHandler.ListAllProducts)
	logger.LogInfo("GET /api/v1/products route registered", logutil.Route("GET", "/api/v1/products"))

	protected.Get("/products/:id", productHandler.GetProductByID)
	logger.LogInfo("GET /api/v1/products/:id route registered", logutil.Route("GET", "/api/v1/products/:id"))

	protected.Post("/products", productHandler.CreateProduct)
	logger.LogInfo("POST /api/v1/products/with-image route registered", logutil.Route("POST", "/api/v1/products/with-image"))

	protected.Put("/products/:id", productHandler.UpdateProduct)
	logger.LogInfo("PUT /api/v1/products/:id route registered", logutil.Route("PUT", "/api/v1/products/:id"))

	protected.Delete("/products/:id", productHandler.DeleteProduct)
	logger.LogInfo("DELETE /api/v1/products/:id route registered", logutil.Route("DELETE", "/api/v1/products/:id"))

	// Variation Routes (Not tied to a specific product)
	protected.Get("/variations", variationHandler.ListAllVariations)
	logger.LogInfo("GET /api/v1/variations route registered", logutil.Route("GET", "/api/v1/variations"))

	protected.Get("/variations/:id", variationHandler.GetVariationByID)
	logger.LogInfo("GET /api/v1/variations/:id route registered", logutil.Route("GET", "/api/v1/variations/:id"))

	protected.Post("/variations", variationHandler.CreateVariation)
	logger.LogInfo("POST /api/v1/variations route registered", logutil.Route("POST", "/api/v1/variations"))

	protected.Put("/variations/:id", variationHandler.UpdateVariation)
	logger.LogInfo("PUT /api/v1/variations/:id route registered", logutil.Route("PUT", "/api/v1/variations/:id"))

	protected.Delete("/variations/:id", variationHandler.DeleteVariation)
	logger.LogInfo("DELETE /api/v1/variations/:id route registered", logutil.Route("DELETE", "/api/v1/variations/:id"))

	// Variation Routes (Tied to a specific product)
	protected.Get("/products/:product_id/variations", variationHandler.GetProductVariations)
	logger.LogInfo("GET /api/v1/products/:product_id/variations route registered", logutil.Route("GET", "/api/v1/products/:product_id/variations"))

	protected.Post("/products/:product_id/variations", variationHandler.CreateProductVariation)
	logger.LogInfo("POST /api/v1/products/:product_id/variations route registered", logutil.Route("POST", "/api/v1/products/:product_id/variations"))

	protected.Put("/products/:product_id/variations/:id", variationHandler.UpdateProductVariation)
	logger.LogInfo("PUT /api/v1/products/:product_id/variations/:id route registered", logutil.Route("PUT", "/api/v1/products/:product_id/variations/:id"))

	protected.Delete("/products/:product_id/variations/:id", variationHandler.DeleteProductVariation)
	logger.LogInfo("DELETE /api/v1/products/:product_id/variations/:id route registered", logutil.Route("DELETE", "/api/v1/products/:product_id/variations/:id"))

	// Order routes
	v1.Post("/orders", orderHandler.CreateOrder)
	logger.LogInfo("POST /api/v1/orders route registered (public)", logutil.Route("POST", "/api/v1/orders"))

	protected.Get("/orders", orderHandler.ListAllOrders)
	logger.LogInfo("GET /api/v1/orders route registered", logutil.Route("GET", "/api/v1/orders"))

	protected.Get("/orders/:id", orderHandler.GetOrderByID)
	logger.LogInfo("GET /api/v1/orders/:id route registered", logutil.Route("GET", "/api/v1/orders/:id"))

	protected.Put("/orders/:id", orderHandler.UpdateOrder)
	logger.LogInfo("PUT /api/v1/orders/:id route registered", logutil.Route("PUT", "/api/v1/orders/:id"))

	// Order Status
	protected.Post("/orders/:orderID/complete", orderHandler.MarkOrderAsCompleted)
	logger.LogInfo("POST /api/v1/orders/:orderID/complete route registered", logutil.Route("POST", "/api/v1/orders/:orderID/complete"))

	protected.Post("/orders/:orderID/cancel", orderHandler.MarkOrderAsCanceled)
	logger.LogInfo("POST /api/v1/orders/:orderID/cancel route registered", logutil.Route("POST", "/api/v1/orders/:orderID/cancel"))

	// Order Payment
	protected.Get("/payments", paymentHandler.ListAllOrderPayments)
	logger.LogInfo("GET /api/v1/payments route registered", logutil.Route("GET", "/api/v1/payments"))

	protected.Get("/orders/:orderID/payments", paymentHandler.GetOrderPayments)
	logger.LogInfo("GET /api/v1/orders/:orderID/payments route registered", logutil.Route("GET", "/api/v1/orders/:orderID/payments"))

	protected.Post("/orders/:orderID/payments", paymentHandler.ProcessOrderPayment)
	logger.LogInfo("POST /api/v1/orders/:orderID/payments route registered", logutil.Route("POST", "/api/v1/orders/:orderID/payments"))

	// Utility routes
	v1.Get("/themes", themeHandler.ListAllThemes)
	logger.LogInfo("GET /api/v1/themes route registered (public)", logutil.Route("GET", "/api/v1/themes"))

	protected.Get("/themes/:id", themeHandler.GetThemeByID)
	logger.LogInfo("GET /api/v1/themes/:id route registered", logutil.Route("GET", "/api/v1/themes/:id"))

	protected.Post("/themes", themeHandler.CreateTheme)
	logger.LogInfo("POST /api/v1/themes route registered", logutil.Route("POST", "/api/v1/themes"))

	protected.Put("/themes/:id", themeHandler.UpdateTheme)
	logger.LogInfo("PUT /api/v1/themes/:id route registered", logutil.Route("PUT", "/api/v1/themes/:id"))

	protected.Delete("/themes/:id", themeHandler.DeleteTheme)
	logger.LogInfo("DELETE /api/v1/themes/:id route registered", logutil.Route("DELETE", "/api/v1/themes/:id"))

	// Log routes
	protected.Get("/logs", func(c *fiber.Ctx) error {
		var logs []domain.Log
		if err := db.Find(&logs).Error; err != nil {
			logger.LogError("Failed to fetch logs", logutil.Route("GET", "/api/v1/logs"))
			return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to retrieve Logs", fiber.StatusInternalServerError))
		}

		return c.Status(fiber.StatusOK).JSON(utils.Success("Logs fetched successfully", logs))
	})
	logger.LogInfo("GET /api/v1/logs route registered", logutil.Route("GET", "/api/v1/logs"))
}
