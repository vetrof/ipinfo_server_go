package service

import (
	"ip_info_server/internal/models"
	"ip_info_server/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) RegisterUser(username, password string) (*models.User, error) {
	return s.userRepo.CreateUser(username, password)
}

func (s *UserService) LoginUser(username, password string) (*models.User, error) {
	return s.userRepo.GetUserByCredentials(username, password)
}

func (s *UserService) ValidateToken(token string) (*models.User, error) {
	return s.userRepo.GetUserByToken(token)
}

func (s *UserService) GetUserIDByToken(token string) (int, error) {
	return s.userRepo.GetUserIDByToken(token)
}
