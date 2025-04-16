package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/latoulicious/siresto-backend/internal/domain"
	"github.com/latoulicious/siresto-backend/internal/handler"
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
	// QR Code domain
	qrRepo := &repository.QRCodeRepository{DB: db}
	qrService := &service.QRCodeService{Repo: qrRepo}
	qrHandler := &handler.QRCodeHandler{Service: qrService}

	// Category domain
	categoryRepo := &repository.CategoryRepository{DB: db}
	categoryService := &service.CategoryService{Repo: categoryRepo}
	categoryHandler := &handler.CategoryHandler{Service: categoryService}

	// Product domain
	productRepo := &repository.ProductRepository{DB: db}
	productService := &service.ProductService{Repo: productRepo}
	productHandler := &handler.ProductHandler{Service: productService}

	// Variation domain
	variationRepo := &repository.VariationRepository{DB: db}
	variationService := &service.VariationService{Repo: variationRepo}
	variationHandler := &handler.VariationHandler{Service: variationService}

	// User domain
	userRepo := &repository.UserRepository{DB: db}
	userService := &service.UserService{Repo: userRepo}
	userHandler := &handler.UserHandler{Service: userService}

	// Role domain
	roleRepo := &repository.RoleRepository{DB: db}
	roleService := &service.RoleService{Repo: roleRepo}
	roleHandler := &handler.RoleHandler{Service: roleService}

	// API v1
	v1 := app.Group("/api/v1")

	// General routes
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

	// Auth routes
	v1.Post("/auth/login", userHandler.LoginUser)
	logger.LogInfo("POST /api/v1/auth/login route registered", logutil.Route("POST", "/api/v1/auth/login"))

	// User routes
	v1.Get("/users", userHandler.ListAllUsers)
	logger.LogInfo("GET /api/v1/users route registered", logutil.Route("GET", "/api/v1/users"))

	v1.Get("/users/:id", userHandler.GetUserByID)
	logger.LogInfo("GET /api/v1/users/:id route registered", logutil.Route("GET", "/api/v1/users/:id"))

	v1.Post("/users", userHandler.CreateUser)
	logger.LogInfo("POST /api/v1/users route registered", logutil.Route("POST", "/api/v1/users"))

	v1.Put("/users/:id", userHandler.UpdateUser)
	logger.LogInfo("PUT /api/v1/users/:id route registered", logutil.Route("PUT", "/api/v1/users/:id"))

	v1.Delete("/users/:id", userHandler.DeleteUser)
	logger.LogInfo("DELETE /api/v1/users/:id route registered", logutil.Route("DELETE", "/api/v1/users/:id"))

	// Role routes
	v1.Get("/roles", roleHandler.ListAllRoles)
	logger.LogInfo("GET /api/v1/roles route registered", logutil.Route("GET", "/api/v1/roles"))

	v1.Get("/roles/:id", roleHandler.GetRoleByID)
	logger.LogInfo("GET /api/v1/roles/:id route registered", logutil.Route("GET", "/api/v1/roles/:id"))

	v1.Post("/roles", roleHandler.CreateRole)
	logger.LogInfo("POST /api/v1/roles route registered", logutil.Route("POST", "/api/v1/roles"))

	v1.Put("/roles/:id", roleHandler.UpdateRole)
	logger.LogInfo("PUT /api/v1/roles/:id route registered", logutil.Route("PUT", "/api/v1/roles/:id"))

	v1.Delete("/roles/:id", roleHandler.DeleteRole)
	logger.LogInfo("DELETE /api/v1/roles/:id route registered", logutil.Route("DELETE", "/api/v1/roles/:id"))

	// QR Code routes
	v1.Get("/qr-codes", qrHandler.ListAllQRCodesHandler)
	logger.LogInfo("GET /api/v1/qr-codes route registered", logutil.Route("GET", "/api/v1/qr-codes"))

	v1.Get("/qr-codes/:id", qrHandler.GetQRCodeByIDHandler)
	logger.LogInfo("GET /api/v1/qr-codes/:id route registered", logutil.Route("GET", "/api/v1/qr-codes/:id"))

	v1.Get("/qr-codes/store/:store_id", qrHandler.ListQRCodesHandler)
	logger.LogInfo("GET /api/v1/qr-codes/store/:store_id route registered", logutil.Route("GET", "/api/v1/qr-codes/store/:store_id"))

	v1.Post("/qr-codes", qrHandler.CreateQRCodeHandler)
	logger.LogInfo("POST /api/v1/qr-codes route registered", logutil.Route("POST", "/api/v1/qr-codes"))

	v1.Post("/qr-codes/bulk", qrHandler.BulkCreateQRCodeHandler)
	logger.LogInfo("POST /api/v1/qr-codes/bulk route registered", logutil.Route("POST", "/api/v1/qr-codes/bulk"))

	v1.Put("/qr-codes/:id", qrHandler.UpdateQRCodeHandler)
	logger.LogInfo("PUT /api/v1/qr-codes/:id route registered", logutil.Route("PUT", "/api/v1/qr-codes/:id"))

	v1.Delete("/qr-codes/:id", qrHandler.DeleteQRCodeHandler)
	logger.LogInfo("DELETE /api/v1/qr-codes/:id route registered", logutil.Route("DELETE", "/api/v1/qr-codes/:id"))

	// Category routes
	v1.Get("/categories", categoryHandler.ListAllCategories)
	logger.LogInfo("GET /api/v1/categories route registered", logutil.Route("GET", "/api/v1/categories"))

	v1.Get("/categories/:id", categoryHandler.GetCategoryByID)
	logger.LogInfo("GET /api/v1/categories/:id route registered", logutil.Route("GET", "/api/v1/categories/:id"))

	v1.Post("/categories", categoryHandler.CreateCategory)
	logger.LogInfo("POST /api/v1/categories route registered", logutil.Route("POST", "/api/v1/categories"))

	v1.Put("/categories/:id", categoryHandler.UpdateCategory)
	logger.LogInfo("PUT /api/v1/categories/:id route registered", logutil.Route("PUT", "/api/v1/categories/:id"))

	v1.Delete("/categories/:id", categoryHandler.DeleteCategory)
	logger.LogInfo("DELETE /api/v1/categories/:id route registered", logutil.Route("DELETE", "/api/v1/categories/:id"))

	// Product routes
	v1.Get("/products", productHandler.ListAllProducts)
	logger.LogInfo("GET /api/v1/products route registered", logutil.Route("GET", "/api/v1/products"))

	v1.Get("/products/:id", productHandler.GetProductByID)
	logger.LogInfo("GET /api/v1/products/:id route registered", logutil.Route("GET", "/api/v1/products/:id"))

	v1.Post("/products", productHandler.CreateProduct)
	logger.LogInfo("POST /api/v1/products route registered", logutil.Route("POST", "/api/v1/products"))

	v1.Put("/products/:id", productHandler.UpdateProduct)
	logger.LogInfo("PUT /api/v1/products/:id route registered", logutil.Route("PUT", "/api/v1/products/:id"))

	v1.Delete("/products/:id", productHandler.DeleteProduct)
	logger.LogInfo("DELETE /api/v1/products/:id route registered", logutil.Route("DELETE", "/api/v1/products/:id"))

	// Variation Routes (Not tied to a specific product)
	v1.Get("/variations", variationHandler.ListAllVariations)
	logger.LogInfo("GET /api/v1/variations route registered", logutil.Route("GET", "/api/v1/variations"))

	v1.Get("/variations/:id", variationHandler.GetVariationByID)
	logger.LogInfo("GET /api/v1/variations/:id route registered", logutil.Route("GET", "/api/v1/variations/:id"))

	v1.Post("/variations", variationHandler.CreateVariation)
	logger.LogInfo("POST /api/v1/variations route registered", logutil.Route("POST", "/api/v1/variations"))

	v1.Put("/variations/:id", variationHandler.UpdateVariation)
	logger.LogInfo("PUT /api/v1/variations/:id route registered", logutil.Route("PUT", "/api/v1/variations/:id"))

	v1.Delete("/variations/:id", variationHandler.DeleteVariation)
	logger.LogInfo("DELETE /api/v1/variations/:id route registered", logutil.Route("DELETE", "/api/v1/variations/:id"))

	// TODO
	// Variation Routes (Tied to a specific product)

	// v1.Get("/products/:product_id/variations", variationHandler.ListVariations)
	// logger.LogInfo("GET /api/v1/products/:product_id/variations route registered", logutil.Route("GET", "/api/v1/products/:product_id/variations"))

	// v1.Post("/products/:product_id/variations", variationHandler.CreateVariation)
	// logger.LogInfo("POST /api/v1/products/:product_id/variations route registered", logutil.Route("POST", "/api/v1/products/:product_id/variations"))

	// v1.Put("/products/:product_id/variations/:id", variationHandler.UpdateVariation)
	// logger.LogInfo("PUT /api/v1/products/:product_id/variations/:id route registered", logutil.Route("PUT", "/api/v1/products/:product_id/variations/:id"))

	// v1.Delete("/products/:product_id/variations/:id", variationHandler.DeleteVariation)
	// logger.LogInfo("DELETE /api/v1/products/:product_id/variations/:id route registered", logutil.Route("DELETE", "/api/v1/products/:product_id/variations/:id"))

	// Log routes
	v1.Get("/logs", func(c *fiber.Ctx) error {
		var logs []domain.Log
		if err := db.Find(&logs).Error; err != nil {
			logger.LogError("Failed to fetch logs", logutil.Route("GET", "/api/v1/logs"))
			return c.Status(fiber.StatusInternalServerError).JSON(utils.Error("Failed to retrieve Logs", fiber.StatusInternalServerError))
		}

		return c.Status(fiber.StatusOK).JSON(utils.Success("Logs fetched successfully", logs))
	})
	logger.LogInfo("GET /api/v1/logs route registered", logutil.Route("GET", "/api/v1/logs"))
}
