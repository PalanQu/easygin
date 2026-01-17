package controllers

import (
	"context"
	"easygin/internal/models"
	"easygin/internal/services"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (ctrl *UserController) GetUsers(
	ctx context.Context, _ struct{},
) (*models.GetAllUsersResponse, error) {
	usersResponse, err := ctrl.userService.GetAllUsers(ctx)
	return usersResponse, err
}

func (ctrl *UserController) CreateUser(
	ctx context.Context,
	req *models.CreateUserRequest,
) (*models.CreateUserResponse, error) {
	resp, err := ctrl.userService.CreateUser(ctx, req)
	return resp, err
}
