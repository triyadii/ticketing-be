package service

import (
	"errors"
	"time"

	"ticketing-be/internal/models"
	"ticketing-be/internal/repository"
	"ticketing-be/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateTicketInput struct {
	// Data Pelapor
	NIKPelapor string `form:"nik_pelapor" json:"nik_pelapor"`
	Nama       string `form:"nama" json:"nama"`
	Instansi   string `form:"instansi" json:"instansi"`
	NoHP       string `form:"no_hp" json:"no_hp"`

	// Data Tiket
	IdPelayanan  uuid.UUID `form:"id_pelayanan" json:"id_pelayanan"`
	JudulKendala string    `form:"judul_kendala" json:"judul_kendala"`
	Keterangan   string    `form:"keterangan" json:"keterangan"`
	FotoPaths    []string  `json:"foto_paths"`
}

type UpdateTicketStatusInput struct {
	Status models.TicketStatus `json:"status" example:"selesai"`
	Solusi *string             `json:"solusi,omitempty" example:"Kabel LAN diganti dan konfigurasi IP direset"`
}

type TicketService interface {
	CreateTicket(input CreateTicketInput) (*models.Ticket, error)
	GetTickets(filter repository.TicketFilter) ([]models.Ticket, int64, error)
	GetTicketByIDOrNumber(identifier string) (*models.Ticket, error)
	AssignTicket(ticketID uuid.UUID, petugasID uuid.UUID) (*models.Ticket, error)
	UpdateStatus(ticketID uuid.UUID, input UpdateTicketStatusInput) (*models.Ticket, error)
	DeleteTicket(ticketID uuid.UUID) error
}

type ticketService struct {
	db          *gorm.DB
	ticketRepo  repository.TicketRepository
	pelaporRepo repository.PelaporRepository
	layananRepo repository.LayananRepository
	userRepo    repository.UserRepository
}

func NewTicketService(
	db *gorm.DB,
	ticketRepo repository.TicketRepository,
	pelaporRepo repository.PelaporRepository,
	layananRepo repository.LayananRepository,
	userRepo repository.UserRepository,
) TicketService {
	return &ticketService{
		db:          db,
		ticketRepo:  ticketRepo,
		pelaporRepo: pelaporRepo,
		layananRepo: layananRepo,
		userRepo:    userRepo,
	}
}

func (s *ticketService) CreateTicket(input CreateTicketInput) (*models.Ticket, error) {
	if input.NIKPelapor == "" || input.Nama == "" || input.Instansi == "" || input.NoHP == "" {
		return nil, errors.New("data pelapor (nik, nama, instansi, no_hp) wajib diisi lengkap")
	}

	if input.IdPelayanan == uuid.Nil {
		return nil, errors.New("kategori layanan wajib dipilih")
	}

	if input.JudulKendala == "" || input.Keterangan == "" {
		return nil, errors.New("judul kendala dan keterangan rincian keluhan wajib diisi")
	}

	if len(input.FotoPaths) > 3 {
		return nil, errors.New("maksimal foto pendukung keluhan adalah 3 file")
	}

	// 1. Validasi Layanan
	layanan, err := s.layananRepo.FindByID(input.IdPelayanan)
	if err != nil || layanan == nil {
		return nil, errors.New("kategori layanan tidak valid")
	}
	if layanan.Status != models.LayananStatusAktif {
		return nil, errors.New("kategori layanan saat ini sedang nonaktif")
	}

	// 2. Simpan / Perbarui Data Pelapor
	pelapor := &models.Pelapor{
		NIK:      input.NIKPelapor,
		Nama:     input.Nama,
		Instansi: input.Instansi,
		NoHP:     input.NoHP,
	}
	if err := s.pelaporRepo.CreateOrUpdate(pelapor); err != nil {
		return nil, errors.New("gagal menyimpan profil pelapor: " + err.Error())
	}

	// 3. Generate Nomor Tiket
	nomorTicket, err := utils.GenerateTicketNumber(s.db)
	if err != nil {
		return nil, errors.New("gagal menghasilkan nomor tiket: " + err.Error())
	}

	// 4. Buat Tiket
	ticket := &models.Ticket{
		NomorTicket:  nomorTicket,
		NikPelapor:   pelapor.NIK,
		IdPelayanan:  layanan.ID,
		JudulKendala: input.JudulKendala,
		Keterangan:   input.Keterangan,
		Foto:         models.StringArray(input.FotoPaths),
		Status:       models.TicketStatusMenunggu,
	}

	if err := s.ticketRepo.Create(ticket); err != nil {
		return nil, errors.New("gagal menyimpan tiket: " + err.Error())
	}

	// Ambil data lengkap yang sudah ter-preload
	return s.ticketRepo.FindByID(ticket.ID)
}

func (s *ticketService) GetTickets(filter repository.TicketFilter) ([]models.Ticket, int64, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 || filter.PerPage > 100 {
		filter.PerPage = 10
	}
	return s.ticketRepo.FindAll(filter)
}

func (s *ticketService) GetTicketByIDOrNumber(identifier string) (*models.Ticket, error) {
	// Cek apakah identifier adalah UUID
	if uid, err := uuid.Parse(identifier); err == nil {
		return s.ticketRepo.FindByID(uid)
	}
	return s.ticketRepo.FindByNomorTicket(identifier)
}

func (s *ticketService) AssignTicket(ticketID uuid.UUID, petugasID uuid.UUID) (*models.Ticket, error) {
	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return nil, errors.New("tiket tidak ditemukan")
	}

	petugas, err := s.userRepo.FindByID(petugasID)
	if err != nil || petugas == nil {
		return nil, errors.New("petugas IT tidak valid")
	}

	if ticket.Status == models.TicketStatusSelesai || ticket.Status == models.TicketStatusDibatalkan {
		return nil, errors.New("tiket yang sudah selesai atau dibatalkan tidak dapat ditugaskan ulang")
	}

	ticket.IdPetugas = &petugasID
	ticket.Status = models.TicketStatusDiproses

	if err := s.ticketRepo.Update(ticket); err != nil {
		return nil, err
	}

	return s.ticketRepo.FindByID(ticketID)
}

func (s *ticketService) UpdateStatus(ticketID uuid.UUID, input UpdateTicketStatusInput) (*models.Ticket, error) {
	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return nil, errors.New("tiket tidak ditemukan")
	}

	validStatus := map[models.TicketStatus]bool{
		models.TicketStatusMenunggu:   true,
		models.TicketStatusDiproses:   true,
		models.TicketStatusSelesai:    true,
		models.TicketStatusDibatalkan: true,
	}

	if !validStatus[input.Status] {
		return nil, errors.New("status tiket tidak valid. Pilihan: menunggu, diproses, selesai, dibatalkan")
	}

	if input.Status == models.TicketStatusSelesai {
		if input.Solusi == nil || *input.Solusi == "" {
			return nil, errors.New("catatan solusi penanganan wajib diisi saat menyelesaikan tiket")
		}
		now := time.Now()
		ticket.FinishedAt = &now
		ticket.Solusi = input.Solusi
	} else if input.Solusi != nil {
		ticket.Solusi = input.Solusi
	}

	ticket.Status = input.Status

	if err := s.ticketRepo.Update(ticket); err != nil {
		return nil, err
	}

	return s.ticketRepo.FindByID(ticketID)
}

func (s *ticketService) DeleteTicket(ticketID uuid.UUID) error {
	_, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return errors.New("tiket tidak ditemukan")
	}
	return s.ticketRepo.Delete(ticketID)
}
