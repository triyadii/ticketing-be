package handler

import (
	"mime/multipart"
	"strconv"

	"ticketing-be/internal/config"
	"ticketing-be/internal/middleware"
	"ticketing-be/internal/repository"
	"ticketing-be/internal/service"
	"ticketing-be/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TicketHandler struct {
	ticketService service.TicketService
	cfg           *config.Config
}

func NewTicketHandler(ticketService service.TicketService, cfg *config.Config) *TicketHandler {
	return &TicketHandler{ticketService: ticketService, cfg: cfg}
}

// CreateTicket godoc
// @Summary      Buat Tiket Kendala Baru (Public)
// @Description  Pelapor mengirimkan tiket kendala baru disertai foto bukti (maksimal 3 foto)
// @Tags         Tickets
// @Accept       multipart/form-data
// @Produce      json
// @Param        nik_pelapor formData string true "NIK Pelapor (16 digit)"
// @Param        nama formData string true "Nama Lengkap Pelapor"
// @Param        instansi formData string true "OPD / Instansi / Bidang"
// @Param        no_hp formData string true "Nomor WhatsApp / HP Aktif"
// @Param        id_pelayanan formData string true "UUID Kategori Layanan"
// @Param        judul_kendala formData string true "Ringkasan Kendala (maks. 150 karakter)"
// @Param        keterangan formData string true "Rincian Detail Kendala"
// @Param        foto formData file false "Foto Bukti Kendala (maksimal 3 file, format jpg/png/webp)"
// @Success      201 {object} utils.APIResponse{data=models.Ticket}
// @Failure      400 {object} utils.APIResponse
// @Router       /api/v1/tickets [post]
func (h *TicketHandler) CreateTicket(c *fiber.Ctx) error {
	nikPelapor := c.FormValue("nik_pelapor")
	nama := c.FormValue("nama")
	instansi := c.FormValue("instansi")
	noHP := c.FormValue("no_hp")
	idPelayananStr := c.FormValue("id_pelayanan")
	judulKendala := c.FormValue("judul_kendala")
	keterangan := c.FormValue("keterangan")

	idPelayanan, err := uuid.Parse(idPelayananStr)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "ID Layanan tidak valid", nil)
	}

	// Tangani upload foto (maksimal 3 foto)
	var photoHeaders []*multipart.FileHeader
	form, err := c.MultipartForm()
	if err == nil && form != nil {
		if files, exists := form.File["foto"]; exists {
			photoHeaders = files
		}
	}

	maxSizeBytes := int64(h.cfg.MaxFileSizeMB) * 1024 * 1024
	savedPaths, err := utils.SaveUploadedPhotos(photoHeaders, h.cfg.UploadDir, maxSizeBytes)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	input := service.CreateTicketInput{
		NIKPelapor:   nikPelapor,
		Nama:         nama,
		Instansi:     instansi,
		NoHP:         noHP,
		IdPelayanan:  idPelayanan,
		JudulKendala: judulKendala,
		Keterangan:   keterangan,
		FotoPaths:    savedPaths,
	}

	ticket, err := h.ticketService.CreateTicket(input)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SendSuccess(c, fiber.StatusCreated, "Tiket kendala berhasil dibuat", ticket)
}

// GetTickets godoc
// @Summary      Daftar Tiket Kendala
// @Description  Mengambil seluruh tiket dengan filter status, layanan, petugas, pelapor, range tanggal, dan kata kunci
// @Tags         Tickets
// @Produce      json
// @Param        page query int false "Nomor halaman (default: 1)"
// @Param        per_page query int false "Item per halaman (default: 10)"
// @Param        status query string false "Filter status (menunggu, diproses, selesai, dibatalkan)"
// @Param        id_pelayanan query string false "Filter ID Layanan (UUID)"
// @Param        id_petugas query string false "Filter ID Petugas (UUID)"
// @Param        nik_pelapor query string false "Filter NIK Pelapor"
// @Param        start_date query string false "Filter tanggal awal (YYYY-MM-DD)"
// @Param        end_date query string false "Filter tanggal akhir (YYYY-MM-DD)"
// @Param        search query string false "Cari nomor tiket, judul kendala, atau keterangan"
// @Success      200 {object} utils.APIResponse{data=utils.PaginatedData}
// @Router       /api/v1/tickets [get]
func (h *TicketHandler) GetTickets(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))

	filter := repository.TicketFilter{
		Page:        page,
		PerPage:     perPage,
		Status:      c.Query("status", ""),
		IdPelayanan: c.Query("id_pelayanan", ""),
		IdPetugas:   c.Query("id_petugas", ""),
		NikPelapor:  c.Query("nik_pelapor", ""),
		StartDate:   c.Query("start_date", ""),
		EndDate:     c.Query("end_date", ""),
		Search:      c.Query("search", ""),
	}

	tickets, total, err := h.ticketService.GetTickets(filter)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal mengambil daftar tiket", err.Error())
	}

	return utils.SendPaginated(c, fiber.StatusOK, "Berhasil mengambil daftar tiket", tickets, page, perPage, total)
}

// GetTicketDetail godoc
// @Summary      Detail Tiket Kendala
// @Description  Mengambil rincian tiket berdasarkan UUID atau Nomor Tiket (misal IT-202611-0001)
// @Tags         Tickets
// @Produce      json
// @Param        identifier path string true "UUID Tiket atau Nomor Tiket (IT-YYYYMM-XXXX)"
// @Success      200 {object} utils.APIResponse{data=models.Ticket}
// @Failure      404 {object} utils.APIResponse
// @Router       /api/v1/tickets/{identifier} [get]
func (h *TicketHandler) GetTicketDetail(c *fiber.Ctx) error {
	identifier := c.Params("identifier")

	ticket, err := h.ticketService.GetTicketByIDOrNumber(identifier)
	if err != nil || ticket == nil {
		return utils.SendError(c, fiber.StatusNotFound, "Tiket tidak ditemukan", nil)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Berhasil mengambil rincian tiket", ticket)
}

// AssignTicket godoc
// @Summary      Ambil Tiket oleh Petugas IT
// @Description  Petugas IT menugaskan tiket ke diri sendiri (status berubah menjadi diproses)
// @Tags         Tickets
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "UUID Tiket"
// @Success      200 {object} utils.APIResponse{data=models.Ticket}
// @Failure      400 {object} utils.APIResponse
// @Router       /api/v1/tickets/{id}/assign [patch]
func (h *TicketHandler) AssignTicket(c *fiber.Ctx) error {
	idParam := c.Params("id")
	ticketID, err := uuid.Parse(idParam)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format UUID tiket tidak valid", nil)
	}

	claims, ok := c.Locals("user").(*middleware.JWTClaims)
	if !ok || claims == nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "Sesi autentikasi tidak valid", nil)
	}

	ticket, err := h.ticketService.AssignTicket(ticketID, claims.UserID)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Tiket berhasil ditugaskan dan berstatus diproses", ticket)
}

// UpdateTicketStatus godoc
// @Summary      Update Status Tiket
// @Description  Mengubah status tiket (menunggu, diproses, selesai, dibatalkan). Jika status selesai, wajib mengisi field solusi!
// @Tags         Tickets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "UUID Tiket"
// @Param        request body service.UpdateTicketStatusInput true "Status dan Catatan Solusi"
// @Success      200 {object} utils.APIResponse{data=models.Ticket}
// @Failure      400 {object} utils.APIResponse
// @Router       /api/v1/tickets/{id}/status [patch]
func (h *TicketHandler) UpdateTicketStatus(c *fiber.Ctx) error {
	idParam := c.Params("id")
	ticketID, err := uuid.Parse(idParam)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format UUID tiket tidak valid", nil)
	}

	var input service.UpdateTicketStatusInput
	if err := c.BodyParser(&input); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format payload tidak valid", err.Error())
	}

	ticket, err := h.ticketService.UpdateStatus(ticketID, input)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Status tiket berhasil diperbarui", ticket)
}

// DeleteTicket godoc
// @Summary      Hapus Tiket Kendala
// @Description  Menghapus tiket kendala dari sistem
// @Tags         Tickets
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "UUID Tiket"
// @Success      200 {object} utils.APIResponse
// @Failure      400 {object} utils.APIResponse
// @Router       /api/v1/tickets/{id} [delete]
func (h *TicketHandler) DeleteTicket(c *fiber.Ctx) error {
	idParam := c.Params("id")
	ticketID, err := uuid.Parse(idParam)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format UUID tiket tidak valid", nil)
	}

	if err := h.ticketService.DeleteTicket(ticketID); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Tiket berhasil dihapus", nil)
}
