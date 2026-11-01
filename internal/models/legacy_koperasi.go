package models

import (
	"time"
)

type LegacyKoperasi struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	KodeAnggota    string    `gorm:"type:varchar(50);index;not null" json:"kode_anggota"`
	NamaAnggota    string    `gorm:"type:varchar(100);not null" json:"nama_anggota"`
	UnitKerja      string    `gorm:"type:varchar(150)" json:"unit_kerja"`
	JenisTransaksi string    `gorm:"type:varchar(100);not null" json:"jenis_transaksi"`
	Nominal        float64   `gorm:"type:numeric(15,2);not null" json:"nominal"`
	Keterangan     string    `gorm:"type:text" json:"keterangan"`
	Tanggal        time.Time `gorm:"type:date;not null" json:"tanggal"`
	CreatedAt      time.Time `json:"created_at"`
}

func (LegacyKoperasi) TableName() string {
	return "legacy_koperasi"
}
