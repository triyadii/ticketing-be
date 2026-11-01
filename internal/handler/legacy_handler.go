package handler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"ticketing-be/internal/repository"
	"ticketing-be/internal/service"
	"ticketing-be/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type LegacyHandler struct {
	legacyService service.LegacyService
}

func NewLegacyHandler(legacyService service.LegacyService) *LegacyHandler {
	return &LegacyHandler{legacyService: legacyService}
}

// GetLegacyData godoc
// @Summary      Menampilkan Data Lama Koperasi
// @Description  Mengambil data lama koperasi dengan pagination, search keyword, filter jenis transaksi, dan range tanggal
// @Tags         Koperasi (Data Lama)
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "Nomor halaman (default: 1)"
// @Param        per_page query int false "Item per halaman (default: 10)"
// @Param        search query string false "Cari kode anggota, nama, unit kerja, keterangan"
// @Param        jenis_transaksi query string false "Filter jenis transaksi (misal: Simpanan Pokok, Pinjaman Reguler)"
// @Param        start_date query string false "Filter tanggal awal (YYYY-MM-DD)"
// @Param        end_date query string false "Filter tanggal akhir (YYYY-MM-DD)"
// @Success      200 {object} utils.APIResponse{data=utils.PaginatedData}
// @Router       /api/v1/koperasi/data-lama [get]
func (h *LegacyHandler) GetLegacyData(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))

	filter := repository.LegacyFilter{
		Page:           page,
		PerPage:        perPage,
		Search:         c.Query("search", ""),
		JenisTransaksi: c.Query("jenis_transaksi", ""),
		StartDate:      c.Query("start_date", ""),
		EndDate:        c.Query("end_date", ""),
	}

	data, total, err := h.legacyService.GetLegacyData(filter)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal mengambil data lama koperasi", err.Error())
	}

	return utils.SendPaginated(c, fiber.StatusOK, "Berhasil mengambil data lama koperasi", data, page, perPage, total)
}

// ExportLegacyData godoc
// @Summary      Tarik Data Lama (Export Excel / CSV)
// @Description  Mengunduh seluruh data lama koperasi ke berkas spreadsheet Excel (.xlsx) atau CSV
// @Tags         Koperasi (Data Lama)
// @Produce      application/vnd.openxmlformats-officedocument.spreadsheetml.sheet,text/csv
// @Security     BearerAuth
// @Param        format query string false "Format file export: excel atau csv (default: excel)" Enums(excel, csv)
// @Param        search query string false "Cari kata kunci"
// @Param        jenis_transaksi query string false "Filter jenis transaksi"
// @Param        start_date query string false "Filter tanggal awal (YYYY-MM-DD)"
// @Param        end_date query string false "Filter tanggal akhir (YYYY-MM-DD)"
// @Success      200 {file} binary "File Excel atau CSV hasil export"
// @Failure      500 {object} utils.APIResponse
// @Router       /api/v1/koperasi/data-lama/export [get]
func (h *LegacyHandler) ExportLegacyData(c *fiber.Ctx) error {
	format := strings.ToLower(c.Query("format", "excel"))

	filter := repository.LegacyFilter{
		Search:         c.Query("search", ""),
		JenisTransaksi: c.Query("jenis_transaksi", ""),
		StartDate:      c.Query("start_date", ""),
		EndDate:        c.Query("end_date", ""),
	}

	timestamp := time.Now().Format("20060102_150405")

	if format == "csv" {
		buf, err := h.legacyService.ExportToCSV(filter)
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "Gagal mengexport data ke CSV", err.Error())
		}

		filename := fmt.Sprintf("data_lama_koperasi_%s.csv", timestamp)
		c.Set("Content-Type", "text/csv; charset=utf-8")
		c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		return c.Send(buf.Bytes())
	}

	// Default: Excel .xlsx
	buf, err := h.legacyService.ExportToExcel(filter)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal mengexport data ke Excel", err.Error())
	}

	filename := fmt.Sprintf("data_lama_koperasi_%s.xlsx", timestamp)
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	return c.Send(buf.Bytes())
}
