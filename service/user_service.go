package service

import (
	"bytes"
	"fmt"
	"github.com/akasyuka/go-gin-gorm-example/model"
	"github.com/akasyuka/go-gin-gorm-example/repository"
	"io"
	"net/http"
	"time"
)

type UserService interface {
	CreateUser(email, name string) (*model.User, error)
	GetUsers() ([]model.User, error)
	SendUser(u model.User)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(email, name string) (*model.User, error) {
	user := &model.User{
		Email: email,
		Name:  name,
	}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUsers() ([]model.User, error) {
	return s.repo.FindAll()
}

func (s *userService) SendUser(u model.User) {
	url := "https://example.com/api"
	payload := []byte(`{"name":"John","email":"john@example.com"}`)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer your-token") // если нужен токен

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Status:", resp.Status)
	fmt.Println("Body:", string(body))
}
