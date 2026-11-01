package handler

import (
	"ticketing-be/internal/middleware"
	"ticketing-be/internal/service"
	"ticketing-be/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login godoc
// @Summary      Login Petugas IT
// @Description  Autentikasi akun petugas IT untuk mendapatkan access token JWT
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body service.LoginRequest true "Kredensial Login"
// @Success      200 {object} utils.APIResponse{data=service.LoginResponse}
// @Failure      400 {object} utils.APIResponse
// @Failure      401 {object} utils.APIResponse
// @Router       /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req service.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format payload request tidak valid", err.Error())
	}

	res, err := h.authService.Login(req)
	if err != nil {
		return utils.SendError(c, fiber.StatusUnauthorized, err.Error(), nil)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Login berhasil", res)
}

// GetProfile godoc
// @Summary      Profil Petugas IT yang Sedang Login
// @Description  Mengambil data detail akun petugas IT berdasarkan token JWT
// @Tags         Auth
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} utils.APIResponse{data=models.User}
// @Failure      401 {object} utils.APIResponse
// @Router       /api/v1/auth/profile [get]
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(*middleware.JWTClaims)
	if !ok || claims == nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "Sesi autentikasi tidak valid", nil)
	}

	user, err := h.authService.GetProfile(claims.UserID)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Data profil petugas tidak ditemukan", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Berhasil mengambil profil petugas", user)
}
