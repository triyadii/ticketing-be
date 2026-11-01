package repository

import (
	"ticketing-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketFilter struct {
	Page        int
	PerPage     int
	Status      string
	IdPelayanan string
	IdPetugas   string
	NikPelapor  string
	StartDate   string
	EndDate     string
	Search      string
}

type TicketRepository interface {
	Create(ticket *models.Ticket) error
	FindByID(id uuid.UUID) (*models.Ticket, error)
	FindByNomorTicket(nomor string) (*models.Ticket, error)
	FindAll(filters TicketFilter) ([]models.Ticket, int64, error)
	Update(ticket *models.Ticket) error
	Delete(id uuid.UUID) error
}

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) Create(ticket *models.Ticket) error {
	return r.db.Create(ticket).Error
}

func (r *ticketRepository) FindByID(id uuid.UUID) (*models.Ticket, error) {
	var ticket models.Ticket
	err := r.db.Preload("Pelapor").Preload("Layanan").Preload("Petugas").First(&ticket, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) FindByNomorTicket(nomor string) (*models.Ticket, error) {
	var ticket models.Ticket
	err := r.db.Preload("Pelapor").Preload("Layanan").Preload("Petugas").First(&ticket, "nomor_ticket = ?", nomor).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) FindAll(filters TicketFilter) ([]models.Ticket, int64, error) {
	var tickets []models.Ticket
	var total int64

	query := r.db.Model(&models.Ticket{}).Preload("Pelapor").Preload("Layanan").Preload("Petugas")

	if filters.Status != "" {
		query = query.Where("tickets.status = ?", filters.Status)
	}

	if filters.IdPelayanan != "" {
		query = query.Where("tickets.id_pelayanan = ?", filters.IdPelayanan)
	}

	if filters.IdPetugas != "" {
		query = query.Where("tickets.id_petugas = ?", filters.IdPetugas)
	}

	if filters.NikPelapor != "" {
		query = query.Where("tickets.nik_pelapor = ?", filters.NikPelapor)
	}

	if filters.StartDate != "" {
		query = query.Where("DATE(tickets.created_at) >= ?", filters.StartDate)
	}

	if filters.EndDate != "" {
		query = query.Where("DATE(tickets.created_at) <= ?", filters.EndDate)
	}

	if filters.Search != "" {
		searchPattern := "%" + filters.Search + "%"
		query = query.Where("tickets.nomor_ticket ILIKE ? OR tickets.judul_kendala ILIKE ? OR tickets.keterangan ILIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filters.Page - 1) * filters.PerPage
	err := query.Order("tickets.created_at DESC").Offset(offset).Limit(filters.PerPage).Find(&tickets).Error
	return tickets, total, err
}

func (r *ticketRepository) Update(ticket *models.Ticket) error {
	return r.db.Save(ticket).Error
}

func (r *ticketRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Ticket{}, "id = ?", id).Error
}
