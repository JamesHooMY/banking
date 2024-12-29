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
				Balance: user.Balance,
			},
		})
	}
}

func (h *UserHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required,min=8,max=20"`
		}

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
		c.JSON(http.StatusOK, gin.H{"token": token})
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
// func (h *UserHandler) GetUser() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		span, ctx := apm.StartSpan(c.Request.Context(), "UserHandler.GetUser", "handler")
// 		defer span.End()

// 		id := c.Param("id")

// 		userID, err := strconv.ParseUint(id, 10, 64)
// 		if err != nil {
// 			apm.CaptureError(ctx, err).Send()
// 			c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrResponse{
// 				Msg: "invalid user id",
// 			})
// 			return
// 		}

// 		authUserId := c.GetUint("authUserId")
// 		if uint(userID) != authUserId {
// 			apm.CaptureError(ctx, fmt.Errorf("unauthorized")).Send()
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, &v1.ErrResponse{
// 				Msg: "unauthorized",
// 			})
// 			return
// 		}

// 		user, err := h.UserService.GetUser(ctx, uint(userID))
// 		if err != nil {
// 			apm.CaptureError(ctx, err).Send()
// 			c.AbortWithStatusJSON(http.StatusInternalServerError, &v1.ErrResponse{
// 				Msg: err.Error(),
// 			})
// 			return
// 		}

// 		c.JSON(http.StatusOK, &GetUserResp{
// 			Data: &User{
// 				ID:      user.ID,
// 				Name:    user.Name,
// 				Email:   user.Email,
// 				Balance: user.Balance,
// 			},
// 		})
// 	}
// }

func (h *UserHandler) GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := apm.StartSpan(c.Request.Context(), "UserHandler.GetUsers", "handler")
		defer span.End()

		userId := c.Param("userId")
		authUserId := c.GetUint("authUserId")

		idUint, err := strconv.ParseUint(userId, 10, 64)
		if err != nil {
			apm.CaptureError(ctx, err).Send()
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		if !c.GetBool("isAdmin") && authUserId != uint(idUint) {
			apm.CaptureError(ctx, fmt.Errorf("unauthorized")).Send()
			c.AbortWithStatusJSON(http.StatusUnauthorized, &v1.ErrResponse{
				Msg: "unauthorized",
			})
			return
		}

		users, err := h.userService.GetUsers(ctx, uint(idUint))
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

func (h *UserHandler) CreateAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := apm.StartSpan(c.Request.Context(), "UserHandler.CreateAPIKey", "handler")
		defer span.End()

		authUserId := c.GetUint("authUserId")
		apiKey, secretKey, err := h.apiKeyService.CreateAPIKey(ctx, authUserId)
		if err != nil {
			apm.CaptureError(ctx, err).Send()
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, CreateAPIKeyResp{
			Data: &APIKey{
				Key:    apiKey,
				Secret: secretKey,
				UserID: authUserId,
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

		authUserId := c.GetUint("authUserId")
		if err := h.apiKeyService.DeleteAPIKey(ctx, authUserId, input.Key); err != nil {
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
		authUserId := c.GetUint("authUserId")

		if !isAdmin && authUserId != uint(userIdUint) {
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
