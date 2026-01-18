package services

import (
	"context"
	"easygin/internal/models"
	"easygin/pkg/apperror"
	"easygin/pkg/ent"
	"easygin/pkg/logging"
	"easygin/pkg/prom"

	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

type UserService struct {
	db *ent.Client
}

func NewUserService(db *ent.Client) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) CreateUser(
	ctx context.Context,
	user *models.CreateUserRequest,
) (*models.CreateUserResponse, error) {
	logger := logging.GetLoggerFromContext(ctx)
	logger.Info("Creating user", zap.String("name", user.Name))
	u, err := s.db.User.Create().SetName(user.Name).Save(ctx)
	if err != nil {
		return nil, apperror.InvalidRequest("failed to create user", err)
	}
	resp := &models.CreateUserResponse{
		ID: u.ID,
	}
	return resp, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) (*models.GetAllUsersResponse, error) {
	logger := logging.GetLoggerFromContext(ctx)
	logger.Info("Getting all users")
	dbUsers, err := s.db.User.Query().All(ctx)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			prom.GetInstance().IncrementCounterValue(prom.DatabaseErrorTotal, []string{"custom_error_type"})
			logger.Error("MySQL error occurred", zap.String("error_message", mysqlErr.Message))
			return nil, apperror.InternalError("failed to get all users", mysqlErr)
		}
		logger.Error("Error occurred while querying users", zap.Error(err))
		return nil, apperror.InvalidRequest("failed to get all users", err)
	}
	users := []*models.User{}
	for _, u := range dbUsers {
		users = append(users, &models.User{
			ID:   u.ID,
			Name: u.Name,
		})
	}
	resp := &models.GetAllUsersResponse{
		Users: users,
	}
	return resp, nil
}
