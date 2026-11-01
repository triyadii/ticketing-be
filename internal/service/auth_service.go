package service

import (
	"errors"

	"ticketing-be/internal/config"
	"ticketing-be/internal/middleware"
	"ticketing-be/internal/models"
	"ticketing-be/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username" example:"admin"`
	Password string `json:"password" example:"password123"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

type AuthService interface {
	Login(req LoginRequest) (*LoginResponse, error)
	GetProfile(userID uuid.UUID) (*models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{userRepo: userRepo, cfg: cfg}
}

func (s *authService) Login(req LoginRequest) (*LoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("username dan password wajib diisi")
	}

	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("username atau password salah")
	}

	if user.Status != models.UserStatusAktif {
		return nil, errors.New("akun anda dinonaktifkan, hubungi administrator")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("username atau password salah")
	}

	token, err := middleware.GenerateToken(user, s.cfg)
	if err != nil {
		return nil, errors.New("gagal membuat token autentikasi")
	}

	return &LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *authService) GetProfile(userID uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}
