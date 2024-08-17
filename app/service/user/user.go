package user

import (
	"context"
	"errors"
	"time"

	"banking/domain"
	mysqlModel "banking/model/mysql"

	"go.elastic.co/apm/v2"
)

var ErrPasswordIncorrect = errors.New("password incorrect")

type userService struct {
	userQryRepo domain.IUserQueryRepo
	userCmdRepo domain.IUserCommandRepo
}

// add database repo here
func NewUserService(userCmdRepo domain.IUserCommandRepo, userQryRepo domain.IUserQueryRepo) domain.IUserService {
	return &userService{
		userQryRepo: userQryRepo,
		userCmdRepo: userCmdRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *mysqlModel.User) (err error) {
	span, ctx := apm.StartSpan(ctx, "userService.CreateUser", "service")
	defer span.End()

	return s.userCmdRepo.CreateUser(ctx, user)
}

func (s *userService) GetUser(ctx context.Context, userID uint) (user *mysqlModel.User, err error) {
	span, ctx := apm.StartSpan(ctx, "userService.GetUser", "service")
	defer span.End()

	return s.userQryRepo.GetUser(ctx, userID)
}

func (s *userService) GetUsers(ctx context.Context) (users []*mysqlModel.User, err error) {
	span, ctx := apm.StartSpan(ctx, "userService.GetUsers", "service")
	defer span.End()

	// Simulate slow query for demo
	time.Sleep(2 * time.Millisecond)

	return s.userQryRepo.GetUsers(ctx)
}
