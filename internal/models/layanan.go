package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LayananStatus string

const (
	LayananStatusAktif    LayananStatus = "aktif"
	LayananStatusNonaktif LayananStatus = "nonaktif"
)

type Layanan struct {
	ID          uuid.UUID     `gorm:"type:uuid;primaryKey" json:"id"`
	NamaLayanan string        `gorm:"type:varchar(100);not null" json:"nama_layanan"`
	Status      LayananStatus `gorm:"type:varchar(20);default:'aktif';not null" json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

func (l *Layanan) BeforeCreate(tx *gorm.DB) (err error) {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	if l.Status == "" {
		l.Status = LayananStatusAktif
	}
	return
}
