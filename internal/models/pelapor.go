package models

import (
	"time"
)

type Pelapor struct {
	NIK       string    `gorm:"type:varchar(16);primaryKey" json:"nik"`
	Nama      string    `gorm:"type:varchar(100);not null" json:"nama"`
	Instansi  string    `gorm:"type:varchar(150);not null" json:"instansi"`
	NoHP      string    `gorm:"type:varchar(20);not null" json:"no_hp"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Tickets []Ticket `gorm:"foreignKey:NikPelapor;references:NIK" json:"tickets,omitempty"`
}
