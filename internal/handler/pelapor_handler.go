package handler

import (
	"strconv"

	"ticketing-be/internal/service"
	"ticketing-be/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type PelaporHandler struct {
	pelaporService service.PelaporService
}

func NewPelaporHandler(pelaporService service.PelaporService) *PelaporHandler {
	return &PelaporHandler{pelaporService: pelaporService}
}

// CreatePelapor godoc
// @Summary      Pendaftaran / Tambah Pelapor
// @Description  Mendaftarkan master data pelapor pegawai/masyarakat
// @Tags         Pelapor
// @Accept       json
// @Produce      json
// @Param        request body service.PelaporRequest true "Data Pelapor"
// @Success      201 {object} utils.APIResponse{data=models.Pelapor}
// @Failure      400 {object} utils.APIResponse
// @Router       /api/v1/pelapor [post]
func (h *PelaporHandler) CreatePelapor(c *fiber.Ctx) error {
	var req service.PelaporRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format payload tidak valid", err.Error())
	}

	pelapor, err := h.pelaporService.CreateOrUpdate(req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SendSuccess(c, fiber.StatusCreated, "Data pelapor berhasil disimpan", pelapor)
}

// GetPelapors godoc
// @Summary      Daftar Pelapor
// @Description  Mengambil data seluruh pelapor dengan pagination dan pencarian
// @Tags         Pelapor
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "Nomor halaman (default: 1)"
// @Param        per_page query int false "Item per halaman (default: 10)"
// @Param        search query string false "Cari NIK, Nama, Instansi, atau No HP"
// @Success      200 {object} utils.APIResponse{data=utils.PaginatedData}
// @Router       /api/v1/pelapor [get]
func (h *PelaporHandler) GetPelapors(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	search := c.Query("search", "")

	pelapors, total, err := h.pelaporService.GetPelapors(page, perPage, search)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal mengambil data pelapor", err.Error())
	}

	return utils.SendPaginated(c, fiber.StatusOK, "Berhasil mengambil data pelapor", pelapors, page, perPage, total)
}

// GetPelaporByNIK godoc
// @Summary      Detail Pelapor & Riwayat Tiket
// @Description  Mengambil profil pelapor beserta seluruh riwayat tiket yang pernah dibuat
// @Tags         Pelapor
// @Produce      json
// @Param        nik path string true "NIK Pelapor (16 digit)"
// @Success      200 {object} utils.APIResponse{data=models.Pelapor}
// @Failure      404 {object} utils.APIResponse
// @Router       /api/v1/pelapor/{nik} [get]
func (h *PelaporHandler) GetPelaporByNIK(c *fiber.Ctx) error {
	nik := c.Params("nik")
	pelapor, err := h.pelaporService.GetByNIK(nik)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Data pelapor tidak ditemukan", nil)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Berhasil mengambil data pelapor", pelapor)
}

// UpdatePelapor godoc
// @Summary      Update Profil Pelapor
// @Description  Memperbarui nama, instansi, atau kontak nomor HP pelapor
// @Tags         Pelapor
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        nik path string true "NIK Pelapor"
// @Param        request body service.PelaporRequest true "Data Perubahan Pelapor"
// @Success      200 {object} utils.APIResponse{data=models.Pelapor}
// @Failure      400 {object} utils.APIResponse
// @Router       /api/v1/pelapor/{nik} [put]
func (h *PelaporHandler) UpdatePelapor(c *fiber.Ctx) error {
	nik := c.Params("nik")

	var req service.PelaporRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format payload tidak valid", err.Error())
	}

	pelapor, err := h.pelaporService.UpdatePelapor(nik, req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Profil pelapor berhasil diperbarui", pelapor)
}

// DeletePelapor godoc
// @Summary      Hapus Pelapor
// @Description  Menghapus master data pelapor
// @Tags         Pelapor
// @Produce      json
// @Security     BearerAuth
// @Param        nik path string true "NIK Pelapor"
// @Success      200 {object} utils.APIResponse
// @Failure      400 {object} utils.APIResponse
// @Router       /api/v1/pelapor/{nik} [delete]
func (h *PelaporHandler) DeletePelapor(c *fiber.Ctx) error {
	nik := c.Params("nik")

	if err := h.pelaporService.DeletePelapor(nik); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Data pelapor berhasil dihapus", nil)
}
