package service

import (
	"bytes"
	"encoding/csv"
	"fmt"

	"ticketing-be/internal/models"
	"ticketing-be/internal/repository"

	"github.com/xuri/excelize/v2"
)

type LegacyService interface {
	GetLegacyData(filters repository.LegacyFilter) ([]models.LegacyKoperasi, int64, error)
	ExportToExcel(filters repository.LegacyFilter) (*bytes.Buffer, error)
	ExportToCSV(filters repository.LegacyFilter) (*bytes.Buffer, error)
}

type legacyService struct {
	legacyRepo repository.LegacyRepository
}

func NewLegacyService(legacyRepo repository.LegacyRepository) LegacyService {
	return &legacyService{legacyRepo: legacyRepo}
}

func (s *legacyService) GetLegacyData(filters repository.LegacyFilter) ([]models.LegacyKoperasi, int64, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PerPage < 1 || filters.PerPage > 100 {
		filters.PerPage = 10
	}
	return s.legacyRepo.FindAll(filters)
}

func (s *legacyService) ExportToExcel(filters repository.LegacyFilter) (*bytes.Buffer, error) {
	data, err := s.legacyRepo.FindAllForExport(filters)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data untuk export: %w", err)
	}

	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Data Lama Koperasi"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Header
	headers := []string{"No", "ID", "Kode Anggota", "Nama Anggota", "Unit Kerja", "Jenis Transaksi", "Nominal (Rp)", "Keterangan", "Tanggal Transaksi"}
	for colIdx, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// Style Header
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#1F4E79"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	f.SetRowStyle(sheetName, 1, 1, headerStyle)
	f.SetRowHeight(sheetName, 1, 26)

	// Isi Data
	for rowIdx, item := range data {
		row := rowIdx + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), rowIdx+1)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), item.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), item.KodeAnggota)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), item.NamaAnggota)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), item.UnitKerja)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), item.JenisTransaksi)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), item.Nominal)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), item.Keterangan)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), item.Tanggal.Format("2006-01-02"))
	}

	// Auto fit column width
	for colIdx := 1; colIdx <= len(headers); colIdx++ {
		colName, _ := excelize.ColumnNumberToName(colIdx)
		f.SetColWidth(sheetName, colName, colName, 18)
	}

	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		return nil, fmt.Errorf("gagal menulis buffer excel: %w", err)
	}

	return buf, nil
}

func (s *legacyService) ExportToCSV(filters repository.LegacyFilter) (*bytes.Buffer, error) {
	data, err := s.legacyRepo.FindAllForExport(filters)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data untuk export: %w", err)
	}

	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	// Tulis Header
	headers := []string{"No", "ID", "Kode Anggota", "Nama Anggota", "Unit Kerja", "Jenis Transaksi", "Nominal", "Keterangan", "Tanggal"}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	// Tulis Rows
	for idx, item := range data {
		record := []string{
			fmt.Sprintf("%d", idx+1),
			fmt.Sprintf("%d", item.ID),
			item.KodeAnggota,
			item.NamaAnggota,
			item.UnitKerja,
			item.JenisTransaksi,
			fmt.Sprintf("%.2f", item.Nominal),
			item.Keterangan,
			item.Tanggal.Format("2006-01-02"),
		}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf, nil
}
