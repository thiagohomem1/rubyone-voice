package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"rubyone-voice/config"
	"rubyone-voice/database"
	"rubyone-voice/controllers"
	"rubyone-voice/services"
	"rubyone-voice/routes"
)

func main() {
	// Carregar configuração
	config.LoadConfig()

	// Conectar ao banco de dados
	database.Connect()

	// Inicializar aplicação Fiber
	app := fiber.New(fiber.Config{
		AppName:      "RubyOne Voice SaaS",
		ServerHeader: "RubyOne Voice",
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return ctx.Status(code).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		},
	})

	// Inicializar Auth
	authService := services.NewAuthService(database.DB)
	authController := controllers.NewAuthController(authService)
	routes.SetupAuthRoutes(app, authController)

	// Inicializar Permissions (Stage 3)
	permissionService := services.NewPermissionService(database.DB)
	permissionController := controllers.NewPermissionController(permissionService)
	routes.SetupPermissionRoutes(app, permissionController, permissionService)


	// Middlewares
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - ${latency}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Endpoint de health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "RubyOne Voice SaaS is running 🚀",
		})
	})

	// Configurar graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Gracefully shutting down...")
		app.Shutdown()
	}()

	// Iniciar servidor
	port := config.AppConfig.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 RubyOne Voice SaaS iniciando na porta %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Falha ao iniciar servidor:", err)
	}

	log.Println("RubyOne Voice SaaS encerrado.")
} 