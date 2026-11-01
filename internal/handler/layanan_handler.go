package handler

import (
	"ticketing-be/internal/service"
	"ticketing-be/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LayananHandler struct {
	layananService service.LayananService
}

func NewLayananHandler(layananService service.LayananService) *LayananHandler {
	return &LayananHandler{layananService: layananService}
}

// CreateLayanan godoc
// @Summary      Tambah Layanan Baru
// @Description  Membuat jenis kategori kendala / layanan IT baru
// @Tags         Layanan
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body service.LayananRequest true "Data Layanan Baru"
// @Success      201 {object} utils.APIResponse{data=models.Layanan}
// @Failure      400 {object} utils.APIResponse
// @Router       /api/v1/layanan [post]
func (h *LayananHandler) CreateLayanan(c *fiber.Ctx) error {
	var req service.LayananRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format payload tidak valid", err.Error())
	}

	layanan, err := h.layananService.CreateLayanan(req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SendSuccess(c, fiber.StatusCreated, "Layanan berhasil dibuat", layanan)
}

// GetLayanans godoc
// @Summary      Daftar Layanan IT
// @Description  Mengambil seluruh kategori layanan IT (bisa difilter status)
// @Tags         Layanan
// @Produce      json
// @Param        status query string false "Filter status (aktif / nonaktif)"
// @Success      200 {object} utils.APIResponse{data=[]models.Layanan}
// @Router       /api/v1/layanan [get]
func (h *LayananHandler) GetLayanans(c *fiber.Ctx) error {
	status := c.Query("status", "")
	layanans, err := h.layananService.GetLayanans(status)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal mengambil daftar layanan", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Berhasil mengambil daftar layanan", layanans)
}

// GetLayananByID godoc
// @Summary      Detail Layanan IT
// @Description  Mengambil detail kategori layanan IT berdasarkan ID
// @Tags         Layanan
// @Produce      json
// @Param        id path string true "Layanan UUID"
// @Success      200 {object} utils.APIResponse{data=models.Layanan}
// @Failure      404 {object} utils.APIResponse
// @Router       /api/v1/layanan/{id} [get]
func (h *LayananHandler) GetLayananByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	uid, err := uuid.Parse(idParam)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format UUID tidak valid", nil)
	}

	layanan, err := h.layananService.GetLayananByID(uid)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Layanan tidak ditemukan", nil)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Berhasil mengambil data layanan", layanan)
}

// UpdateLayanan godoc
// @Summary      Update Layanan IT
// @Description  Memperbarui nama layanan atau status aktif/nonaktif
// @Tags         Layanan
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Layanan UUID"
// @Param        request body service.LayananRequest true "Data Perubahan Layanan"
// @Success      200 {object} utils.APIResponse{data=models.Layanan}
// @Failure      400 {object} utils.APIResponse
// @Router       /api/v1/layanan/{id} [put]
func (h *LayananHandler) UpdateLayanan(c *fiber.Ctx) error {
	idParam := c.Params("id")
	uid, err := uuid.Parse(idParam)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format UUID tidak valid", nil)
	}

	var req service.LayananRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format payload tidak valid", err.Error())
	}

	layanan, err := h.layananService.UpdateLayanan(uid, req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Layanan berhasil diperbarui", layanan)
}

// DeleteLayanan godoc
// @Summary      Hapus Layanan IT
// @Description  Menghapus kategori layanan IT
// @Tags         Layanan
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Layanan UUID"
// @Success      200 {object} utils.APIResponse
// @Failure      400 {object} utils.APIResponse
// @Router       /api/v1/layanan/{id} [delete]
func (h *LayananHandler) DeleteLayanan(c *fiber.Ctx) error {
	idParam := c.Params("id")
	uid, err := uuid.Parse(idParam)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format UUID tidak valid", nil)
	}

	if err := h.layananService.DeleteLayanan(uid); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Layanan berhasil dihapus", nil)
}
