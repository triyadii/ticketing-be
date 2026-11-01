package repository

import (
	"ticketing-be/internal/models"

	"gorm.io/gorm"
)

type LegacyFilter struct {
	Page           int
	PerPage        int
	Search         string
	JenisTransaksi string
	StartDate      string
	EndDate        string
}

type LegacyRepository interface {
	FindAll(filters LegacyFilter) ([]models.LegacyKoperasi, int64, error)
	FindAllForExport(filters LegacyFilter) ([]models.LegacyKoperasi, error)
	FindByID(id uint) (*models.LegacyKoperasi, error)
	Create(data *models.LegacyKoperasi) error
}

type legacyRepository struct {
	db *gorm.DB
}

func NewLegacyRepository(db *gorm.DB) LegacyRepository {
	return &legacyRepository{db: db}
}

func (r *legacyRepository) buildQuery(filters LegacyFilter) *gorm.DB {
	query := r.db.Model(&models.LegacyKoperasi{})

	if filters.Search != "" {
		pattern := "%" + filters.Search + "%"
		query = query.Where("kode_anggota ILIKE ? OR nama_anggota ILIKE ? OR unit_kerja ILIKE ? OR keterangan ILIKE ?",
			pattern, pattern, pattern, pattern)
	}

	if filters.JenisTransaksi != "" {
		query = query.Where("jenis_transaksi = ?", filters.JenisTransaksi)
	}

	if filters.StartDate != "" {
		query = query.Where("tanggal >= ?", filters.StartDate)
	}

	if filters.EndDate != "" {
		query = query.Where("tanggal <= ?", filters.EndDate)
	}

	return query
}

func (r *legacyRepository) FindAll(filters LegacyFilter) ([]models.LegacyKoperasi, int64, error) {
	var results []models.LegacyKoperasi
	var total int64

	query := r.buildQuery(filters)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filters.Page - 1) * filters.PerPage
	err := query.Order("tanggal DESC, id DESC").Offset(offset).Limit(filters.PerPage).Find(&results).Error
	return results, total, err
}

func (r *legacyRepository) FindAllForExport(filters LegacyFilter) ([]models.LegacyKoperasi, error) {
	var results []models.LegacyKoperasi
	query := r.buildQuery(filters)
	err := query.Order("tanggal DESC, id DESC").Find(&results).Error
	return results, err
}

func (r *legacyRepository) FindByID(id uint) (*models.LegacyKoperasi, error) {
	var data models.LegacyKoperasi
	err := r.db.First(&data, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *legacyRepository) Create(data *models.LegacyKoperasi) error {
	return r.db.Create(data).Error
}
