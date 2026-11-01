package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketStatus string

const (
	TicketStatusMenunggu   TicketStatus = "menunggu"
	TicketStatusDiproses   TicketStatus = "diproses"
	TicketStatusSelesai    TicketStatus = "selesai"
	TicketStatusDibatalkan TicketStatus = "dibatalkan"
)

type Ticket struct {
	ID           uuid.UUID    `gorm:"type:uuid;primaryKey" json:"id"`
	NomorTicket  string       `gorm:"type:varchar(50);uniqueIndex;not null" json:"nomor_ticket"`
	NikPelapor   string       `gorm:"type:varchar(16);not null;index" json:"nik_pelapor"`
	IdPelayanan  uuid.UUID    `gorm:"type:uuid;not null;index" json:"id_pelayanan"`
	IdPetugas    *uuid.UUID   `gorm:"type:uuid;index" json:"id_petugas"`
	JudulKendala string       `gorm:"type:varchar(150);not null" json:"judul_kendala"`
	Keterangan   string       `gorm:"type:text;not null" json:"keterangan"`
	Foto         StringArray  `gorm:"type:json;not null" json:"foto"`
	Solusi       *string      `gorm:"type:text" json:"solusi"`
	Status       TicketStatus `gorm:"type:varchar(20);default:'menunggu';not null;index" json:"status"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	FinishedAt   *time.Time   `json:"finished_at"`

	// Relasi
	Pelapor *Pelapor `gorm:"foreignKey:NikPelapor;references:NIK" json:"pelapor,omitempty"`
	Layanan *Layanan `gorm:"foreignKey:IdPelayanan;references:ID" json:"layanan,omitempty"`
	Petugas *User    `gorm:"foreignKey:IdPetugas;references:ID" json:"petugas,omitempty"`
}

func (t *Ticket) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	if t.Status == "" {
		t.Status = TicketStatusMenunggu
	}
	if t.Foto == nil {
		t.Foto = []string{}
	}
	return
}
