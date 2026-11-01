package repository

import (
	"ticketing-be/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LayananRepository interface {
	Create(layanan *models.Layanan) error
	FindByID(id uuid.UUID) (*models.Layanan, error)
	FindAll(status string) ([]models.Layanan, error)
	Update(layanan *models.Layanan) error
	Delete(id uuid.UUID) error
}

type layananRepository struct {
	db *gorm.DB
}

func NewLayananRepository(db *gorm.DB) LayananRepository {
	return &layananRepository{db: db}
}

func (r *layananRepository) Create(layanan *models.Layanan) error {
	return r.db.Create(layanan).Error
}

func (r *layananRepository) FindByID(id uuid.UUID) (*models.Layanan, error) {
	var layanan models.Layanan
	err := r.db.First(&layanan, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &layanan, nil
}

func (r *layananRepository) FindAll(status string) ([]models.Layanan, error) {
	var layanans []models.Layanan
	query := r.db.Model(&models.Layanan{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Order("nama_layanan ASC").Find(&layanans).Error
	return layanans, err
}

func (r *layananRepository) Update(layanan *models.Layanan) error {
	return r.db.Save(layanan).Error
}

func (r *layananRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Layanan{}, "id = ?", id).Error
}
