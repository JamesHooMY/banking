package user

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"banking/global"
	"banking/model"

	v1 "banking/app/api/v1"
	userRepo "banking/app/repo/mysql/user"
	userSrv "banking/app/service/user"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type IUserHandler interface {
	CreateUser() gin.HandlerFunc
	GetUser() gin.HandlerFunc
	GetUsers() gin.HandlerFunc
	Transfer() gin.HandlerFunc
	Deposit() gin.HandlerFunc
	Withdraw() gin.HandlerFunc
}

type UserHandler struct {
	UserService userSrv.IUserService
}

func NewUserHandler(userService userSrv.IUserService) IUserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

// @Tags User
// @Router /api/v1/user [post]
// @Summary Create User
// @Description Create User
// @Accept json
// @Produce json
// @Param CreateUserReq body CreateUserReq user "create user request"
// @Success 201 {object} CreateUserResp "success created user"
// @Failure 400 {object} v1.ErrResponse "bad request"
// @Failure 500 {object} v1.ErrResponse "internal server error"
func (h *UserHandler) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreateUserReq
		if err := c.ShouldBindJSON(&input); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrResponse{
				Msg: err.Error(),
			})
			return
		}

		user := &model.User{
			Name:    input.Name,
			Balance: decimal.NewFromFloat(0),
		}
		if err := h.UserService.CreateUser(c, user); err != nil {
			if errors.Is(err, userRepo.ErrUserExisted) {
				c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrResponse{
					Msg: userRepo.ErrUserExisted.Error(),
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, &v1.ErrResponse{
				Msg: err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, &CreateUserResp{
			Data: &User{
				ID:      user.ID,
				Name:    user.Name,
				Balance: user.Balance,
			},
			Msg: "user created",
		})
	}
}

type CreateUserReq struct {
	Name string `json:"name" binding:"required,min=3,max=20,alphanumunicode"`
}

type User struct {
	ID      uint            `json:"id"`
	Name    string          `json:"name"`
	Balance decimal.Decimal `json:"balance"`
}

type CreateUserResp struct {
	Data *User  `json:"data"`
	Msg  string `json:"msg"`
}

// @Tags User
// @Router /api/v1/user/{id} [get]
// @Summary Get User
// @Description Get User
// @Accept json
// @Produce json
// @Param id path int true "user id"
// @Success 200 {object} GetUserResp "success"
// @Failure 400 {object} v1.ErrResponse "bad request"
// @Failure 500 {object} v1.ErrResponse "internal server error"
func (h *UserHandler) GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		userID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrResponse{
				Msg: "invalid user id",
			})
			return
		}

		user, err := h.UserService.GetUser(c, uint(userID))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, &v1.ErrResponse{
				Msg: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, &GetUserResp{
			Data: &User{
				ID:      user.ID,
				Name:    user.Name,
				Balance: user.Balance,
			},
			Msg: "user found",
		})
	}
}

type GetUserResp struct {
	Data *User  `json:"data"`
	Msg  string `json:"msg"`
}

func (h *UserHandler) GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := h.UserService.GetUsers(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, v1.ErrResponse{
			Data: users,
			Msg:  "users found",
		})
	}
}

func (h *UserHandler) Transfer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			FromUserID uint    `json:"fromUserId" binding:"required,min=1,number"`
			ToUserID   uint    `json:"toUserId" binding:"required,min=1,number"`
			Amount     float64 `json:"amount" binding:"required,gt=0,number"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrResponse{
				Msg: err.Error(),
			})
			return
		}

		if input.FromUserID == input.ToUserID {
			c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrResponse{
				Msg: "fromUserId and toUserId should not be the same",
			})
			return
		}

		user, err := h.UserService.Transfer(c, input.FromUserID, input.ToUserID, decimal.NewFromFloat(input.Amount))
		if err != nil {
			if errors.Is(err, userRepo.ErrInsufficientBalance) {
				c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrResponse{
					Msg: userRepo.ErrInsufficientBalance.Error(),
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, &v1.ErrResponse{
				Msg: err.Error(),
			})
			return
		}

		global.Logger.Info(fmt.Sprintf("[Transfer]: fromUserId: %d, toUserId: %d, amount: %f", input.FromUserID, input.ToUserID, input.Amount))

		c.JSON(http.StatusOK, &v1.ErrResponse{
			Data: user,
			Msg:  "transfer success",
		})
	}
}

func (h *UserHandler) Deposit() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			UserID uint    `json:"userId" binding:"required,min=1,number"`
			Amount float64 `json:"amount" binding:"required,gt=0,number"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrResponse{
				Msg: err.Error(),
			})
			return
		}

		user, err := h.UserService.Deposit(c, input.UserID, decimal.NewFromFloat(input.Amount))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, &v1.ErrResponse{
				Msg: err.Error(),
			})
			return
		}

		global.Logger.Info(fmt.Sprintf("[Deposit]: userId: %d, amount: %f", input.UserID, input.Amount))

		c.JSON(http.StatusOK, &v1.ErrResponse{
			Data: user,
			Msg:  "deposit success",
		})
	}
}

func (h *UserHandler) Withdraw() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			UserID uint    `json:"userId" binding:"required,min=1,number"`
			Amount float64 `json:"amount" binding:"required,gt=0,number"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrResponse{
				Msg: err.Error(),
			})
			return
		}

		user, err := h.UserService.Withdraw(c, input.UserID, decimal.NewFromFloat(input.Amount))
		if err != nil {
			if errors.Is(err, userRepo.ErrInsufficientBalance) {
				c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrResponse{
					Msg: userRepo.ErrInsufficientBalance.Error(),
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, &v1.ErrResponse{
				Msg: err.Error(),
			})
			return
		}

		global.Logger.Info(fmt.Sprintf("[Withdraw]: userId: %d, amount: %f", input.UserID, input.Amount))

		c.JSON(http.StatusOK, &v1.ErrResponse{
			Data: user,
			Msg:  "withdraw success",
		})
	}
}
