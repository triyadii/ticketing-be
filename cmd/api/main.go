package main

import (
	"fmt"
	"log"
	"os"

	"ticketing-be/internal/config"
	"ticketing-be/internal/handler"
	"ticketing-be/internal/middleware"
	"ticketing-be/internal/repository"
	"ticketing-be/internal/service"

	_ "ticketing-be/docs" // Import swagger docs yang digenerate

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

// @title                      IT Helpdesk Ticketing & Koperasi API
// @version                    1.0
// @description                Dokumentasi RESTful API Backend Sistem IT Helpdesk Ticketing dan Manajemen Data Lama Koperasi berbasis Go Fiber & PostgreSQL.
// @contact.name               Tim IT Support
// @contact.email              it-support@deliserdangkab.go.id
// @host                       localhost:8086
// @BasePath                   /
// @schemes                    http https
// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
// @description                Ketik 'Bearer <token_jwt>' untuk mengakses endpoint yang dilindungi.
func main() {
	// 1. Muat Konfigurasi
	cfg := config.LoadConfig()

	// 2. Inisialisasi Database PostgreSQL
	db := config.InitDB(cfg)

	// 3. Pastikan folder uploads tersedia
	if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
		log.Printf("Peringatan: Gagal membuat direktori upload '%s': %v\n", cfg.UploadDir, err)
	}

	// 4. Inisialisasi Repository
	userRepo := repository.NewUserRepository(db)
	layananRepo := repository.NewLayananRepository(db)
	pelaporRepo := repository.NewPelaporRepository(db)
	ticketRepo := repository.NewTicketRepository(db)
	legacyRepo := repository.NewLegacyRepository(db)

	// 5. Inisialisasi Service
	authSvc := service.NewAuthService(userRepo, cfg)
	userSvc := service.NewUserService(userRepo)
	layananSvc := service.NewLayananService(layananRepo)
	pelaporSvc := service.NewPelaporService(pelaporRepo)
	ticketSvc := service.NewTicketService(db, ticketRepo, pelaporRepo, layananRepo, userRepo)
	legacySvc := service.NewLegacyService(legacyRepo)

	// 6. Inisialisasi Handler
	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)
	layananHandler := handler.NewLayananHandler(layananSvc)
	pelaporHandler := handler.NewPelaporHandler(pelaporSvc)
	ticketHandler := handler.NewTicketHandler(ticketSvc, cfg)
	legacyHandler := handler.NewLegacyHandler(legacySvc)

	// 7. Setup Fiber App
	app := fiber.New(fiber.Config{
		AppName:      cfg.AppName,
		ServerHeader: "GoFiber",
		BodyLimit:    50 * 1024 * 1024, // 50MB batas upload multipart
	})

	// Middleware Global
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))

	// Serve Static Files (Untuk akses foto bukti tiket)
	app.Static("/uploads", "./storage/uploads")

	// Endpoint Dokumentasi Swagger UI
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Welcome / Healthcheck
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"app":         cfg.AppName,
			"environment": cfg.AppEnv,
			"status":      "online",
			"docs":        fmt.Sprintf("%s/swagger/index.html", cfg.AppURL),
		})
	})

	// 8. Registrasi API Routes (/api/v1)
	api := app.Group("/api/v1")

	// --- Auth Route ---
	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Get("/profile", middleware.Protected(), authHandler.GetProfile)

	// --- Users (Petugas IT) Route ---
	users := api.Group("/users", middleware.Protected())
	users.Get("/", userHandler.GetUsers)
	users.Post("/", userHandler.CreateUser)
	users.Get("/:id", userHandler.GetUserByID)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)

	// --- Layanan Kendala IT Route ---
	layanan := api.Group("/layanan")
	layanan.Get("/", layananHandler.GetLayanans)
	layanan.Get("/:id", layananHandler.GetLayananByID)
	layanan.Post("/", middleware.Protected(), layananHandler.CreateLayanan)
	layanan.Put("/:id", middleware.Protected(), layananHandler.UpdateLayanan)
	layanan.Delete("/:id", middleware.Protected(), layananHandler.DeleteLayanan)

	// --- Pelapor Route ---
	pelapor := api.Group("/pelapor")
	pelapor.Post("/", pelaporHandler.CreatePelapor)
	pelapor.Get("/:nik", pelaporHandler.GetPelaporByNIK)
	pelapor.Get("/", middleware.Protected(), pelaporHandler.GetPelapors)
	pelapor.Put("/:nik", middleware.Protected(), pelaporHandler.UpdatePelapor)
	pelapor.Delete("/:nik", middleware.Protected(), pelaporHandler.DeletePelapor)

	// --- Tickets Route ---
	tickets := api.Group("/tickets")
	tickets.Post("/", ticketHandler.CreateTicket)
	tickets.Get("/", ticketHandler.GetTickets)
	tickets.Get("/:identifier", ticketHandler.GetTicketDetail)
	tickets.Patch("/:id/assign", middleware.Protected(), ticketHandler.AssignTicket)
	tickets.Patch("/:id/status", middleware.Protected(), ticketHandler.UpdateTicketStatus)
	tickets.Delete("/:id", middleware.Protected(), ticketHandler.DeleteTicket)

	// --- Koperasi (Data Lama & Export) Route ---
	koperasi := api.Group("/koperasi", middleware.Protected())
	koperasi.Get("/data-lama", legacyHandler.GetLegacyData)
	koperasi.Get("/data-lama/export", legacyHandler.ExportLegacyData)

	// 9. Jalankan Server
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server %s berjalan di port %s\n", cfg.AppName, cfg.AppPort)
	log.Printf("Dokumentasi Swagger tersedia di: %s/swagger/index.html\n", cfg.AppURL)

	if err := app.Listen(addr); err != nil {
		log.Fatalf("Server gagal dijalankan: %v", err)
	}
}
