package transaction

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	v1 "banking/app/api/restful/v1"
	transactionRepo "banking/app/repo/mysql/transaction"
	"banking/domain"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go.elastic.co/apm/v2"
)

type TransactionHandler struct {
	transactionService domain.ITransactionService
}

func NewTransactionHandler(TransactionService domain.ITransactionService) domain.ITransactionHandler {
	return &TransactionHandler{
		transactionService: TransactionService,
	}
}

func (h *TransactionHandler) Transfer() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := apm.StartSpan(c.Request.Context(), "TransactionHandler.Transfer", "handler")
		defer span.End()

		var input struct {
			FromUserID uint    `json:"fromUserId" binding:"required,min=1,number"`
			ToUserID   uint    `json:"toUserId" binding:"required,min=1,number"`
			Amount     float64 `json:"amount" binding:"required,gt=0,number"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			apm.CaptureError(ctx, err).Send()
			c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrResponse{
				Msg: err.Error(),
			})
			return
		}

		if input.FromUserID != c.GetUint("authedUserId") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &v1.ErrResponse{
				Msg: "fromUserId is not authorized",
			})
			return
		}

		if input.FromUserID == input.ToUserID {
			c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrResponse{
				Msg: "fromUserId and toUserId should not be the same",
			})
			return
		}

		transaction, err := h.transactionService.Transfer(ctx, input.FromUserID, input.ToUserID, decimal.NewFromFloat(input.Amount))
		if err != nil {
			if errors.Is(err, transactionRepo.ErrInsufficientBalance) {
				apm.CaptureError(ctx, err).Send()
				c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrResponse{
					Msg: transactionRepo.ErrInsufficientBalance.Error(),
				})
				return
			}

			apm.CaptureError(ctx, err).Send()
			c.AbortWithStatusJSON(http.StatusInternalServerError, &v1.ErrResponse{
				Msg: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, &TransferResp{
			Data: &Transaction{
				FromUserID:      transaction.FromUserID,
				FromUserBalance: transaction.FromUserBalance,
				ToUserID:        transaction.ToUserID,
				Amount:          transaction.Amount,
				TransactionType: transaction.TransactionType,
				Details:         transaction.Details,
			},
		})
	}
}

func (h *TransactionHandler) Deposit() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := apm.StartSpan(c.Request.Context(), "TransactionHandler.Deposit", "handler")
		defer span.End()

		var input struct {
			UserID uint    `json:"userId" binding:"required,min=1,number"`
			Amount float64 `json:"amount" binding:"required,gt=0,number"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			apm.CaptureError(ctx, err).Send()
			c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrResponse{
				Msg: err.Error(),
			})
			return
		}

		authedUserID := c.GetUint("authedUserId")
		if input.UserID != authedUserID {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &v1.ErrResponse{
				Msg: "userId is not authorized",
			})
			return
		}

		transaction, err := h.transactionService.Deposit(ctx, input.UserID, decimal.NewFromFloat(input.Amount))
		if err != nil {
			apm.CaptureError(ctx, err).Send()
			c.AbortWithStatusJSON(http.StatusInternalServerError, &v1.ErrResponse{
				Msg: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, &DepositResp{
			Data: &Transaction{
				FromUserID:      transaction.FromUserID,
				FromUserBalance: transaction.FromUserBalance,
				ToUserID:        transaction.ToUserID,
				ToUserBalance:   transaction.ToUserBalance,
				Amount:          transaction.Amount,
				TransactionType: transaction.TransactionType,
				Details:         transaction.Details,
			},
		})
	}
}

func (h *TransactionHandler) Withdraw() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := apm.StartSpan(c.Request.Context(), "TransactionHandler.Withdraw", "handler")
		defer span.End()

		var input struct {
			UserID uint    `json:"userId" binding:"required,min=1,number"`
			Amount float64 `json:"amount" binding:"required,gt=0,number"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			apm.CaptureError(ctx, err).Send()
			c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrResponse{
				Msg: err.Error(),
			})
			return
		}

		authedUserID := c.GetUint("authedUserId")
		if input.UserID != authedUserID {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &v1.ErrResponse{
				Msg: "userId is not authorized",
			})
			return
		}

		transaction, err := h.transactionService.Withdraw(ctx, input.UserID, decimal.NewFromFloat(input.Amount))
		if err != nil {
			if errors.Is(err, transactionRepo.ErrInsufficientBalance) {
				apm.CaptureError(ctx, err).Send()
				c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrResponse{
					Msg: transactionRepo.ErrInsufficientBalance.Error(),
				})
				return
			}

			apm.CaptureError(ctx, err).Send()
			c.AbortWithStatusJSON(http.StatusInternalServerError, &v1.ErrResponse{
				Msg: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, &WithdrawResp{
			Data: &Transaction{
				FromUserID:      transaction.FromUserID,
				FromUserBalance: transaction.FromUserBalance,
				ToUserID:        transaction.ToUserID,
				ToUserBalance:   transaction.ToUserBalance,
				Amount:          transaction.Amount,
				TransactionType: transaction.TransactionType,
				Details:         transaction.Details,
			},
		})
	}
}

func (h *TransactionHandler) GetTransactions() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := apm.StartSpan(c.Request.Context(), "TransactionHandler.GetTransactions", "handler")
		defer span.End()

		userId := c.Param("userId")
		authedUserID := c.GetUint("authedUserId")

		userIdUint, err := strconv.ParseUint(userId, 10, 64)
		if err != nil {
			apm.CaptureError(ctx, err).Send()
			c.AbortWithStatusJSON(http.StatusBadRequest, &v1.ErrResponse{
				Msg: "invalid user id",
			})
			return
		}

		if uint(userIdUint) != authedUserID {
			apm.CaptureError(ctx, fmt.Errorf("unauthorized")).Send()
			c.AbortWithStatusJSON(http.StatusUnauthorized, &v1.ErrResponse{
				Msg: "unauthorized",
			})
			return
		}

		transactions, err := h.transactionService.GetTransactions(ctx, uint(userIdUint))
		if err != nil {
			apm.CaptureError(ctx, err).Send()
			c.AbortWithStatusJSON(http.StatusInternalServerError, &v1.ErrResponse{
				Msg: err.Error(),
			})
			return
		}

		transactionList := make([]*Transaction, 0, len(transactions))
		for _, t := range transactions {
			transactionList = append(transactionList, &Transaction{
				FromUserID:      t.FromUserID,
				FromUserBalance: t.FromUserBalance,
				ToUserID:        t.ToUserID,
				ToUserBalance:   t.ToUserBalance,
				Amount:          t.Amount,
				TransactionType: t.TransactionType,
				Details:         t.Details,
			})
		}
		c.JSON(http.StatusOK, &GetTransactionsResp{
			Data: transactionList,
		})
	}
}
