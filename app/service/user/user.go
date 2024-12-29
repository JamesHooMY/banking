package user

import (
	"context"
	"errors"
	"time"

	"banking/domain"
	mysqlModel "banking/model/mysql"
	"banking/utils"

	"github.com/go-redis/redis/v8"
	"go.elastic.co/apm/v2"
	"golang.org/x/crypto/bcrypt"
)

var ErrPasswordIncorrect = errors.New("password incorrect")

type userService struct {
	userQryRepo       domain.IUserQueryRepo
	userCmdRepo       domain.IUserCommandRepo
	jwtRedisCmdRepo   domain.IRedisJWTCommandRepo
	jwtRedisQueryRepo domain.IRedisJWTQueryRepo
}

// add database repo here
func NewUserService(
	UserCmdRepo domain.IUserCommandRepo,
	UserQryRepo domain.IUserQueryRepo,
	JWTRedisCmdRepo domain.IRedisJWTCommandRepo,
	JWTRedisQueryRepo domain.IRedisJWTQueryRepo,
) domain.IUserService {
	return &userService{
		userQryRepo:       UserQryRepo,
		userCmdRepo:       UserCmdRepo,
		jwtRedisCmdRepo:   JWTRedisCmdRepo,
		jwtRedisQueryRepo: JWTRedisQueryRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *mysqlModel.User) (err error) {
	span, ctx := apm.StartSpan(ctx, "userService.CreateUser", "service")
	defer span.End()

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	return s.userCmdRepo.CreateUser(ctx, user)
}

func (s *userService) Login(ctx context.Context, email, password string) (token string, err error) {
	span, ctx := apm.StartSpan(ctx, "userService.Login", "service")
	defer span.End()

	// Check if token exists in redis
	tokenRedis, err := s.jwtRedisQueryRepo.GetRedisJWT(ctx, email)
	if err != redis.Nil && err != nil {
		return "", err
	}

	if tokenRedis != "" {
		return tokenRedis, nil
	}

	// Get user by email
	user, err := s.userQryRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if compareErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); compareErr != nil {
		return "", ErrPasswordIncorrect
	}

	// Generate token
	token, err = utils.GenerateJWT(user.ID, email, user.IsAdmin)
	if err != nil {
		return "", err
	}

	// Save token to redis
	if err := s.jwtRedisCmdRepo.SetRedisJWT(ctx, email, token); err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) GetUsers(ctx context.Context, userID uint) (users []*mysqlModel.User, err error) {
	span, ctx := apm.StartSpan(ctx, "userService.GetUsers", "service")
	defer span.End()

	// Simulate slow query for demo
	time.Sleep(2 * time.Millisecond)

	return s.userQryRepo.GetUsers(ctx, userID)
}
