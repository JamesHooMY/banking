package user

//go:generate mockgen -destination ./mock/user.go -source=./user.go -package=mock

import (
	"context"
	"errors"

	"banking/model"

	"github.com/shopspring/decimal"
)

var ErrPasswordIncorrect = errors.New("password incorrect")

type IUserService interface {
	CreateUser(ctx context.Context, user *model.User) (err error)
	GetUser(ctx context.Context, userID uint) (user *model.User, err error)
	GetUsers(ctx context.Context) (users []*model.User, err error)
	Transfer(ctx context.Context, fromUserID, toUserID uint, amount decimal.Decimal) (user *model.User, err error)
	Deposit(ctx context.Context, userID uint, amount decimal.Decimal) (user *model.User, err error)
	Withdraw(ctx context.Context, userID uint, amount decimal.Decimal) (user *model.User, err error)
}

type IUserQueryRepo interface {
	GetUser(ctx context.Context, userID uint) (user *model.User, err error)
	GetUsers(ctx context.Context) (users []*model.User, err error)
}

type IUserCommandRepo interface {
	CreateUser(ctx context.Context, user *model.User) (err error)
	Transfer(ctx context.Context, fromUserID, toUserID uint, amount decimal.Decimal) (user *model.User, err error)
	Deposit(ctx context.Context, userID uint, amount decimal.Decimal) (user *model.User, err error)
	Withdraw(ctx context.Context, userID uint, amount decimal.Decimal) (user *model.User, err error)
}

type userService struct {
	userQryRepo IUserQueryRepo
	userCmdRepo IUserCommandRepo
}

// add database repo here
func NewUserService(userQryRepo IUserQueryRepo, userCmdRepo IUserCommandRepo) IUserService {
	return &userService{
		userQryRepo: userQryRepo,
		userCmdRepo: userCmdRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *model.User) (err error) {
	return s.userCmdRepo.CreateUser(ctx, user)
}

func (s *userService) GetUser(ctx context.Context, userID uint) (user *model.User, err error) {
	return s.userQryRepo.GetUser(ctx, userID)
}

func (s *userService) GetUsers(ctx context.Context) (users []*model.User, err error) {
	return s.userQryRepo.GetUsers(ctx)
}

func (s *userService) Transfer(ctx context.Context, fromUserID, toUserID uint, amount decimal.Decimal) (user *model.User, err error) {
	return s.userCmdRepo.Transfer(ctx, fromUserID, toUserID, amount)
}

func (s *userService) Deposit(ctx context.Context, userID uint, amount decimal.Decimal) (user *model.User, err error) {
	return s.userCmdRepo.Deposit(ctx, userID, amount)
}

func (s *userService) Withdraw(ctx context.Context, userID uint, amount decimal.Decimal) (user *model.User, err error) {
	return s.userCmdRepo.Withdraw(ctx, userID, amount)
}
