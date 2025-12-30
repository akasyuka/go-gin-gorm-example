package service

import (
	"github.com/akasyuka/service-a/model"
	"github.com/akasyuka/service-a/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(id int64) (*model.User, error) {
	return s.repo.FindByID(id)
}
