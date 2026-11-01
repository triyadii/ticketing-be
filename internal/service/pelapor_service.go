package service

import (
	"errors"

	"ticketing-be/internal/models"
	"ticketing-be/internal/repository"
)

type PelaporRequest struct {
	NIK      string `json:"nik" example:"1207012345670001"`
	Nama     string `json:"nama" example:"Rahmat Hidayat"`
	Instansi string `json:"instansi" example:"Dinas Kependudukan dan Catatan Sipil"`
	NoHP     string `json:"no_hp" example:"081234567890"`
}

type PelaporService interface {
	CreateOrUpdate(req PelaporRequest) (*models.Pelapor, error)
	GetByNIK(nik string) (*models.Pelapor, error)
	GetPelapors(page, perPage int, search string) ([]models.Pelapor, int64, error)
	UpdatePelapor(nik string, req PelaporRequest) (*models.Pelapor, error)
	DeletePelapor(nik string) error
}

type pelaporService struct {
	pelaporRepo repository.PelaporRepository
}

func NewPelaporService(pelaporRepo repository.PelaporRepository) PelaporService {
	return &pelaporService{pelaporRepo: pelaporRepo}
}

func (s *pelaporService) CreateOrUpdate(req PelaporRequest) (*models.Pelapor, error) {
	if req.NIK == "" || req.Nama == "" || req.Instansi == "" || req.NoHP == "" {
		return nil, errors.New("nik, nama, instansi, dan no_hp wajib diisi lengkap")
	}

	if len(req.NIK) < 10 || len(req.NIK) > 16 {
		return nil, errors.New("panjang NIK harus antara 10 hingga 16 karakter")
	}

	pelapor := &models.Pelapor{
		NIK:      req.NIK,
		Nama:     req.Nama,
		Instansi: req.Instansi,
		NoHP:     req.NoHP,
	}

	if err := s.pelaporRepo.CreateOrUpdate(pelapor); err != nil {
		return nil, err
	}

	return pelapor, nil
}

func (s *pelaporService) GetByNIK(nik string) (*models.Pelapor, error) {
	pelapor, err := s.pelaporRepo.FindByNIK(nik)
	if err != nil {
		return nil, errors.New("data pelapor tidak ditemukan")
	}
	return pelapor, nil
}

func (s *pelaporService) GetPelapors(page, perPage int, search string) ([]models.Pelapor, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}
	return s.pelaporRepo.FindAll(page, perPage, search)
}

func (s *pelaporService) UpdatePelapor(nik string, req PelaporRequest) (*models.Pelapor, error) {
	pelapor, err := s.pelaporRepo.FindByNIK(nik)
	if err != nil {
		return nil, errors.New("data pelapor tidak ditemukan")
	}

	if req.Nama != "" {
		pelapor.Nama = req.Nama
	}
	if req.Instansi != "" {
		pelapor.Instansi = req.Instansi
	}
	if req.NoHP != "" {
		pelapor.NoHP = req.NoHP
	}

	if err := s.pelaporRepo.Update(pelapor); err != nil {
		return nil, err
	}

	return pelapor, nil
}

func (s *pelaporService) DeletePelapor(nik string) error {
	_, err := s.pelaporRepo.FindByNIK(nik)
	if err != nil {
		return errors.New("data pelapor tidak ditemukan")
	}
	return s.pelaporRepo.Delete(nik)
}
