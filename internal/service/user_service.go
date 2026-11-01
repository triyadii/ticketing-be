package service

import (
	"errors"

	"ticketing-be/internal/models"
	"ticketing-be/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserRequest struct {
	Nama     string            `json:"nama" example:"Ahmad Dani"`
	Username string            `json:"username" example:"ahmaddani"`
	Password string            `json:"password" example:"rahasia123"`
	Status   models.UserStatus `json:"status" example:"aktif"`
}

type UpdateUserRequest struct {
	Nama     string            `json:"nama" example:"Ahmad Dani, S.Kom"`
	Password string            `json:"password,omitempty" example:"passwordBaru123"`
	Status   models.UserStatus `json:"status" example:"aktif"`
}

type UserService interface {
	CreateUser(req CreateUserRequest) (*models.User, error)
	GetUsers(page, perPage int, search string, status string) ([]models.User, int64, error)
	GetUserByID(id uuid.UUID) (*models.User, error)
	UpdateUser(id uuid.UUID, req UpdateUserRequest) (*models.User, error)
	DeleteUser(id uuid.UUID) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) CreateUser(req CreateUserRequest) (*models.User, error) {
	if req.Nama == "" || req.Username == "" || req.Password == "" {
		return nil, errors.New("nama, username, dan password wajib diisi")
	}

	existing, _ := s.userRepo.FindByUsername(req.Username)
	if existing != nil {
		return nil, errors.New("username sudah digunakan")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("gagal mengenkripsi password")
	}

	user := &models.User{
		Nama:     req.Nama,
		Username: req.Username,
		Password: string(hashedPassword),
		Status:   req.Status,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUsers(page, perPage int, search string, status string) ([]models.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}
	return s.userRepo.FindAll(page, perPage, search, status)
}

func (s *userService) GetUserByID(id uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *userService) UpdateUser(id uuid.UUID, req UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("petugas IT tidak ditemukan")
	}

	if req.Nama != "" {
		user.Nama = req.Nama
	}
	if req.Status != "" {
		user.Status = req.Status
	}
	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("gagal mengenkripsi password baru")
		}
		user.Password = string(hashed)
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) DeleteUser(id uuid.UUID) error {
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("petugas IT tidak ditemukan")
	}
	return s.userRepo.Delete(id)
}
