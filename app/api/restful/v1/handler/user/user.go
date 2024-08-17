package user

import (
	"errors"
	"net/http"
	"strconv"

	mysqlModel "banking/model/mysql"

	v1 "banking/app/api/restful/v1"
	userRepo "banking/app/repo/mysql/user"
	domain "banking/domain"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go.elastic.co/apm/v2"
)

type UserHandler struct {
	UserService domain.IUserService
}

func NewUserHandler(userService domain.IUserService) domain.IUserHandler {
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
		span, ctx := apm.StartSpan(c.Request.Context(), "UserHandler.CreateUser", "handler")
		defer span.End()

		var input CreateUserReq
		if err := c.ShouldBindJSON(&input); err != nil {
			apm.CaptureError(ctx, err).Send()
			c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrResponse{
				Msg: err.Error(),
			})
			return
		}

		user := &mysqlModel.User{
			Name:    input.Name,
			Balance: decimal.NewFromFloat(0),
		}
		if err := h.UserService.CreateUser(ctx, user); err != nil {
			if errors.Is(err, userRepo.ErrUserExisted) {
				apm.CaptureError(ctx, err).Send()
				c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrResponse{
					Msg: userRepo.ErrUserExisted.Error(),
				})
				return
			}

			apm.CaptureError(ctx, err).Send()
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
		})
	}
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
		span, ctx := apm.StartSpan(c.Request.Context(), "UserHandler.GetUser", "handler")
		defer span.End()

		id := c.Param("id")

		userID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			apm.CaptureError(ctx, err).Send()
			c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrResponse{
				Msg: "invalid user id",
			})
			return
		}

		user, err := h.UserService.GetUser(ctx, uint(userID))
		if err != nil {
			apm.CaptureError(ctx, err).Send()
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
		})
	}
}

func (h *UserHandler) GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := apm.StartSpan(c.Request.Context(), "UserHandler.GetUsers", "handler")
		defer span.End()

		users, err := h.UserService.GetUsers(ctx)
		if err != nil {
			apm.CaptureError(ctx, err).Send()
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		data := make([]*User, 0, len(users))
		for _, user := range users {
			data = append(data, &User{
				ID:      user.ID,
				Name:    user.Name,
				Balance: user.Balance,
			})
		}

		c.JSON(http.StatusOK, GetUsersResp{
			Data: data,
		})
	}
}
