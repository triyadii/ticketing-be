package repository

import (
	"ticketing-be/internal/models"

	"gorm.io/gorm"
)

type PelaporRepository interface {
	CreateOrUpdate(pelapor *models.Pelapor) error
	FindByNIK(nik string) (*models.Pelapor, error)
	FindAll(page, perPage int, search string) ([]models.Pelapor, int64, error)
	Update(pelapor *models.Pelapor) error
	Delete(nik string) error
}

type pelaporRepository struct {
	db *gorm.DB
}

func NewPelaporRepository(db *gorm.DB) PelaporRepository {
	return &pelaporRepository{db: db}
}

func (r *pelaporRepository) CreateOrUpdate(pelapor *models.Pelapor) error {
	var existing models.Pelapor
	err := r.db.Where("nik = ?", pelapor.NIK).First(&existing).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return r.db.Create(pelapor).Error
		}
		return err
	}

	// Update data jika pelapor sudah ada
	existing.Nama = pelapor.Nama
	existing.Instansi = pelapor.Instansi
	existing.NoHP = pelapor.NoHP
	return r.db.Save(&existing).Error
}

func (r *pelaporRepository) FindByNIK(nik string) (*models.Pelapor, error) {
	var pelapor models.Pelapor
	err := r.db.Preload("Tickets", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC")
	}).Preload("Tickets.Layanan").Preload("Tickets.Petugas").Where("nik = ?", nik).First(&pelapor).Error
	if err != nil {
		return nil, err
	}
	return &pelapor, nil
}

func (r *pelaporRepository) FindAll(page, perPage int, search string) ([]models.Pelapor, int64, error) {
	var pelapors []models.Pelapor
	var total int64

	query := r.db.Model(&models.Pelapor{})

	if search != "" {
		query = query.Where("nik ILIKE ? OR nama ILIKE ? OR instansi ILIKE ? OR no_hp ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	err := query.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&pelapors).Error
	return pelapors, total, err
}

func (r *pelaporRepository) Update(pelapor *models.Pelapor) error {
	return r.db.Save(pelapor).Error
}

func (r *pelaporRepository) Delete(nik string) error {
	return r.db.Delete(&models.Pelapor{}, "nik = ?", nik).Error
}
