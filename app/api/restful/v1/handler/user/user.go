package user

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	v1 "banking/app/api/restful/v1"
	userRepo "banking/app/repo/mysql/user"
	"banking/domain"
	mysqlModel "banking/model/mysql"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go.elastic.co/apm/v2"
)

type UserHandler struct {
	userService   domain.IUserService
	apiKeyService domain.IAPIKeyService
}

func NewUserHandler(UserService domain.IUserService, APIKeyService domain.IAPIKeyService) domain.IUserHandler {
	return &UserHandler{
		userService:   UserService,
		apiKeyService: APIKeyService,
	}
}

// @Tags User
// @Router /api/v1/user/register [post]
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
			Name:     input.Name,
			Email:    input.Email,
			Balance:  decimal.NewFromFloat(0),
			Password: input.Password,
		}
		if err := h.userService.CreateUser(ctx, user); err != nil {
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
				Email:   user.Email,
				Balance: user.Balance,
			},
		})
	}
}

// @Tags User
// @Router /api/v1/user/login [post]
// @Summary Login
// @Description Login
// @Accept json
// @Produce json
// @Param LoginReq body LoginReq true "login request"
// @Success 200 {object} LoginResp "success"
// @Failure 400 {object} v1.ErrResponse "bad request"
// @Failure 401 {object} v1.ErrResponse "unauthorized"
func (h *UserHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input LoginReq

		// Bind JSON input
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid input"})
			return
		}

		token, err := h.userService.Login(c.Request.Context(), input.Email, input.Password)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Invalid credentials"})
			return
		}

		// Return JWT token
		c.JSON(http.StatusOK, &LoginResp{
			Token: token,
		})
	}
}

// @Tags User
// @Router /api/v1/user/{userId} [get]
// @Summary Get Users
// @Description Get Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path uint true "user id"
// @Success 200 {object} GetUsersResp "success"
// @Failure 400 {object} v1.ErrResponse "bad request"
// @Failure 500 {object} v1.ErrResponse "internal server error"
func (h *UserHandler) GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := apm.StartSpan(c.Request.Context(), "UserHandler.GetUsers", "handler")
		defer span.End()

		userId := c.Param("userId")
		authedUserId := c.GetUint("authedUserId")

		userIdUint, err := strconv.ParseUint(userId, 10, 64)
		if err != nil {
			apm.CaptureError(ctx, err).Send()
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		if !c.GetBool("isAdmin") && authedUserId != uint(userIdUint) {
			apm.CaptureError(ctx, fmt.Errorf("unauthorized")).Send()
			c.AbortWithStatusJSON(http.StatusUnauthorized, &v1.ErrResponse{
				Msg: "unauthorized",
			})
			return
		}

		users, err := h.userService.GetUsers(ctx, uint(userIdUint))
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

// @Tags User
// @Router /api/v1/user/apikey [post]
// @Summary Create API Key
// @Description Create API Key
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} CreateAPIKeyResp "success created api key"
// @Failure 500 {object} v1.ErrResponse "internal server error"
func (h *UserHandler) CreateAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := apm.StartSpan(c.Request.Context(), "UserHandler.CreateAPIKey", "handler")
		defer span.End()

		authedUserId := c.GetUint("authedUserId")
		apiKey, secretKey, err := h.apiKeyService.CreateAPIKey(ctx, authedUserId)
		if err != nil {
			apm.CaptureError(ctx, err).Send()
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, CreateAPIKeyResp{
			Data: &APIKey{
				Key:    apiKey,
				Secret: secretKey,
				UserID: authedUserId,
			},
		})
	}
}

func (h *UserHandler) DeleteAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := apm.StartSpan(c.Request.Context(), "UserHandler.DeleteAPIKey", "handler")
		defer span.End()

		input := struct {
			Key string `json:"key" binding:"required"`
		}{}

		if err := c.ShouldBindJSON(&input); err != nil {
			apm.CaptureError(ctx, err).Send()
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		authedUserId := c.GetUint("authedUserId")
		if err := h.apiKeyService.DeleteAPIKey(ctx, authedUserId, input.Key); err != nil {
			apm.CaptureError(ctx, err).Send()
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"msg": "API key deleted"})
	}
}

func (h *UserHandler) GetAPIKeys() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := apm.StartSpan(c.Request.Context(), "UserHandler.GetAPIKeys", "handler")
		defer span.End()

		userId := c.Query("userId")
		fmt.Printf("userId: %s\n", userId)
		userIdUint, err := strconv.ParseUint(userId, 10, 64)
		if err != nil {
			apm.CaptureError(ctx, err).Send()
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		isAdmin := c.GetBool("isAdmin")
		authedUserId := c.GetUint("authedUserId")

		if !isAdmin && authedUserId != uint(userIdUint) {
			apm.CaptureError(ctx, fmt.Errorf("unauthorized")).Send()
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		key, err := url.QueryUnescape(c.Query("key"))
		if err != nil {
			apm.CaptureError(ctx, err).Send()
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid key"})
			return
		}

		apiKeys, err := h.apiKeyService.GetAPIKeys(ctx, uint(userIdUint), key)
		if err != nil {
			apm.CaptureError(ctx, err).Send()
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		data := make([]*APIKey, 0, len(apiKeys))
		for _, apiKey := range apiKeys {
			data = append(data, &APIKey{
				Key:    apiKey.APIKey,
				UserID: apiKey.UserID,
			})
		}

		c.JSON(http.StatusOK, gin.H{"data": data})
	}
}
