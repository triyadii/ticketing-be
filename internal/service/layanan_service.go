package service

import (
	"errors"

	"ticketing-be/internal/models"
	"ticketing-be/internal/repository"

	"github.com/google/uuid"
)

type LayananRequest struct {
	NamaLayanan string               `json:"nama_layanan" example:"Jaringan & WiFi"`
	Status      models.LayananStatus `json:"status" example:"aktif"`
}

type LayananService interface {
	CreateLayanan(req LayananRequest) (*models.Layanan, error)
	GetLayanans(status string) ([]models.Layanan, error)
	GetLayananByID(id uuid.UUID) (*models.Layanan, error)
	UpdateLayanan(id uuid.UUID, req LayananRequest) (*models.Layanan, error)
	DeleteLayanan(id uuid.UUID) error
}

type layananService struct {
	layananRepo repository.LayananRepository
}

func NewLayananService(layananRepo repository.LayananRepository) LayananService {
	return &layananService{layananRepo: layananRepo}
}

func (s *layananService) CreateLayanan(req LayananRequest) (*models.Layanan, error) {
	if req.NamaLayanan == "" {
		return nil, errors.New("nama layanan wajib diisi")
	}

	layanan := &models.Layanan{
		NamaLayanan: req.NamaLayanan,
		Status:      req.Status,
	}

	if err := s.layananRepo.Create(layanan); err != nil {
		return nil, err
	}

	return layanan, nil
}

func (s *layananService) GetLayanans(status string) ([]models.Layanan, error) {
	return s.layananRepo.FindAll(status)
}

func (s *layananService) GetLayananByID(id uuid.UUID) (*models.Layanan, error) {
	return s.layananRepo.FindByID(id)
}

func (s *layananService) UpdateLayanan(id uuid.UUID, req LayananRequest) (*models.Layanan, error) {
	layanan, err := s.layananRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("layanan tidak ditemukan")
	}

	if req.NamaLayanan != "" {
		layanan.NamaLayanan = req.NamaLayanan
	}
	if req.Status != "" {
		layanan.Status = req.Status
	}

	if err := s.layananRepo.Update(layanan); err != nil {
		return nil, err
	}

	return layanan, nil
}

func (s *layananService) DeleteLayanan(id uuid.UUID) error {
	_, err := s.layananRepo.FindByID(id)
	if err != nil {
		return errors.New("layanan tidak ditemukan")
	}
	return s.layananRepo.Delete(id)
}
