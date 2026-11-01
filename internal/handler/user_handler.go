package handler

import (
	"strconv"

	"ticketing-be/internal/service"
	"ticketing-be/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// CreateUser godoc
// @Summary      Tambah Petugas IT Baru
// @Description  Membuat akun petugas IT baru
// @Tags         Users (Petugas IT)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body service.CreateUserRequest true "Data Petugas Baru"
// @Success      201 {object} utils.APIResponse{data=models.User}
// @Failure      400 {object} utils.APIResponse
// @Router       /api/v1/users [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req service.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format payload tidak valid", err.Error())
	}

	user, err := h.userService.CreateUser(req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SendSuccess(c, fiber.StatusCreated, "Petugas IT berhasil ditambahkan", user)
}

// GetUsers godoc
// @Summary      Daftar Petugas IT
// @Description  Mengambil daftar seluruh akun tim IT dengan pagination dan filter
// @Tags         Users (Petugas IT)
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "Nomor halaman (default: 1)"
// @Param        per_page query int false "Item per halaman (default: 10)"
// @Param        search query string false "Cari berdasarkan nama atau username"
// @Param        status query string false "Filter status (aktif / nonaktif)"
// @Success      200 {object} utils.APIResponse{data=utils.PaginatedData}
// @Router       /api/v1/users [get]
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	search := c.Query("search", "")
	status := c.Query("status", "")

	users, total, err := h.userService.GetUsers(page, perPage, search, status)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal mengambil data user", err.Error())
	}

	return utils.SendPaginated(c, fiber.StatusOK, "Berhasil mengambil data petugas IT", users, page, perPage, total)
}

// GetUserByID godoc
// @Summary      Detail Petugas IT
// @Description  Mengambil detail satu petugas IT berdasarkan UUID
// @Tags         Users (Petugas IT)
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User UUID"
// @Success      200 {object} utils.APIResponse{data=models.User}
// @Failure      404 {object} utils.APIResponse
// @Router       /api/v1/users/{id} [get]
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	uid, err := uuid.Parse(idParam)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format UUID tidak valid", nil)
	}

	user, err := h.userService.GetUserByID(uid)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Petugas IT tidak ditemukan", nil)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Berhasil mengambil data petugas IT", user)
}

// UpdateUser godoc
// @Summary      Update Petugas IT
// @Description  Memperbarui data profil, password, atau status akun petugas IT
// @Tags         Users (Petugas IT)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User UUID"
// @Param        request body service.UpdateUserRequest true "Data Perubahan Petugas"
// @Success      200 {object} utils.APIResponse{data=models.User}
// @Failure      400 {object} utils.APIResponse
// @Router       /api/v1/users/{id} [put]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	uid, err := uuid.Parse(idParam)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format UUID tidak valid", nil)
	}

	var req service.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format payload tidak valid", err.Error())
	}

	user, err := h.userService.UpdateUser(uid, req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Data petugas IT berhasil diperbarui", user)
}

// DeleteUser godoc
// @Summary      Hapus Petugas IT
// @Description  Menghapus data akun petugas IT
// @Tags         Users (Petugas IT)
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User UUID"
// @Success      200 {object} utils.APIResponse
// @Failure      400 {object} utils.APIResponse
// @Router       /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	uid, err := uuid.Parse(idParam)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format UUID tidak valid", nil)
	}

	if err := h.userService.DeleteUser(uid); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Petugas IT berhasil dihapus", nil)
}
