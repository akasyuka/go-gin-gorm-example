package controller

import (
	"context"

	userv1 "github.com/akasyuka/service-a/gen/user/v1"
	"github.com/akasyuka/service-a/service"
)

type UserGrpcController struct {
	userv1.UnimplementedUserServiceServer
	userService *service.UserService
}

func NewUserGrpcController(s *service.UserService) *UserGrpcController {
	return &UserGrpcController{userService: s}
}

func (c *UserGrpcController) GetUser(
	ctx context.Context,
	req *userv1.GetUserRequest,
) (*userv1.GetUserResponse, error) {

	user, err := c.userService.GetUser(req.Id)
	if err != nil {
		return nil, err
	}

	return &userv1.GetUserResponse{
		Id:     user.ID,
		Email:  user.Email,
		Active: user.Active,
	}, nil
}
