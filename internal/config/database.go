package config

import (
	"fmt"
	"log"
	"time"

	"ticketing-be/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg *Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode, cfg.DBTimeZone,
	)

	dbLogger := logger.Default.LogMode(logger.Info)
	if cfg.AppEnv == "production" {
		dbLogger = logger.Default.LogMode(logger.Error)
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get generic database object: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("PostgreSQL connection established successfully")

	// Run Auto Migration
	AutoMigrate(DB)

	// Seed Initial Data
	SeedInitialData(DB)

	return DB
}

func AutoMigrate(db *gorm.DB) {
	log.Println("Running database auto-migrations...")
	err := db.AutoMigrate(
		&models.User{},
		&models.Layanan{},
		&models.Pelapor{},
		&models.Ticket{},
		&models.LegacyKoperasi{},
	)
	if err != nil {
		log.Fatalf("Auto-migration failed: %v", err)
	}
	log.Println("Auto-migrations completed successfully")
}

func SeedInitialData(db *gorm.DB) {
	// 1. Seed Default Admin User
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	if userCount == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		if err == nil {
			admin := models.User{
				Nama:     "Super Admin IT",
				Username: "admin",
				Password: string(hashedPassword),
				Status:   models.UserStatusAktif,
			}
			if err := db.Create(&admin).Error; err == nil {
				log.Println("Default Admin IT created (username: admin, password: password123)")
			}
		}
	}

	// 2. Seed Default Layanan IT
	var layananCount int64
	db.Model(&models.Layanan{}).Count(&layananCount)
	if layananCount == 0 {
		layanans := []models.Layanan{
			{NamaLayanan: "Jaringan & Wi-Fi", Status: models.LayananStatusAktif},
			{NamaLayanan: "Hardware & PC / Laptop", Status: models.LayananStatusAktif},
			{NamaLayanan: "Aplikasi Internal & Website", Status: models.LayananStatusAktif},
			{NamaLayanan: "Printer & Scanner", Status: models.LayananStatusAktif},
			{NamaLayanan: "Email & Akun Sistem", Status: models.LayananStatusAktif},
		}
		for _, l := range layanans {
			db.Create(&l)
		}
		log.Println("Initial Layanan IT categories created")
	}

	// 3. Seed Default Sample Data Lama Koperasi jika kosong
	var legacyCount int64
	db.Model(&models.LegacyKoperasi{}).Count(&legacyCount)
	if legacyCount == 0 {
		now := time.Now()
		legacyData := []models.LegacyKoperasi{
			{KodeAnggota: "KOP-001", NamaAnggota: "Budi Santoso", UnitKerja: "Dinas Pendidikan", JenisTransaksi: "Simpanan Pokok", Nominal: 500000, Keterangan: "Setoran awal anggota baru", Tanggal: now.AddDate(0, -5, -10), CreatedAt: now},
			{KodeAnggota: "KOP-002", NamaAnggota: "Siti Rahma", UnitKerja: "Dinas Kesehatan", JenisTransaksi: "Simpanan Wajib", Nominal: 150000, Keterangan: "Iuran rutin bulanan", Tanggal: now.AddDate(0, -4, -5), CreatedAt: now},
			{KodeAnggota: "KOP-003", NamaAnggota: "Ahmad Dahlan", UnitKerja: "Inspektorat", JenisTransaksi: "Pinjaman Reguler", Nominal: 5000000, Keterangan: "Pinjaman modal usaha", Tanggal: now.AddDate(0, -3, -15), CreatedAt: now},
			{KodeAnggota: "KOP-004", NamaAnggota: "Dewi Lestari", UnitKerja: "Badan Kepegawaian", JenisTransaksi: "Simpanan Sukarela", Nominal: 300000, Keterangan: "Tabungan sukarela", Tanggal: now.AddDate(0, -2, -2), CreatedAt: now},
			{KodeAnggota: "KOP-005", NamaAnggota: "Rudi Hartono", UnitKerja: "Dinas Kominfo", JenisTransaksi: "Angsuran Pinjaman", Nominal: 850000, Keterangan: "Angsuran ke-3 pinjaman", Tanggal: now.AddDate(0, -1, -20), CreatedAt: now},
		}
		for _, item := range legacyData {
			db.Create(&item)
		}
		log.Println("Sample data lama koperasi created")
	}
}
