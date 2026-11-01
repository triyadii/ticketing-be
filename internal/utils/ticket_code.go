package utils

import (
	"fmt"
	"time"

	"ticketing-be/internal/models"

	"gorm.io/gorm"
)

// GenerateTicketNumber menghasilkan nomor tiket dengan format: IT-YYYYMM-XXXX
// Contoh: IT-202611-0001
func GenerateTicketNumber(db *gorm.DB) (string, error) {
	now := time.Now()
	prefix := fmt.Sprintf("IT-%s", now.Format("200601"))

	var count int64
	// Hitung tiket yang dibuat pada bulan ini yang diawali dengan prefix yang sama
	pattern := prefix + "-%"
	err := db.Model(&models.Ticket{}).Where("nomor_ticket LIKE ?", pattern).Count(&count).Error
	if err != nil {
		return "", err
	}

	nextNumber := count + 1
	ticketCode := fmt.Sprintf("%s-%04d", prefix, nextNumber)
	return ticketCode, nil
}
